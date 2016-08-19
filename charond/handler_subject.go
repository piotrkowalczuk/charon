package charond

import (
	"database/sql"
	"fmt"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type subjectHandler struct {
	*handler
}

func (sh *subjectHandler) handle(ctx context.Context, r *charon.SubjectRequest) (*charon.SubjectResponse, error) {
	var (
		ses *mnemosynerpc.Session
		err error
	)
	if r.AccessToken == "" {
		res, err := sh.session.Context(ctx, none())
		if err != nil {
			return nil, handleMnemosyneError(err)
		}
		ses = res.Session
	} else {
		res, err := sh.session.Get(ctx, &mnemosynerpc.GetRequest{
			AccessToken: r.AccessToken,
		})
		if err != nil {
			return nil, handleMnemosyneError(err)
		}
		ses = res.Session
	}

	id, err := charon.SubjectID(ses.SubjectId).UserID()
	if err != nil {
		return nil, fmt.Errorf("invalid session subject id: %s", ses.SubjectId)
	}

	ent, err := sh.repository.user.findOneByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "subject does not exists with id: %d", id)
		}

		return nil, err
	}

	permissionEntities, err := sh.repository.permission.findByUserID(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, grpc.Errorf(codes.Internal, "subject list of permissions failure: %s", err.Error())
	}

	permissions := make([]string, 0, len(permissionEntities))
	for _, e := range permissionEntities {
		permissions = append(permissions, e.Permission().String())
	}

	sh.loggerWith("subject_id", ses.SubjectId)

	return &charon.SubjectResponse{
		Id:          int64(ent.id),
		Username:    ent.username,
		FirstName:   ent.firstName,
		LastName:    ent.lastName,
		Permissions: permissions,
		IsActive:    ent.isActive,
		IsConfirmed: ent.isConfirmed,
		IsStuff:     ent.isStaff,
		IsSuperuser: ent.isSuperuser,
	}, nil
}
