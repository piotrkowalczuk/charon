package charond

import (
	"sort"
	"testing"

	"reflect"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListGroupPermissionsHandler_ListPermissions(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	createGroupResp, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
		Name: "existing-group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	_, err = suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: testPermissionsDataUserService,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: testPermissionsDataCustomerService,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: testPermissionsDataBigService,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: testPermissionsDataImageService,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	dataset := [][]string{
		testPermissionsDataUserService,
	}
	dataset = append(dataset, append(dataset[0], testPermissionsDataCustomerService...))
	dataset = append(dataset, append(dataset[1], testPermissionsDataBigService...))
	dataset = append(dataset, append(dataset[2], testPermissionsDataImageService...))

	for _, permissions := range dataset {
		_, err = suite.charon.group.SetPermissions(ctx, &charonrpc.SetGroupPermissionsRequest{
			GroupId:     createGroupResp.Group.Id,
			Permissions: permissions,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		res, err := suite.charon.group.ListPermissions(ctx, &charonrpc.ListGroupPermissionsRequest{
			Id: createGroupResp.Group.Id,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() != codes.NotFound {
					t.Fatalf("wrong error code, expected %s but got %s for error: %s", codes.NotFound, st.Code(), err.Error())
				}
			} else {
				t.Errorf("wrong error type: %T", err)
			}
		}

		sort.Strings(res.Permissions)
		sort.Strings(permissions)
		if !reflect.DeepEqual(permissions, res.Permissions) {
			t.Errorf("wrong permissions returend, expected:\n	%v\nbut got:\n	 %v", permissions, res.Permissions)
		} else {
			t.Logf("equal number of permissions: %d", len(res.Permissions))
		}
	}
}

func TestListGroupPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.ListGroupPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.GroupPermissionCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &listGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListGroupPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.ListGroupPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{},
				Permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
					charon.UserPermissionCanRetrieve,
				},
			},
		},
	}

	h := &listGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}

var testPermissionsDataCustomerService = []string{
	"customerserv:account:can create",
	"customerserv:account:can modify as creator",
	"customerserv:account:can modify as manager",
	"customerserv:account:can delete as creator",
	"customerserv:account:can delete as manager",
	"customerserv:account:can retrieve as stranger",

	"customerserv:profile:can create",
	"customerserv:profile:can modify as creator",
	"customerserv:profile:can modify as manager",
	"customerserv:profile:can delete as creator",
	"customerserv:profile:can delete as manager",
	"customerserv:profile:can retrieve as stranger",

	"customerserv:contract:can create",
	"customerserv:contract:can modify as creator",
	"customerserv:contract:can modify as manager",
	"customerserv:contract:can delete as creator",
	"customerserv:contract:can delete as manager",
	"customerserv:contract:can retrieve as stranger",

	"customerserv:package:can create",
	"customerserv:package:can modify as creator",
	"customerserv:package:can modify as manager",
	"customerserv:package:can delete as creator",
	"customerserv:package:can delete as manager",
	"customerserv:package:can retrieve as stranger",

	"customerserv:campaign:can create",
	"customerserv:campaign:can modify as creator",
	"customerserv:campaign:can modify as manager",
	"customerserv:campaign:can delete as creator",
	"customerserv:campaign:can delete as manager",
	"customerserv:campaign:can retrieve as stranger",
}
var testPermissionsDataUserService = []string{
	"userserv:account:can create",
	"userserv:account:can modify as creator",
	"userserv:account:can modify as manager",
	"userserv:account:can delete as creator",
	"userserv:account:can delete as manager",
	"userserv:account:can retrieve as creator",
	"userserv:account:can retrieve as stranger",
	"userserv:account:can retrieve as manager",

	"userserv:profile:can create",
	"userserv:profile:can create as parent entity creator",
	"userserv:profile:can create as parent entity manager",
	"userserv:profile:can modify as creator",
	"userserv:profile:can modify as manager",
	"userserv:profile:can modify as parent entity creator",
	"userserv:profile:can modify as parent entity manager",
	"userserv:profile:can delete as creator",
	"userserv:profile:can delete as manager",
	"userserv:profile:can retrieve as creator",
	"userserv:profile:can retrieve as stranger",
	"userserv:profile:can retrieve as manager",

	"userserv:url:can create",
	"userserv:url:can create as parent entity creator",
	"userserv:url:can create as parent entity manager",
	"userserv:url:can modify as creator",
	"userserv:url:can modify as manager",
	"userserv:url:can modify as parent entity creator",
	"userserv:url:can modify as parent entity manager",
	"userserv:url:can delete as creator",
	"userserv:url:can delete as manager",
	"userserv:url:can retrieve as creator",
	"userserv:url:can retrieve as stranger",
	"userserv:url:can retrieve as manager",
}
var testPermissionsDataImageService = []string{
	"imageserv:file:can create",
	"imageserv:file:can modify as creator",
	"imageserv:file:can modify as stranger",
	"imageserv:file:can delete as creator",
	"imageserv:file:can delete as stranger",
	"imageserv:file:can retrieve as creator",
	"imageserv:file:can retrieve as stranger",

	"imageserv:mime-type:can create",
	"imageserv:mime-type:can modify as creator",
	"imageserv:mime-type:can modify as stranger",
	"imageserv:mime-type:can delete as creator",
	"imageserv:mime-type:can delete as stranger",
	"imageserv:mime-type:can retrieve as creator",
	"imageserv:mime-type:can retrieve as stranger",

	"imageserv:template:can create",
	"imageserv:template:can modify as creator",
	"imageserv:template:can modify as stranger",
	"imageserv:template:can delete as creator",
	"imageserv:template:can delete as stranger",
	"imageserv:template:can retrieve as creator",
	"imageserv:template:can retrieve as stranger",

	"imageserv:css:can create",
	"imageserv:css:can modify as creator",
	"imageserv:css:can modify as stranger",
	"imageserv:css:can delete as creator",
	"imageserv:css:can delete as stranger",
	"imageserv:css:can retrieve as creator",
	"imageserv:css:can retrieve as stranger",

	"imageserv:settings:can create as stranger",
	"imageserv:settings:can create as manager",
	"imageserv:settings:can retrieve as stranger",
	"imageserv:settings:can retrieve as creator",
	"imageserv:settings:can retrieve as manager",
	"imageserv:settings:can delete as stranger",
	"imageserv:settings:can delete as creator",
	"imageserv:settings:can delete as manager",
	"imageserv:settings:can modify as stranger",
	"imageserv:settings:can modify as creator",
	"imageserv:settings:can modify as manager",
}
var testPermissionsDataBigService = []string{
	"bigserv:forum-parameter:can create",
	"bigserv:forum-parameter:can modify as creator",
	"bigserv:forum-parameter:can modify as stranger",
	"bigserv:forum-parameter:can delete as creator",
	"bigserv:forum-parameter:can delete as stranger",
	"bigserv:forum-parameter:can retrieve as creator",
	"bigserv:forum-parameter:can retrieve as stranger",

	"bigserv:forum:can create",
	"bigserv:forum:can modify as creator",
	"bigserv:forum:can modify as stranger",
	"bigserv:forum:can delete as creator",
	"bigserv:forum:can delete as stranger",
	"bigserv:forum:can retrieve as creator",
	"bigserv:forum:can retrieve as stranger",

	"bigserv:post:can create",
	"bigserv:post:can modify as creator",
	"bigserv:post:can modify as stranger",
	"bigserv:post:can delete as creator",
	"bigserv:post:can delete as stranger",
	"bigserv:post:can retrieve as creator",
	"bigserv:post:can retrieve as stranger",

	"bigserv:comment:can create",
	"bigserv:comment:can modify as creator",
	"bigserv:comment:can modify as stranger",
	"bigserv:comment:can delete as creator",
	"bigserv:comment:can delete as stranger",
	"bigserv:comment:can retrieve as creator",
	"bigserv:comment:can retrieve as stranger",
}
