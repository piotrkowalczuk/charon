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
