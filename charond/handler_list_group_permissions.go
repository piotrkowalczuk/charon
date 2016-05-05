package charond

import (
	"database/sql"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
)

type listGroupPermissionsHandler struct {
	*handler
}

func (luph *listGroupPermissionsHandler) handle(ctx context.Context, req *charon.ListGroupPermissionsRequest) (*charon.ListGroupPermissionsResponse, error) {
	luph.loggerWith("group_id", req.Id)

	permissions, err := luph.repository.permission.FindByGroupID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			sklog.Debug(luph.logger, "group permissions retrieved", "group_id", req.Id, "count", len(permissions))

			return &charon.ListGroupPermissionsResponse{}, nil
		}
		return nil, err
	}

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Permission().String())
	}

	luph.loggerWith("results", len(permissions))

	return &charon.ListGroupPermissionsResponse{
		Permissions: perms,
	}, nil
}

func (luph *listGroupPermissionsHandler) firewall(req *charon.ListGroupPermissionsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.GroupPermissionCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: list of group permissions cannot be retrieved, missing permission")
}
