package charon

import (
	"bytes"
	"fmt"
)

const (
	UserCanCreate      Permission = "charon:user:can create"
	UserCanCreateStaff Permission = "charon:user:can create staff"

	UserCanDeleteAsStranger      Permission = "charon:user:can delete as stranger"
	UserCanDeleteAsOwner         Permission = "charon:user:can delete as owner"
	UserCanDeleteStaffAsStranger Permission = "charon:user:can delete staff as stranger"
	UserCanDeleteStaffAsOwner    Permission = "charon:user:can delete staff as owner"

	UserCanModifyAsStranger      Permission = "charon:user:can modify as stranger"
	UserCanModifyAsOwner         Permission = "charon:user:can modify as owner"
	UserCanModifyStaffAsStranger Permission = "charon:user:can modify staff as stranger"
	UserCanModifyStaffAsOwner    Permission = "charon:user:can modify staff as owner"

	UserCanRetrieveAsOwner         Permission = "charon:user:can retrieve as owner"
	UserCanRetrieveAsStranger      Permission = "charon:user:can retrieve as stranger"
	UserCanRetrieveStaffAsOwner    Permission = "charon:user:can retrieve staff as owner"
	UserCanRetrieveStaffAsStranger Permission = "charon:user:can retrieve staff as stranger"

	UserPermissionCanCreate   Permission = "charon:user_permission:can create"
	UserPermissionCanDelete   Permission = "charon:user_permission:can delete"
	UserPermissionCanModify   Permission = "charon:user_permission:can modify"
	UserPermissionCanRetrieve Permission = "charon:user_permission:can retrieve"

	PermissionCanCreate   Permission = "charon:permission:can create"
	PermissionCanDelete   Permission = "charon:permission:can delete"
	PermissionCanModify   Permission = "charon:permission:can modify"
	PermissionCanRetrieve Permission = "charon:permission:can retrieve"

	GroupCanCreate   Permission = "charon:group:can create"
	GroupCanDelete   Permission = "charon:group:can delete"
	GroupCanModify   Permission = "charon:group:can modify"
	GroupCanRetrieve Permission = "charon:group:can retrieve"

	GroupPermissionCanCreate   Permission = "charon:group_permission:can create"
	GroupPermissionCanDelete   Permission = "charon:group_permission:can delete"
	GroupPermissionCanModify   Permission = "charon:group_permission:can modify"
	GroupPermissionCanRetrieve Permission = "charon:group_permission:can retrieve"
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
func NewPermissions(ss ...string) Permissions {
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
