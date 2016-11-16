package ldap

import (
	"encoding/json"
	"os"

	"github.com/go-ldap/ldap"
)

// Mappings ...
type Mappings []Mapping

// Mapping ...
type Mapping struct {
	From MappingFrom `json:"from"`
	To   MappingTo   `json:"to"`
}

// MappingFrom ...
type MappingFrom struct {
	CommonNames         []string `json:"cn"`
	OrganizationalUnits []string `json:"ou"`
	DomainComponents    []string `json:"dc"`
}

// MappingTo ...
type MappingTo struct {
	Groups      []string `json:"groups"`
	Permissions []string `json:"permissions"`
}

// NewMappingsFromFile reads json file and allocates new Mappings based on the output.
func NewMappingsFromFile(p string) (Mappings, error) {
	var m Mappings
	if p == "" {
		return m, nil
	}

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(file).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// Map search groups and permissions that given LDAP entry match.
func (m Mappings) Map(attrs []*ldap.EntryAttribute) ([]string, []string, bool) {
	var (
		groups, permissions []string
		expected, valid     int
	)

MappingLoop:
	for _, mapping := range m {
		expected, valid = 0, 0
		if len(mapping.From.CommonNames) > 0 {
			expected++
		}
		if len(mapping.From.OrganizationalUnits) > 0 {
			expected++
		}
		if len(mapping.From.DomainComponents) > 0 {
			expected++
		}

	AttributesLoop:
		for _, attr := range attrs {
			if valid >= expected {
				break AttributesLoop
			}
			switch attr.Name {
			case "cn":
				if !m.compare(attr.Values, mapping.From.CommonNames) {
					continue MappingLoop
				}
				valid++
			case "ou":
				if !m.compare(attr.Values, mapping.From.OrganizationalUnits) {
					continue MappingLoop
				}
				valid++
			case "dc":
				if !m.compare(attr.Values, mapping.From.DomainComponents) {
					continue MappingLoop
				}
				valid++
			}
		}

		if valid >= expected {
			groups = append(groups, mapping.To.Groups...)
			permissions = append(permissions, mapping.To.Permissions...)
		}
	}

	return groups, permissions, len(groups) > 0 || len(permissions) > 0
}

func (lm Mappings) compare(givens, expexteds []string) bool {
	if len(expexteds) == 0 {
		return true
	}
	var match int
	for _, given := range givens {
		for _, expected := range expexteds {
			if given == expected {
				match++
			}
		}
	}

	return match == len(expexteds)
}
