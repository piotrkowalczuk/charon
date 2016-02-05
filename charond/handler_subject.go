package main

import (
	"database/sql"
	"fmt"

	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type subjectHandler struct {
	*handler
}

func (sh *subjectHandler) handle(ctx context.Context, r *charon.SubjectRequest) (*charon.SubjectResponse, error) {
	var (
		ok bool
		md metadata.MD
	)
	if md, ok = metadata.FromContext(ctx); ok {
		grpc.Header(&md)
	}

	ses, err := sh.session.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	id, err := charon.SessionSubjectID(ses.SubjectId).UserID()
	if err != nil {
		return nil, fmt.Errorf("charond: invalid session subject id: %s", ses.SubjectId)
	}

	user, err := sh.repository.user.FindOneByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "charond: user does not exists with id: %d", id)
		}

		return nil, err
	}

	permissionEntities, err := sh.repository.permission.FindByUserID(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, grpc.Errorf(codes.Internal, "charond: list of permissions failure: %s", err.Error())
	}

	permissions := make([]string, 0, len(permissionEntities))
	for _, e := range permissionEntities {
		permissions = append(permissions, e.Permission().String())
	}

	sh.loggerWith("subject_id", ses.SubjectId)

	return &charon.SubjectResponse{
		Id:          int64(user.ID),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Permissions: permissions,
		IsActive:    user.IsActive,
		IsConfirmed: user.IsConfirmed,
		IsStuff:     user.IsStaff,
		IsSuperuser: user.IsSuperuser,
	}, nil
}
