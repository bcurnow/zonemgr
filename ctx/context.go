/**
 * Copyright (C) 2025 Brian Curnow
 *
 * This file is part of zonemgr.
 *
 * zonemgr is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * zonemgr is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with zonemgr.  If not, see <https://www.gnu.org/licenses/>.
 */
package ctx

import (
	"github.com/spf13/pflag"
)

const (
	FlagGenerateReverseLookupZone  = "generate-reverse-lookup-zones"
	FlagGenerateSerial             = "generate-serial"
	FlagPluginsDirectory           = "plugins-dir"
	FlagSerialChangeIndexDirectory = "serial-change-index-directory"
)

// Holds all the common values that a plugin needs, provide access without errors
type PluginContext interface {
	GenerateReverseLookupZones() bool
	GenerateSerial() bool
	PluginsDirectory() string
	SerialChangeIndexDirectory() string
}

// The actual implementation, there will be only one of these instances
type pluginContext struct {
	generateReverseLookupZones bool
	generateSerial             bool
	pluginsDirectory           string
	serialChangeIndexDirectory string
}

// The single instance
var (
	pc  *pluginContext = &pluginContext{}
	env Environment    = &environment{}
)

func C() PluginContext {
	return pc
}

// Initializes the plugin context. This should be called once at the start of application setup to ensure a fully populated context
// This allows the errors to be handled up front in the command and all the rest of the code simply uses the values.
// The values will be by one of the following, these are in priority order, the first one we find becomes the value
//   - Command Line Flag
//   - Environment Variable
//   - Default Value
func InitPluginContext(flags *pflag.FlagSet) error {
	if err := env.LoadValues(); err != nil {
		return err
	}

	if flags.Changed(FlagGenerateReverseLookupZone) {
		generateReverseLookupZones, err := flags.GetBool(FlagGenerateReverseLookupZone)
		if err != nil {
			return err
		}
		pc.generateReverseLookupZones = generateReverseLookupZones
	} else {
		pc.generateReverseLookupZones = env.GenerateReverseLookupZones()
	}

	if flags.Changed(FlagGenerateSerial) {
		generateSerial, err := flags.GetBool(FlagGenerateSerial)
		if err != nil {
			return err
		}
		pc.generateSerial = generateSerial
	} else {
		pc.generateSerial = env.GenerateSerial()
	}

	if flags.Changed(FlagPluginsDirectory) {
		pluginsDirectory, err := flags.GetString(FlagPluginsDirectory)
		if err != nil {
			return err
		}
		pc.pluginsDirectory = pluginsDirectory
	} else {
		pc.pluginsDirectory = env.PluginsDirectory()
	}

	if flags.Changed(FlagSerialChangeIndexDirectory) {
		serialChangeIndexDirectory, err := flags.GetString(FlagSerialChangeIndexDirectory)
		if err != nil {
			return err
		}
		pc.serialChangeIndexDirectory = serialChangeIndexDirectory
	} else {
		pc.serialChangeIndexDirectory = env.SerialChangeIndexDirectory()
	}

	return nil
}

func (c *pluginContext) GenerateReverseLookupZones() bool {
	return c.generateReverseLookupZones
}

func (c *pluginContext) GenerateSerial() bool {
	return c.generateSerial
}

func (c *pluginContext) PluginsDirectory() string {
	return c.pluginsDirectory
}

func (c *pluginContext) SerialChangeIndexDirectory() string {
	return c.serialChangeIndexDirectory
}
