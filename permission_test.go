package charon

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPermission_String(t *testing.T) {
	data := map[string]Permission{
		"charon:user:can create":       UserCanCreate,
		"charon:user_group:can delete": UserGroupCanDelete,
	}

	for expected, d := range data {
		if d.String() != expected {
			t.Errorf("wrong output, expected %s but got %s", expected, d.String())
		}
	}
}

func TestPermission_Split(t *testing.T) {
	data := []struct {
		subsystem, module, action string
		permission                Permission
	}{
		{},
		{
			subsystem:  "charon",
			module:     "user",
			action:     "can create",
			permission: UserCanCreate,
		},
		{
			subsystem:  "charon",
			module:     "user_permission",
			action:     "can check granting as a stranger",
			permission: UserPermissionCanCheckGrantingAsStranger,
		},
		{
			module:     "module",
			action:     "action as someone",
			permission: Permission("module:action as someone"),
		},
		{
			action:     "action as somebody else",
			permission: Permission("action as somebody else"),
		},
	}

	for _, d := range data {
		subsystem, module, action := d.permission.Split()

		if d.subsystem != subsystem {
			t.Errorf("wrong subsystem, expected %s but got %s", d.subsystem, subsystem)
		}
		if d.module != module {
			t.Errorf("wrong module, expected %s but got %s", d.module, module)
		}
		if d.action != action {
			t.Errorf("wrong action, expected %s but got %s", d.action, action)
		}
	}
}

func TestPermission_Subsystem(t *testing.T) {
	data := map[string]Permission{
		"charon": UserPermissionCanCheckGrantingAsStranger,
	}

	for expected, d := range data {
		if expected != d.Subsystem() {
			t.Errorf("wrong subsystem, expected %s but got %s", expected, d.Subsystem())
		}
	}
}

func TestPermission_Module(t *testing.T) {
	data := map[string]Permission{
		"user_permission": UserPermissionCanCheckGrantingAsStranger,
	}

	for expected, d := range data {
		if expected != d.Module() {
			t.Errorf("wrong module, expected %s but got %s", expected, d.Module())
		}
	}
}

func TestPermission_Action(t *testing.T) {
	data := map[string]Permission{
		"can check granting as a stranger": UserPermissionCanCheckGrantingAsStranger,
	}

	for expected, d := range data {
		if expected != d.Action() {
			t.Errorf("wrong action, expected %s but got %s", expected, d.Action())
		}
	}
}

func TestPermission_MarshalJSON(t *testing.T) {
	data := map[string]Permission{
		`""`:                                                        Permission(""),
		`"some action"`:                                             Permission("some action"),
		`"module:some action"`:                                      Permission("module:some action"),
		`"charon:user_permission:can check granting as a stranger"`: UserPermissionCanCheckGrantingAsStranger,
	}

	for expected, d := range data {
		b, err := json.Marshal(d)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}
		if expected != string(b) {
			t.Errorf("wrong json output, expected %s but got %s", expected, string(b))
		}
	}
}

func TestPermissions_Contains_true(t *testing.T) {
	data := map[Permission]Permissions{
		UserCanCreate: {
			UserCanCreate,
			UserCanDeleteAsOwner,
		},
		UserPermissionCanCreate: {
			UserPermissionCanCreate,
		},
	}

	for expected, permissions := range data {
		if !permissions.Contains(expected) {
			t.Errorf("expected permission (%s), is not present", expected)
		}
	}
}

func TestPermissions_Contains_false(t *testing.T) {
	data := map[Permission]Permissions{
		PermissionCanCreate: {
			UserCanCreate,
			UserCanDeleteAsOwner,
		},
		UserCanDeleteAsStranger: {
			UserPermissionCanCreate,
		},
		UserCanCreate: {},
	}

	for unexpected, permissions := range data {
		if permissions.Contains(unexpected) {
			t.Errorf("unexpected permission (%s), should not be present", unexpected)
		}
	}
}

func TestPermissions_Contains_many(t *testing.T) {
	var permissions Permissions

	permissions = Permissions{UserCanCreate, UserCanDeleteAsOwner}
	if c := permissions.Contains(UserCanCreate, UserCanDeleteAsOwner); !c {
		t.Errorf("both permission can be found, but ware not")
	}

	permissions = Permissions{UserCanCreate, UserCanDeleteAsOwner}
	if c := permissions.Contains(UserCanCreate, UserCanDeleteAsStranger); !c {
		t.Errorf("at least one of privided permission is there")
	}

	permissions = Permissions{UserCanCreate, UserCanDeleteAsOwner}
	if c := permissions.Contains(UserCanModifyAsOwner, UserCanDeleteAsStranger); c {
		t.Errorf("none of privided permission are present but was found")
	}
}

func TestPermissions_Strings(t *testing.T) {
	got := Permissions{UserCanCreate, UserCanDeleteAsOwner}.Strings()
	expected := []string{UserCanCreate.String(), UserCanDeleteAsOwner.String()}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("output does not match, expected %v but got %v", expected, got)
	}
}

var (
	permissionTestFixtures = []*permissionEntity{
		{
			Subsystem: "subsystem",
			Module:    "module",
			Action:    "action",
		},
		{
			Module: "module",
			Action: "action",
		},
		{
			Action: "action",
		},
	}
)

func TestPermissionRepository_FindOneByID(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range loadPermissionFixtures(t, suite.repository.permission, permissionTestFixtures) {
		found, err := suite.repository.permission.FindOneByID(res.got.ID)

		if err != nil {
			t.Errorf("permission cannot be found, unexpected error: %s", err.Error())
			continue
		}

		if assert(t, found != nil, "permission was not found, nil object returned") {
			assertf(t, reflect.DeepEqual(res.got, *found), "created and retrieved entity should be equal, but its not\ncreated: %#v\nfounded: %#v", res.got, found)
		}
	}
}

func TestPermissionRepository_Register(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	data := []struct {
		created, removed, untouched int64
		permissions                 Permissions
	}{
		{
			created:     int64(len(AllPermissions)),
			permissions: AllPermissions,
		},
		{
			untouched:   int64(len(AllPermissions)),
			permissions: AllPermissions,
		},
		{
			untouched: 1,
			removed:   int64(len(AllPermissions) - 1),
			permissions: Permissions{
				UserCanCreate,
			},
		},
		{
			created: 1,
			removed: 1,
			permissions: Permissions{
				Permission("charon:fakemodule:fakeaction"),
			},
		},
		{
			created: 1,
			permissions: Permissions{
				Permission("fakesystem:fakemodule:fakeaction"),
			},
		},
		{
			removed:     1,
			created:     int64(len(AllPermissions)),
			permissions: AllPermissions,
		},
	}

	for i, d := range data {
		created, untouched, removed, err := suite.repository.permission.Register(d.permissions)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		if created != d.created {
			t.Errorf("expected different number of created permissions, expected %d got %d for set %d", d.created, created, i)
		}
		if untouched != d.untouched {
			t.Errorf("expected different number of untouched permissions, expected %d got %d for set %d", d.untouched, untouched, i)
		}
		if removed != d.removed {
			t.Errorf("expected different number of removed permissions, expected %d got %d for set %d", d.removed, removed, i)
		}
	}
}

type permissionFixtures struct {
	got, given permissionEntity
}

func loadPermissionFixtures(t *testing.T, r permissionProvider, f []*permissionEntity) chan permissionFixtures {
	data := make(chan permissionFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("permission cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("permission has been created, got id %d", entity.ID)
			}

			data <- permissionFixtures{
				got:   *entity,
				given: *given,
			}
		}

		close(data)
	}()

	return data
}
