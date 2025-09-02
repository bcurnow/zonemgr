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

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package models

import (
	"fmt"
	"sort"
	"strings"
)

// Represents the overall Zone file structure, the YAML file is an array of these
type Zone struct {
	Config                *Config                    `yaml:"config"`
	ResourceRecords       map[string]*ResourceRecord `yaml:"resource_records"`
	TTL                   *TTL                       `yaml:"ttl"`
	resourceRecordsByType map[ResourceRecordType]map[string]*ResourceRecord
}

func (z *Zone) String() string {
	var rrString strings.Builder

	for identifier, rr := range z.ResourceRecords {
		rrString.WriteString("     ")
		rrString.WriteString(identifier)
		rrString.WriteString(" -> ")
		rrString.WriteString(rr.String())
		rrString.WriteString("\n")
	}

	return "Zone{\n" +
		fmt.Sprintf("   Config: %s\n", z.Config) +
		"   ResourceRecords:\n" +
		fmt.Sprintf("%s\n", rrString.String()[:len(rrString.String())-1]) +
		fmt.Sprintf("   TTL: %s\n", z.TTL) +
		"}"
}

func (z *Zone) SOARecord() *ResourceRecord {
	for _, rr := range z.ResourceRecords {
		if rr.Type == SOA {
			return rr
		}
	}

	return nil
}

func (z *Zone) ResourceRecordsByType() map[ResourceRecordType]map[string]*ResourceRecord {
	if nil == z.resourceRecordsByType {
		z.resourceRecordsByType = make(map[ResourceRecordType]map[string]*ResourceRecord, len(z.ResourceRecords))
	}

	for identifier, rr := range z.ResourceRecords {
		rrType := ResourceRecordType(rr.Type)
		_, ok := z.resourceRecordsByType[rrType]
		if !ok {
			rrsOfType := make(map[string]*ResourceRecord)
			z.resourceRecordsByType[rrType] = rrsOfType
		}

		z.resourceRecordsByType[rrType][identifier] = rr
	}

	return z.resourceRecordsByType
}

func (z *Zone) WithSortedResourceRecords(fn func(identifier string, rr *ResourceRecord) error) error {
	for _, identifier := range z.sortedResourceRecordKeys() {
		if err := fn(identifier, z.ResourceRecords[identifier]); err != nil {
			return err
		}
	}
	return nil
}

func (z *Zone) sortedResourceRecordKeys() []string {
	keys := make([]string, 0, len(z.ResourceRecords))
	for k := range z.ResourceRecords {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
