package charond

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-ldap/ldap"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/pqt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type loginHandler struct {
	*handler
	hasher password.Hasher
}

func (lh *loginHandler) Login(ctx context.Context, r *charonrpc.LoginRequest) (*wrappers.StringValue, error) {
	lh.logger = log.NewContext(lh.logger).With("username", r.Username)

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
		var filter string
		if strings.Contains(r.Username, "@") {
			filter = fmt.Sprintf("(&(objectClass=organizationalPerson)(mail=%s))", ldap.EscapeFilter(r.Username))
		} else {
			parts := strings.Split(r.Username, "@")
			if len(parts) != 2 {
				return nil, grpc.Errorf(codes.InvalidArgument, "invalid email address")
			}
			filter = fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(parts[0]))
		}
		res, err := lh.ldap.Search(ldap.NewSearchRequest(
			lh.opts.LDAPSearchDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			filter,
			[]string{"dn", "givenName", "sn", "mail"},
			nil,
		))
		if err != nil {
			return nil, fmt.Errorf("ldap search failure: %s", err.Error())
		}

		if len(res.Entries) != 1 {
			return nil, grpc.Errorf(codes.Unauthenticated, "user does not exist, number of ldap entries found: %d", len(res.Entries))
		}

		conn, err := ldap.Dial("tcp", lh.ldapAddress)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "ldap connection failure: %s", err.Error())
		}
		defer conn.Close()

		if err = conn.Bind(res.Entries[0].DN, r.Password); err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "ldap bind failure: %s", err.Error())
		}

		username := res.Entries[0].GetAttributeValue("mail")
		usr, err = lh.repository.user.FindOneByUsername(username)
		if err != nil {
			if err == sql.ErrNoRows && lh.opts.LDAP {
				usr, err = lh.repository.user.Create(
					username,
					model.ExternalPassword,
					res.Entries[0].GetAttributeValue("givenName"),
					res.Entries[0].GetAttributeValue("sn"),
					[]byte(model.UserConfirmationTokenUsed),
					true,
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
			} else {
				return nil, grpc.Errorf(codes.Unauthenticated, "the username and password do not match")
			}
		}
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

	lh.loggerWith(
		"is_confirmed", usr.IsConfirmed,
		"is_staff", usr.IsStaff,
		"is_superuser", usr.IsSuperuser,
		"is_active", usr.IsActive,
		"first_name", usr.FirstName,
		"last_name", usr.LastName,
	)
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

	lh.loggerWith("token", res.Session.AccessToken)

	_, err = lh.repository.user.UpdateLastLoginAt(usr.ID)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "last login update failure: %s", err)
	}

	return &wrappers.StringValue{Value: res.Session.AccessToken}, nil
}
