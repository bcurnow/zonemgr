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
