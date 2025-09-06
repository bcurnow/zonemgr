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
package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
)

const (
	ResourceRecordNameFormatString             = "%-40s"
	ResourceRecordTypeFormatString             = "%-6s"
	ResourceRecordMultivalueIndentFormatString = "%4s"
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

func (rr *ResourceRecord) String() string {
	return "ResourceRecord{\n" +
		fmt.Sprintf("       Name: %s\n", rr.Name) +
		fmt.Sprintf("       Type: %s\n", rr.Type) +
		fmt.Sprintf("       Class: %s\n", rr.Class) +
		fmt.Sprintf("       TTL: %s\n", int32ToString(rr.TTL)) +
		fmt.Sprintf("       Values: %s\n", rr.Values) +
		fmt.Sprintf("       Value: %s\n", rr.Value) +
		fmt.Sprintf("       Comment: %s\n", rr.Comment) +
		"     }"
}

// There are two possible places to get a value from: Value or Values[0].Value
// This method will validate that only Value or Values is populated, that, if Values is populated, there's only a single item.
// Will return either Value or the Values[0].Value
func (rr *ResourceRecord) RetrieveSingleValue() string {
	valueCount := len(rr.Values)

	if valueCount > 0 {
		if valueCount > 1 {
			hclog.L().Trace("Resource record has more than 1 value, returning the first", "name", rr.Name, "valueCount", valueCount, "values", rr.Values)
		}
		return rr.Values[0].Value
	}

	// We don't have any values so, even if this is empty, return it
	return rr.Value
}

// There are two possible places for a comment to be: Comment or Values[0].Comment
// This method will validate that only Comment or Values is populated, that, if Values is populated, there's only a single item
// Will return either Comment or Values[0].Comment
func (rr *ResourceRecord) RetrieveSingleComment() string {
	valueCount := len(rr.Values)
	if valueCount > 0 && rr.Values[0].Comment != "" {
		if valueCount > 1 {
			hclog.L().Trace("Resource record has more than 1 value, returning the first", "name", rr.Name, "valueCount", valueCount, "values", rr.Values)
		}
		return rr.Values[0].Comment
	}

	// We don't have any values so, even if this is empty, return it
	return rr.Comment
}

// Validates that either Values has more than one element or Value is set, not both
// Allows for Value to be blank and does not check Values[*].Value at all
func (rr *ResourceRecord) IsValueSetInOnePlace() bool {
	// If we have no values, we can't possibly be set in more than one place
	if len(rr.Values) == 0 {
		return true

	}

	// If we do have values, make sure that the value is also empty
	if len(rr.Values) > 0 && rr.Value == "" {
		return true
	}

	return false
}

// Validates that either Values has more than one element or Comment is set, not both
// Allows for Comment to be blank and does not check Values[*].Comment at all
func (rr *ResourceRecord) IsCommentSetInOnePlace() bool {
	// If we have no values, we can't possibly be set in more than one place
	if len(rr.Values) == 0 {
		return true
	}

	// If we do have values, make sure that the comment is also empty
	if len(rr.Values) > 0 && rr.Comment == "" {
		return true
	}

	return false
}

func (rr *ResourceRecord) RenderResourceWithoutValue() string {
	var record strings.Builder

	record.WriteString(fmt.Sprintf(ResourceRecordNameFormatString, rr.Name))
	record.WriteString(" ")
	record.WriteString(fmt.Sprintf(ResourceRecordTypeFormatString, rr.Type))
	record.WriteString(" ")
	if rr.Class != "" {
		record.WriteString(string(rr.Class))
		record.WriteString(" ")
	}

	if rr.TTL != nil {
		record.WriteString(strconv.Itoa(int(*rr.TTL)))
		record.WriteString(" ")
	}

	return record.String()
}

func (rr *ResourceRecord) RenderSingleValueResource() string {
	var record strings.Builder
	record.WriteString(rr.RenderResourceWithoutValue())

	record.WriteString(rr.RetrieveSingleValue())

	if rr.RetrieveSingleComment() != "" {
		record.WriteString(" ;")
		record.WriteString(rr.RetrieveSingleComment())
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
			record.WriteString(" ;")
			record.WriteString(value.Comment)
		}
		record.WriteString("\n")
	}
	record.WriteString(fmt.Sprintf(indentFormatString, ""))
	record.WriteString(")")

	return record.String()
}
