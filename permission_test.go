package charon

import (
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

func TestPermissions_Contains(t *testing.T) {
	positive := map[Permission]Permissions{
		UserCanCreate: {
			UserCanCreate,
			UserCanDeleteAsOwner,
		},
		UserPermissionCanCreate: {
			UserPermissionCanCreate,
		},
	}

	for expected, permissions := range positive {
		if !permissions.Contains(expected) {
			t.Errorf("expected permission (%s), is not present", expected)
		}
	}

	negative := map[Permission]Permissions{
		PermissionCanCreate: {
			UserCanCreate,
			UserCanDeleteAsOwner,
		},
		UserCanDeleteAsStranger: {
			UserPermissionCanCreate,
		},
		UserCanCreate: {},
	}

	for unexpected, permissions := range negative {
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
