package charon

import (
	"bytes"
	"fmt"
)

const (
	UserCanCreate             Permission = "charon:user:can create"
	UserCanCreateSuper        Permission = "charon:user:can create, super"
	UserCanCreateStaff        Permission = "charon:user:can create, staff"
	UserCanCreateActive       Permission = "charon:user:can create, active"
	UserCanCreateConfirmed    Permission = "charon:user:can create, confirmed"
	UserCanDeleteAsStranger   Permission = "charon:user:can delete as stranger"
	UserCanDeleteAsStaff      Permission = "charon:user:can delete as staff"
	UserCanDeleteAsSuper      Permission = "charon:user:can delete as super"
	UserCanDeleteAsOwner      Permission = "charon:user:can delete as owner"
	UserCanEditAsStranger     Permission = "charon:user:can edity as stranger"
	UserCanEditAsOwner        Permission = "charon:user:can edity as owner"
	UserCanRetrieveAsOwner    Permission = "charon:user:can retrieve as owner"
	UserCanRetrieveAsStranger Permission = "charon:user:can retrieve as stranger"

	UserPermissionCanCreate             Permission = "charon:user_permission:can create"
	UserPermissionCanDelete             Permission = "charon:user_permission:can delete"
	UserPermissionCanEdit               Permission = "charon:user_permission:can edity"
	UserPermissionCanRetrieveAsOwner    Permission = "charon:user_permission:can retrieve as owner"
	UserPermissionCanRetrieveAsStranger Permission = "charon:user_permission:can retrieve as stranger"

	PermissionCanCreate             Permission = "charon:permission:can create"
	PermissionCanDelete             Permission = "charon:permission:can create"
	PermissionCanEdit               Permission = "charon:permission:can create"
	PermissionCanRetrieveAsOwner    Permission = "charon:permission:can create as owner"
	PermissionCanRetrieveAsStranger Permission = "charon:permission:can create as stranger"

	GroupCanCreate             Permission = "charon:group:can create"
	GroupCanDelete             Permission = "charon:group:can create"
	GroupCanEdit               Permission = "charon:group:can create"
	GroupCanRetrieveAsOwner    Permission = "charon:group:can create as owner"
	GroupCanRetrieveAsStranger Permission = "charon:group:can create as stranger"

	GroupPermissionCanCreate             Permission = "charon:group_permission:can create"
	GroupPermissionCanDelete             Permission = "charon:group_permission:can delete"
	GroupPermissionCanEdit               Permission = "charon:group_permission:can edity"
	GroupPermissionCanRetrieveAsOwner    Permission = "charon:group_permission:can retrieve as owner"
	GroupPermissionCanRetrieveAsStranger Permission = "charon:group_permission:can retrieve as stranger"
)

var (
	// EmptyPermission is a shorthand
	EmptyPermission = Permission("")
)

// Permission is a string that consist of subsystem, module/content type and an action.
type Permission string

// NewPermission allocate new Permission object using given string.
func NewPermission(s string) Permission {
	return Permission(s)
}

// String implements fmt.Stringer interface.
func (p Permission) String() string {
	return string(p)
}

// Split returns subsystem, module/content ty and action that describes single Permission.
func (p Permission) Split() (string, string, string) {
	if p == "" {
		return "", "", ""
	}

	parts := bytes.Split([]byte(p), []byte(":"))

	switch len(parts) {
	case 1:
		return "", "", string(parts[0])
	case 2:
		return "", string(parts[0]), string(parts[1])
	default:
		return string(parts[0]), string(parts[1]), string(parts[2])
	}
}

// Permission implements Permission interface.
func (p Permission) Permission() string {
	return string(p)
}

// MarshalJSON implements json.Marshaller interface.
func (p Permission) MarshalJSON() ([]byte, error) {
	return []byte(p), nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (p *Permission) UnmarshalJSON(src interface{}) error {
	switch s := src.(type) {
	case string:
		*p = Permission(s)
	case []byte:
		*p = Permission(s)
	default:
		return fmt.Errorf("charon: permission expects string or slice of bytes, got %T", src)
	}

	return nil
}

type Permissions []Permission

// NewPermissions allocates new Permissions using given slice of strings.
// It maps each string in a slice into Permission.
func NewPermissions(ss []string) Permissions {
	ps := make(Permissions, 0, len(ss))
	for _, s := range ss {
		ps = append(ps, Permission(s))
	}

	return ps
}

// Contains returns true if given Permission exists in the collection.
func (p Permissions) Contains(permission Permission) bool {
	for _, perm := range p {
		if perm == permission {
			return true
		}
	}

	return false
}
