package charond

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	libldap "github.com/go-ldap/ldap"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/ldap"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/pqt"
	"github.com/piotrkowalczuk/qtypes"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type loginHandler struct {
	*handler
	hasher   password.Hasher
	mappings *ldap.Mappings
}

func (lh *loginHandler) Login(ctx context.Context, r *charonrpc.LoginRequest) (*wrappers.StringValue, error) {
	if r.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty username")
	}
	if len(r.Password) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty password")
	}

	var (
		err error
		usr *model.UserEntity
	)
	if lh.opts.LDAP {
		got := lh.ldap.Get()
		if err, ok := got.(error); ok {
			return nil, err
		}

		conn := got.(*libldap.Conn)
		usr, err = lh.handleLDAP(conn, r)
		if err != nil {
			if terr, ok := err.(*libldap.Error); ok && terr.ResultCode >= libldap.ErrorNetwork {
				// on network issue, try once again
				<-time.After(1 * time.Second)
				if usr, err = lh.handleLDAP(conn, r); err != nil {
					conn.Close()
					return nil, err
				}
			}
			lh.ldap.Put(conn)
			return nil, err
		}
		lh.ldap.Put(conn)
	} else {
		usr, err = lh.repository.user.FindOneByUsername(r.Username)
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "user does not exists")
		}
	}

	if !lh.opts.LDAP {
		if bytes.Equal(usr.Password, model.ExternalPassword) {
			return nil, grpc.Errorf(codes.FailedPrecondition, "authentication failure, ldap is required")
		}
		if matches := lh.hasher.Compare(usr.Password, []byte(r.Password)); !matches {
			return nil, grpc.Errorf(codes.Unauthenticated, "the username and password do not match")
		}
	}

	if !usr.IsConfirmed {
		return nil, grpc.Errorf(codes.Unauthenticated, "user is not confirmed")
	}

	if !usr.IsActive {
		return nil, grpc.Errorf(codes.Unauthenticated, "user is not active")
	}

	res, err := lh.session.Start(ctx, &mnemosynerpc.StartRequest{
		Session: &mnemosynerpc.Session{
			SubjectId:     session.ActorIDFromInt64(usr.ID).String(),
			SubjectClient: r.Client,
			Bag: map[string]string{
				"username":   usr.Username,
				"first_name": usr.FirstName,
				"last_name":  usr.LastName,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = lh.repository.user.UpdateLastLoginAt(usr.ID)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "last login update failure: %s", err)
	}

	return &wrappers.StringValue{Value: res.Session.AccessToken}, nil
}

func (lh *loginHandler) handleLDAP(conn *libldap.Conn, r *charonrpc.LoginRequest) (*model.UserEntity, error) {
	var filter string
	if strings.Contains(r.Username, "@") {
		filter = fmt.Sprintf("(&(objectClass=organizationalPerson)(mail=%s))", libldap.EscapeFilter(r.Username))
	} else {
		parts := strings.Split(r.Username, "@")
		if len(parts) != 2 {
			return nil, grpc.Errorf(codes.InvalidArgument, "invalid email address")
		}
		filter = fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", libldap.EscapeFilter(parts[0]))
	}
	res, err := conn.Search(libldap.NewSearchRequest(
		lh.opts.LDAPSearchDN,
		libldap.ScopeWholeSubtree, libldap.NeverDerefAliases, 0, 0, false,
		filter,
		append([]string{"dn", "givenName", "sn", "mail"}, lh.mappings.Attributes...),
		nil,
	))
	if err != nil {
		return nil, fmt.Errorf("ldap search failure: %s", err.Error())
	}

	if len(res.Entries) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "user does not exist, number of ldap entries found: %d", len(res.Entries))
	}

	if err = conn.Bind(res.Entries[0].DN, r.Password); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "ldap bind failure: %s", err.Error())
	}

	username := res.Entries[0].GetAttributeValue("mail")

	sklog.Debug(lh.logger, "user found in ldap", "username_given", r.Username, "username_got", username, "attributes_got", res.Entries[0].Attributes, "attributes_asked", lh.mappings.Attributes)

	var usr *model.UserEntity
	usr, err = lh.repository.user.FindOneByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows && lh.opts.LDAP {
			sklog.Debug(lh.logger, "user found in ldap but does not exists in database", "username", username)

			usr, err = lh.repository.user.Create(
				username,
				model.ExternalPassword,
				res.Entries[0].GetAttributeValue("givenName"),
				res.Entries[0].GetAttributeValue("sn"),
				[]byte(model.UserConfirmationTokenUsed),
				false,
				false,
				true,
				true,
			)
			if err != nil {
				switch pqt.ErrorConstraint(err) {
				case model.TableUserConstraintPrimaryKey:
					return nil, grpc.Errorf(codes.AlreadyExists, "user with such id already exists")
				case model.TableUserConstraintUsernameUnique:
					return nil, grpc.Errorf(codes.AlreadyExists, "user with such username already exists")
				default:
					return nil, err
				}
			}
			sklog.Debug(lh.logger, "user found in ldap has been created", "username", username, "first_name", usr.FirstName, "last_name", usr.LastName)
		} else {
			return nil, grpc.Errorf(codes.Unauthenticated, "the username and password do not match")
		}
	}
	if groups, permissions, ok := lh.mappings.Map(res.Entries[0].Attributes); ok {
		sklog.Debug(lh.logger, "ldap mapping found", "username", username, "count_groups", len(groups), "count_permissions", len(permissions))

		if len(permissions) > 0 {
			inserted, _, err := lh.repository.user.SetPermissions(usr.ID, charon.NewPermissions(permissions...)...)
			if err != nil {
				return nil, err
			}
			sklog.Debug(lh.logger, "permissions given to the user", "username", username, "user_id", usr.ID, "inserted", inserted)
		} else {
			if _, err = lh.repository.userPermissions.DeleteByUserID(usr.ID); err != nil {
				return nil, err
			}
		}

		if len(groups) > 0 {
			groupsFound, err := lh.repository.group.Find(&model.GroupCriteria{
				Name: &qtypes.String{
					Values: groups,
					Type:   qtypes.QueryType_IN,
					Valid:  true,
				},
			})
			if err != nil {
				return nil, err
			}
			var ids []int64
			for _, g := range groupsFound {
				ids = append(ids, g.ID)
			}
			_, _, err = lh.repository.userGroups.Set(usr.ID, ids)
			if err != nil {
				return nil, err
			}
			sklog.Debug(lh.logger, "user added to groups", "user_id", usr.ID, "count", len(ids))
		} else {
			if _, err = lh.repository.userGroups.DeleteByUserID(usr.ID); err != nil {
				return nil, err
			}
		}
	} else {
		sklog.Debug(lh.logger, "ldap mapping not found", "username", username, "attributes_got", res.Entries[0].Attributes, "attributes_asked", lh.mappings.Attributes)
	}

	return usr, nil
}
