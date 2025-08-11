/*
Copyright © 2025 Brian Curnow

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
omitempty
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package schema

import (
	"fmt"
)

// A generic type that can represent a variety of records types as many follow this specific format (A, CNAME, etc.	)
type ResourceRecord struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`  //TODO see if we can use something similar to ResourceRecordClass instead, this would simplify validations
	Class   string                 `yaml:"class"` //TODO See if we can use ResourceRecordClass instead, this would simplify validations
	TTL     *int32                 `yaml:"ttl"`
	Values  []*ResourceRecordValue `yaml:"values"`
	Value   string                 `yaml:"value"`
	Comment string                 `yaml:"comment"`
}

type ResourceRecordValue struct {
	Value   string `yaml:"value"`
	Comment string `yaml:"comment"`
}

// Defines the types of classes available in a zone file
type ResourceRecordClass string

const (
	INTERNET ResourceRecordClass = "IN"
	CSNET    ResourceRecordClass = "CS"
	CHAOS    ResourceRecordClass = "CH"
	HESIOD   ResourceRecordClass = "HS"
)

var resourceRecordClassToString = map[string]ResourceRecordClass{
	string(INTERNET): INTERNET,
	string(CSNET):    CSNET,
	string(CHAOS):    CHAOS,
	string(HESIOD):   HESIOD,
}

func (rrc ResourceRecordClass) IsValid() bool {
	switch rrc {
	case INTERNET, CSNET, CHAOS, HESIOD:
		return true
	default:
		return false
	}
}

func ResourceRecordClassFromString(str string) (*ResourceRecordClass, error) {
	class, ok := resourceRecordClassToString[str]
	if !ok {
		return nil, fmt.Errorf("Invalid resource record class '%s'\n", str)
	}
	return &class, nil
}
