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

import "github.com/bcurnow/zonemgr/plugins/proto"

type Config struct {
	GenerateSerial             bool   `yaml:"generate_serial"`
	SerialChangeIndex          uint32 `yaml:"serial_change_index"`
	GenerateReverseLookupZones bool   `yaml:"generate_reverse_lookup_zones"`
}

func (c Config) ToProtoBuf() *proto.Config {
	return &proto.Config{GenerateSerial: c.GenerateSerial, SerialChangeIndex: c.SerialChangeIndex, GenerateReverseLookupZones: c.GenerateReverseLookupZones}
}

func (c Config) FromProtoBuf(p *proto.Config) Config {
	return Config{GenerateSerial: p.GenerateSerial, SerialChangeIndex: p.SerialChangeIndex, GenerateReverseLookupZones: p.GenerateReverseLookupZones}
}
