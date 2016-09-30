package charond

import (
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-ldap/ldap"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type loginHandler struct {
	*handler
	hasher charon.PasswordHasher
}

func (lh *loginHandler) handle(ctx context.Context, r *charon.LoginRequest) (*charon.LoginResponse, error) {
	lh.logger = log.NewContext(lh.logger).With("username", r.Username)

	if r.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty username")
	}
	if len(r.Password) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty password")
	}

	if lh.opts.LDAP {
		parts := strings.Split(r.Username, "@")
		if len(parts) != 2 {
			return nil, grpc.Errorf(codes.InvalidArgument, "invalid email address")
		}
		res, err := lh.ldap.Search(ldap.NewSearchRequest(
			lh.opts.LDAPDistinguishedName,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(parts[0])),
			[]string{"dn"},
			nil,
		))
		if err != nil {
			return nil, err
		}

		if len(res.Entries) != 1 {
			return nil, grpc.Errorf(codes.Unauthenticated, "user does not exist, number of LDAP entries found: %d", len(res.Entries))
		}

		conn, err := ldap.Dial("tcp", lh.opts.LDAPAddress)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "ldap connection failure: %s", err.Error())
		}
		defer conn.Close()

		if err = conn.Bind(res.Entries[0].DN, r.Password); err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "the username and password do not match")
		}
	}
	usr, err := lh.repository.user.findOneByUsername(r.Username)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "the username and password do not match")
	}

	if !lh.opts.LDAP {
		if matches := lh.hasher.Compare(usr.password, []byte(r.Password)); !matches {
			return nil, grpc.Errorf(codes.Unauthenticated, "the username and password do not match")
		}
	}

	lh.loggerWith(
		"is_confirmed", usr.isConfirmed,
		"is_staff", usr.isStaff,
		"is_superuser", usr.isSuperuser,
		"is_active", usr.isActive,
		"first_name", usr.firstName,
		"last_name", usr.lastName,
	)
	if !usr.isConfirmed {
		return nil, grpc.Errorf(codes.Unauthenticated, "user is not confirmed")
	}

	if !usr.isActive {
		return nil, grpc.Errorf(codes.Unauthenticated, "user is not active")
	}

	res, err := lh.session.Start(ctx, &mnemosynerpc.StartRequest{
		Session: &mnemosynerpc.Session{
			SubjectId:     charon.SubjectIDFromInt64(usr.id).String(),
			SubjectClient: r.Client,
			Bag: map[string]string{
				"username":   usr.username,
				"first_name": usr.firstName,
				"last_name":  usr.lastName,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	lh.loggerWith("token", res.Session.AccessToken)

	_, err = lh.repository.user.updateLastLoginAt(usr.id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "last login update failure: %s", err)
	}

	return &charon.LoginResponse{AccessToken: res.Session.AccessToken}, nil
}
