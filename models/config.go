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

import "fmt"

type Config struct {
	GenerateSerial bool `yaml:"generate_serial" validate:"boolean"`
	// Not validated as a dirpath: FileSerialManager.Next() already creates this directory itself via
	// MkdirAll, so requiring it to pre-exist is unnecessary and breaks first-time setup.
	SerialChangeIndexDirectory string `yaml:"serial_change_index_directory" validate:"omitempty"`
	GenerateReverseLookupZones bool   `yaml:"generate_reverse_lookup_zones" validate:"boolean"`
	IsCatalog                  bool   `yaml:"is_catalog" validate:"boolean"`
	CatalogIncludeReverseZones bool   `yaml:"catalog_include_reverse_zones" validate:"boolean"`
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{ GenerateSerial: %t, GenerateReverseLookupZones: %t, SerialChangeIndexDirectory: %s, IsCatalog: %t, CatalogIncludeReverseZones: %t }", c.GenerateSerial, c.GenerateReverseLookupZones, c.SerialChangeIndexDirectory, c.IsCatalog, c.CatalogIncludeReverseZones)
}
