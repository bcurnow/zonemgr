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
	"strconv"
	"strings"
)

// A generic type that can represent a variety of records types as many follow this specific format (A, CNAME, etc.	)
type ResourceRecord struct {
	Name    string                 `yaml:"name"`
	Type    ResourceRecordType     `yaml:"type"`  //TODO see if we can use something similar to ResourceRecordClass instead, this would simplify validations
	Class   ResourceRecordClass    `yaml:"class"` //TODO See if we can use ResourceRecordClass instead, this would simplify validations
	TTL     *int32                 `yaml:"ttl"`
	Values  []*ResourceRecordValue `yaml:"values"`
	Value   string                 `yaml:"value"`
	Comment string                 `yaml:"comment"`
}

// There are two possible places to get a value from: Value or Values[0].Value
// This method will validate that only Value or Values is populated, that, if Values is populated, there's only a single item.
// Will return either Value or the Values[0].Value
func (rr *ResourceRecord) RetrieveSingleValue(identifier string) (string, error) {
	if err := rr.IsValueSetInOnePlace(identifier); err != nil {
		return "", err
	}

	if err := rr.hasSingleValue(identifier); err != nil {
		return "", err
	}

	if len(rr.Values) == 0 {
		return rr.Value, nil
	}

	//Only option left is the first value in Values
	return rr.Values[0].Value, nil
}

// There are two possible places for a comment to be: Comment or Values[0].Comment
// This method will validate that only Comment or Values is populated, that, if Values is populated, there's only a single item
// Will return either Comment or Values[0].Comment
func (rr *ResourceRecord) RetrieveSingleComment(identifier string) (string, error) {
	if err := rr.IsCommentSetInOnePlace(identifier); err != nil {
		return "", err
	}

	if err := rr.hasSingleValue(identifier); err != nil {
		return "", err
	}

	if len(rr.Values) == 0 {
		return rr.Comment, nil
	}

	//Only option left is the first comment in Values
	return rr.Values[0].Comment, nil
}

// Validates that either Values has more than one element or Value is set, not both
// Allows for Value to be blank and does not check Values[*].Value at all
func (rr *ResourceRecord) IsValueSetInOnePlace(identifier string) error {
	if len(rr.Values) > 0 && rr.Value != "" {
		return fmt.Errorf("%s record invalid, can not specify both value and values, identifier: '%s'", rr.Type, identifier)
	}
	return nil
}

// Validates that either Values has more than one element or Comment is set, not both
// Allows for Comment to be blank and does not check Values[*].Comment at all
func (rr *ResourceRecord) IsCommentSetInOnePlace(identifier string) error {
	if len(rr.Values) > 0 && rr.Comment != "" {
		return fmt.Errorf("%s record invalid, can not specify both comment and values, identifier: '%s'", rr.Type, identifier)
	}
	return nil
}

func (rr *ResourceRecord) RenderResourceWithoutValue() string {
	var record strings.Builder

	record.WriteString(fmt.Sprintf(ResourceRecordNameFormatString, rr.Name))
	record.WriteString(fmt.Sprintf(ResourceRecordTypeFormatString, rr.Type))
	if rr.Class != "" {
		record.WriteString(" ")
		record.WriteString(string(rr.Class))
		record.WriteString(" ")
	}

	if rr.TTL != nil {
		record.WriteString(" ")
		record.WriteString(strconv.Itoa(int(*rr.TTL)))
		record.WriteString(" ")
	}

	return record.String()
}

func (rr *ResourceRecord) RenderSingleValueResource() string {
	var record strings.Builder
	record.WriteString(rr.RenderResourceWithoutValue())
	record.WriteString(rr.Value)
	if rr.Comment != "" {
		record.WriteString(" ; ")
		record.WriteString(rr.Comment)
	}

	return record.String()
}

func (rr *ResourceRecord) RenderMultivalueResource() string {
	var record strings.Builder
	record.WriteString(rr.RenderResourceWithoutValue())
	record.WriteString("(\n")
	indentFormatString := "%" + strconv.Itoa(record.Len()-2) + "s"
	for _, value := range rr.Values {
		record.WriteString(fmt.Sprintf(indentFormatString, ""))                         // This will ensure that all the values are indented
		record.WriteString(fmt.Sprintf(ResourceRecordMultivalueIndentFormatString, "")) // This will add an indent inside the parens
		record.WriteString(fmt.Sprintf(ResourceRecordNameFormatString, value.Value))
		if value.Comment != "" {
			record.WriteString(" ; ")
			record.WriteString(value.Comment)
		}
		record.WriteString("\n")
	}
	record.WriteString(fmt.Sprintf(indentFormatString, ""))
	record.WriteString(")")

	return record.String()
}

func (rr *ResourceRecord) hasSingleValue(identifier string) error {
	if len(rr.Values) <= 1 {
		return nil
	}
	return fmt.Errorf("%s record invalid, found more than one value in values element, identifier: '%s'", rr.Type, identifier)
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

func (rrc ResourceRecordClass) IsValid() bool {
	switch rrc {

	case INTERNET, CSNET, CHAOS, HESIOD, "": //It's always valid for the class to be empty
		return true
	default:
		return false
	}
}

func ResourceRecordClassFromString(str string) (*ResourceRecordClass, error) {
	class, ok := resourceRecordClassToString[str]
	if !ok {
		return nil, fmt.Errorf("invalid resource record class '%s'", str)
	}
	return &class, nil
}

var resourceRecordClassToString = map[string]ResourceRecordClass{
	string(INTERNET): INTERNET,
	string(CSNET):    CSNET,
	string(CHAOS):    CHAOS,
	string(HESIOD):   HESIOD,
}
