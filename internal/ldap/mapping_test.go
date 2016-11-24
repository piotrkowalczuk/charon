package ldap_test

import (
	"reflect"
	"strings"
	"testing"

	libldap "github.com/go-ldap/ldap"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/ldap"
)

func TestNewMappings(t *testing.T) {
	m, err := ldap.NewMappings(strings.NewReader(`
[
  {
    "from": {
      "cn": ["cn_3", "cn_1"]
    },
    "to": {
     "groups": ["Admins"],
     "permissions": []
    }
  },
  {
    "from": {
      "ou": ["ou_2"],
      "dn": ["dn_1"]
    },
    "to": {
      "groups": ["Members"],
      "permissions": []
    }
  }
]
`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if !reflect.DeepEqual(m.Attributes, []string{"cn", "dn", "ou"}) {
		t.Error("attributes are not equal")
	}
	t.Logf("got attributes: %v", m.Attributes)
}

func TestMapping_Map(t *testing.T) {
	mappings := ldap.Mappings{
		Mappings: []ldap.Mapping{
			{
				From: map[string][]string{
					"cn": {"cn_1"},
					"dc": {"dc_1"},
					"ou": {"ou_1", "ou_2"},
				},
				To: ldap.MappingTo{
					Groups:      []string{"admins"},
					Permissions: []string{charon.UserCanCreate.String()},
				},
			},
			{
				From: map[string][]string{
					"cn": {"cn_1", "cn_100"},
					"dc": {"dc_1"},
					"ou": {"ou_1", "ou_2"},
				},
				To: ldap.MappingTo{
					Groups:      []string{"deleters"},
					Permissions: []string{charon.UserCanDeleteAsOwner.String()},
				},
			},
			{
				From: map[string][]string{
					"cn": {"cn_4"},
				},
				To: ldap.MappingTo{
					Groups:      []string{"users"},
					Permissions: []string{charon.UserCanRetrieveAsStranger.String()},
				},
			},
			{
				From: map[string][]string{
					"cn": {"sg_1"},
				},
				To: ldap.MappingTo{
					Groups:      []string{"super group"},
					Permissions: []string{charon.UserCanCreate.String()},
				},
			},
			{
				From: map[string][]string{
					"memberOf": {"cn=project1", "cn=company"},
				},
				To: ldap.MappingTo{
					IsStaff: true,
					Groups:  []string{"project-manager"},
				},
			},
		},
	}

	cases := map[string]struct {
		isStaff             bool
		groups, permissions []string
		attributes          []*libldap.EntryAttribute
		ok                  bool
	}{
		"empty": {
			attributes: []*libldap.EntryAttribute{},
		},
		"less": {
			attributes: []*libldap.EntryAttribute{
				{Name: "cn", Values: []string{"cn_1"}},
				{Name: "dc", Values: []string{"dc_1"}},
				{Name: "ou", Values: []string{"ou_1"}},
			},
		},
		"exact": {
			groups:      []string{"admins"},
			permissions: []string{charon.UserCanCreate.String()},
			attributes: []*libldap.EntryAttribute{
				{Name: "cn", Values: []string{"cn_1"}},
				{Name: "dc", Values: []string{"dc_1"}},
				{Name: "ou", Values: []string{"ou_1", "ou_2"}},
			},
			ok: true,
		},
		"exact-joined": {
			groups:      []string{"admins"},
			permissions: []string{charon.UserCanCreate.String()},
			attributes: []*libldap.EntryAttribute{
				{Name: "cn", Values: []string{"cn_1"}},
				{Name: "dc", Values: []string{"dc_1"}},
				{Name: "ou", Values: []string{"ou_1,ou_2"}},
			},
			ok: true,
		},
		"both": {
			groups:      []string{"admins", "users"},
			permissions: []string{charon.UserCanCreate.String(), charon.UserCanRetrieveAsStranger.String()},
			attributes: []*libldap.EntryAttribute{
				{Name: "cn", Values: []string{"cn_1", "cn_4"}},
				{Name: "dc", Values: []string{"dc_1"}},
				{Name: "ou", Values: []string{"ou_1", "ou_2"}},
			},
			ok: true,
		},
		"compound": {
			groups:      []string{"admins"},
			permissions: []string{charon.UserCanCreate.String()},
			attributes: []*libldap.EntryAttribute{
				{Name: "cn", Values: []string{"cn_1"}},
				{Name: "dc", Values: []string{"dc_1"}},
				{Name: "ou", Values: []string{"ou_1,ou_2"}},
			},
			ok: true,
		},
		"whitespace": {
			groups:      []string{"super group"},
			permissions: []string{charon.UserCanCreate.String()},
			attributes: []*libldap.EntryAttribute{
				{Name: "cn", Values: []string{"sg_1"}},
			},
			ok: true,
		},
		"staff": {
			isStaff: true,
			groups:  []string{"project-manager"},
			attributes: []*libldap.EntryAttribute{
				{Name: "memberOf", Values: []string{"cn=project1,cn=company", "cn=project2,cn=company"}},
			},
			ok: true,
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			got, ok := mappings.Map(c.attributes)
			if c.ok && c.ok != ok {
				t.Fatal("expected mappings")
			}
			if !c.ok && c.ok != ok {
				t.Fatal("unexpected mappings")
			}

			if !reflect.DeepEqual(c.groups, got.Groups) {
				t.Errorf("groups do not match, expected %v but got %v", c.groups, got.Groups)
			}
			if !reflect.DeepEqual(c.permissions, got.Permissions) {
				t.Errorf("permissions do not match, expected %v but got %v", c.permissions, got.Permissions)
			}
			if c.isStaff != got.IsStaff {
				t.Errorf("is staff do nto match, expected %t but got %t", c.isStaff, got.IsStaff)
			}
		})
	}
}
