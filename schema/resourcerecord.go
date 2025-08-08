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

	"github.com/bcurnow/zonemgr/plugins/proto"
)

// A generic type that can represent a variety of records types as many follow this specific format (A, CNAME, etc.	)
type ResourceRecord struct {
	Name    string                `yaml:"name"`
	Type    string                `yaml:"type"`  //TODO see if we can use something similar to ResourceRecordClass instead, this would simplify validations
	Class   string                `yaml:"class"` //TODO See if we can use ResourceRecordClass instead, this would simplify validations
	TTL     *int32                `yaml:"ttl"`
	Values  []ResourceRecordValue `yaml:"values"`
	Value   string                `yaml:"value"`
	Comment string                `yaml:"comment"`
}

func (rr ResourceRecord) ToProtoBuf() *proto.ResourceRecord {
	var ttl int32 = -1
	if rr.TTL != nil {
		// We're using a negative number so we can check for it the other way as well and set appropriately
		ttl = *rr.TTL
	}
	ret := &proto.ResourceRecord{
		Name:    rr.Name,
		Type:    rr.Type,
		Class:   rr.Class,
		Ttl:     ttl,
		Value:   rr.Value,
		Values:  rr.toProtoBufValues(),
		Comment: rr.Comment,
	}

	return ret
}

func (rr ResourceRecord) FromProtoBuf(p *proto.ResourceRecord) ResourceRecord {
	var ttl *int32 = nil
	if p.Ttl != -1 {
		ttl = &p.Ttl
	}
	return ResourceRecord{
		Name:    p.Name,
		Type:    p.Type,
		Class:   p.Class,
		TTL:     ttl,
		Value:   p.Value,
		Values:  rr.fromProtoBufValues(p.Values),
		Comment: p.Comment,
	}
}

func (rr ResourceRecord) toProtoBufValues() []*proto.ResourceRecordValue {
	fmt.Println("In toProtoBufValues")
	protoValues := make([]*proto.ResourceRecordValue, len(rr.Values))
	for i, inputValue := range rr.Values {
		protoValues[i] = &proto.ResourceRecordValue{Value: inputValue.Value, Comment: inputValue.Comment}
	}
	fmt.Println("Leaving toProtoBufValues")
	return protoValues
}

func (rr ResourceRecord) fromProtoBufValues(p []*proto.ResourceRecordValue) []ResourceRecordValue {
	values := make([]ResourceRecordValue, len(p))
	for i, value := range p {
		values[i] = ResourceRecordValue{Value: value.Value, Comment: value.Comment}
	}
	return values
}

type ResourceRecordValue struct {
	Value   string `yaml:"value"`
	Comment string `yaml:"comment"`
}

func (rrv ResourceRecordValue) ToProtoBuf() *proto.ResourceRecordValue {
	return &proto.ResourceRecordValue{Value: rrv.Value, Comment: rrv.Comment}
}

func (rrv ResourceRecordValue) FromProtoBuf(p *proto.ResourceRecordValue) ResourceRecordValue {
	return ResourceRecordValue{Value: p.Value, Comment: p.Comment}
}

// Defines the types of classes available in a zone file
type ResourceRecordClass string

const (
	INTERNET ResourceRecordClass = "IN"
	CSNET    ResourceRecordClass = "CS"
	CHAOS    ResourceRecordClass = "CH"
	HESIOD   ResourceRecordClass = "HS"
)

var resourceRecordToString = map[string]ResourceRecordClass{
	"IN": INTERNET,
	"CS": CSNET,
	"CH": CHAOS,
	"HS": HESIOD,
}

func (rrc ResourceRecordClass) IsValid() bool {
	switch rrc {
	case INTERNET, CSNET, CHAOS, HESIOD:
		return true
	default:
		return false
	}
}

func ResourceRecordClassFromString(str string) (ResourceRecordClass, error) {
	class, ok := resourceRecordToString[str]
	if !ok {
		return "", fmt.Errorf("Invalid resource record class '%s'\n", str)
	}
	return class, nil
}
