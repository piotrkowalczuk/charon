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
		},
	}

	cases := map[string]struct {
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
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			groups, permissions, ok := mappings.Map(c.attributes)
			if c.ok && c.ok != ok {
				t.Fatal("expected mappings")
			}
			if !c.ok && c.ok != ok {
				t.Fatal("unexpected mappings")
			}

			if !reflect.DeepEqual(c.groups, groups) {
				t.Errorf("groups do not match, expected %v but got %v", c.groups, groups)
			}
			if !reflect.DeepEqual(c.permissions, permissions) {
				t.Errorf("permissions do not match, expected %v but got %v", c.permissions, permissions)
			}
		})
	}
}
