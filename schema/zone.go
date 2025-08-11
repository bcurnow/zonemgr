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

// Represents the overall Zone file structure, the YAML file is an array of these
type Zone struct {
	Config                *Config                    `yaml:"config"`
	ResourceRecords       map[string]*ResourceRecord `yaml:"resource_records"`
	TTL                   *TTL                       `yaml:"ttl"`
	resourceRecordsByType map[ResourceRecordType]map[string]*ResourceRecord
}

func (z *Zone) SOARecord() *ResourceRecord {
	for _, rr := range z.ResourceRecords {
		if rr.Type == string(SOA) {
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
