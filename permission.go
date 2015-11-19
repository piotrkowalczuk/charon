package charon

import (
	"bytes"
	"fmt"
)

// Permission is a string that consist of subsystem, module/content type and an action.
type Permission string

// NewPermission allocate new Permission object using given string.
func NewPermission(s string) Permission {
	return Permission(s)
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
