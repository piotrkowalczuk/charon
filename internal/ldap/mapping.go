package ldap

import (
	"encoding/json"
	"os"
	"strings"

	"io"

	"sort"

	"github.com/go-ldap/ldap"
)

// Mappings ...
type Mappings struct {
	Mappings []Mapping
	// Attributes is a list of all unique attributes that appear in mappings
	// that needs to be retrieved from server to make comparison.
	Attributes []string
}

// Mapping ...
type Mapping struct {
	From       map[string][]string `json:"from"`
	To         MappingTo           `json:"to"`
	attributes []string
}

// MappingTo ...
type MappingTo struct {
	IsStaff     bool     `json:"isStaff"`
	Groups      []string `json:"groups"`
	Permissions []string `json:"permissions"`
}

// NewMappings ...
func NewMappings(r io.Reader) (*Mappings, error) {
	var m Mappings
	if err := json.NewDecoder(r).Decode(&m.Mappings); err != nil {
		return nil, err
	}

	for _, mapping := range m.Mappings {
		for attr := range mapping.From {
			m.Attributes = append(m.Attributes, attr)
		}
	}

	removeDuplicates(&m.Attributes)
	sort.Sort(sort.StringSlice(m.Attributes))

	return &m, nil
}

// NewMappingsFromFile reads json file and allocates new Mappings based on the output.
func NewMappingsFromFile(p string) (*Mappings, error) {
	if p == "" {
		return &Mappings{}, nil
	}

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return NewMappings(file)
}

// Map search groups and permissions that given LDAP entry match.
func (m *Mappings) Map(attrs []*ldap.EntryAttribute) (*MappingTo, bool) {
	var (
		to              MappingTo
		expected, valid int
		found           bool
	)

MappingLoop:
	for _, mapping := range m.Mappings {
		expected, valid = 0, 0
		for _, values := range mapping.From {
			if len(values) > 0 {
				expected++
			}
		}

	AttributesLoop:
		for _, attr := range attrs {
			if valid >= expected {
				break AttributesLoop
			}

			for attrName, from := range mapping.From {
				if attr.Name != attrName {
					continue
				}
				if !m.compare(attr.Values, from) {
					continue MappingLoop
				}

				valid++
			}
		}

		if valid >= expected {
			found = true
			if mapping.To.IsStaff {
				to.IsStaff = true
			}
			to.Groups = append(to.Groups, mapping.To.Groups...)
			to.Permissions = append(to.Permissions, mapping.To.Permissions...)
		}
	}

	return &to, found
}

func (lm *Mappings) compare(given, expected []string) bool {
	if len(expected) == 0 {
		return true
	}
	var (
		match, localMatch int
	)
	for _, giv := range given {
		parts := strings.Split(giv, ",")

		// if longer than 1 than probably its memberOf type of attribute
		switch len(parts) {
		case 1:
			for _, exp := range expected {
				if giv == exp {
					match++
				}
			}
		default:
			for _, part := range parts {
				for _, exp := range expected {
					if part == exp {
						localMatch++
					}
				}
			}
			if localMatch == len(expected) {
				return true
			}
			localMatch = 0
		}
	}

	return match == len(expected)
}

func removeDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
