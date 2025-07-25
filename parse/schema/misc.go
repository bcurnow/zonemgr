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
package schema

import "sort"

// A generic type to hold an integer value and a comment as many of the Zone fields allow for this
type ZoneValue struct {
	Value   int64  `yaml:"value"`
	Comment string `yaml:"comment"`
}

// Defines the order in which the records will be sorted in the zonefile
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
