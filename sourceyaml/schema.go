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
package sourceyaml

import "sort"

// Defines the structs which are necessary to parse the YAML. For an example of the YAML format, see examples/zones.yaml.:

// Represents the overall Zone file structure, the YAML file is an array of these
type Zone struct {
	TTL                        ZoneValue                 `yaml:"ttl"`
	Class                      string                    `yaml:"class,omitempty"`
	Nameserver                 string                    `yaml:"nameserver"`
	Administrator              string                    `yaml:"administrator"`
	Refresh                    ZoneValue                 `yaml:"refresh"`
	Retry                      ZoneValue                 `yaml:"retry"`
	Expire                     ZoneValue                 `yaml:"expire"`
	Minimum                    ZoneValue                 `yaml:"minimum"`
	Serial                     uint32                    `yaml:"serial,omitempty"`
	GenerateSerial             bool                      `yaml:"generate_serial,omitempty"`
	SerialChangeIndex          uint32                    `yaml:"serial_change_index"`
	GenerateReverseLookupZones bool                      `yaml:"generate_reverse_lookup_zones,omitempty"`
	ResourceRecords            map[string]ResourceRecord `yaml:"resource_records,omitempty"`
}

// A generic type to hold an integer value and a comment as many of the Zone fields allow for this
type ZoneValue struct {
	Value   int64  `yaml:"value"`
	Comment string `yaml:"comment"`
}

// A generic type that can represent a variety of records types as many follow this specific format (A, CNAME, etc.	)
type ResourceRecord struct {
	Type    string `yaml:"type"`
	Class   string `yaml:"class,omitempty"`
	Value   string `yaml:"value"`
	TTL     int64  `yaml:"ttl,omitempty"`
	Comment string `yaml:"comment,omitempty"`
}

var resourceRecordTypeSortOrder = []string{
	"NS", "A", "CNAME",
}

// Returns the keys of the ResourceRecords in a specific order, in this case, the order specified in resourceRecordTypeSortOrder, each type is then sorted alphabetically
// NOTE: This is based on my personal preference on to sort the records, this should really be more of a strategy pattern so that other options (e.g. strict alphabetical, nosort/map order) can be used
func (z Zone) ResourceRecordsSortOrder() []string {
	// Returns the keys of the ResourceRecords map in sorted order
	keys := make([]string, 0, len(z.ResourceRecords))
	for _, recordType := range resourceRecordTypeSortOrder {
		keys = append(keys, getKeysByType(z.ResourceRecords, recordType)...)
	}
	return keys
}

func getKeysByType(resourceRecords map[string]ResourceRecord, recordType string) []string {
	// Returns the keys of the ResourceRecords map that match the specified type
	keys := make([]string, 0)
	for k, v := range resourceRecords {
		if v.Type == recordType {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}
