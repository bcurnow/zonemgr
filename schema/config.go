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

import (
	"strconv"

	"github.com/bcurnow/zonemgr/env"
)

type Config struct {
	PluginsDirectory           string `yaml:"plugins_directory"`
	GenerateSerial             bool   `yaml:"generate_serial"`
	SerialChangeIndexDirectory string `yaml:"serial_change_index_directory"`
	GenerateReverseLookupZones bool   `yaml:"generate_reverse_lookup_zones"`
}

func (c *Config) ConfigDefaults() error {
	c.PluginsDirectory = env.PluginsDirectory.Value
	val, err := strconv.ParseBool(env.GenerateSerial.Value)
	if err != nil {
		return err
	}
	c.GenerateSerial = val

	c.SerialChangeIndexDirectory = env.SerialChangeIndexDirectory.Value
	val, err = strconv.ParseBool(env.GenerateReverseLookupZones.Value)
	if err != nil {
		return err
	}
	c.GenerateReverseLookupZones = val
	return nil
}
