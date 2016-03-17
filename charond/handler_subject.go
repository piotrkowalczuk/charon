package main

import (
	"database/sql"
	"fmt"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type subjectHandler struct {
	*handler
}

func (sh *subjectHandler) handle(ctx context.Context, r *charon.SubjectRequest) (*charon.SubjectResponse, error) {
	var (
		ses *mnemosyne.Session
		err error
	)
	if r.AccessToken == nil {
		if ses, err = sh.session.FromContext(ctx); err != nil {
			return nil, err
		}
	} else {
		if ses, err = sh.session.Get(ctx, *r.AccessToken); err != nil {
			return nil, err
		}
	}

	id, err := charon.SubjectID(ses.SubjectId).UserID()
	if err != nil {
		return nil, fmt.Errorf("charond: invalid session subject id: %s", ses.SubjectId)
	}

	userEntity, err := sh.repository.user.FindOneByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "charond: subject does not exists with id: %d", id)
		}

		return nil, err
	}

	permissionEntities, err := sh.repository.permission.FindByUserID(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, grpc.Errorf(codes.Internal, "charond: subject list of permissions failure: %s", err.Error())
	}

	permissions := make([]string, 0, len(permissionEntities))
	for _, e := range permissionEntities {
		permissions = append(permissions, e.Permission().String())
	}

	sh.loggerWith("subject_id", ses.SubjectId)

	return &charon.SubjectResponse{
		Id:          int64(userEntity.ID),
		Username:    userEntity.Username,
		FirstName:   userEntity.FirstName,
		LastName:    userEntity.LastName,
		Permissions: permissions,
		IsActive:    userEntity.IsActive,
		IsConfirmed: userEntity.IsConfirmed,
		IsStuff:     userEntity.IsStaff,
		IsSuperuser: userEntity.IsSuperuser,
	}, nil
}
