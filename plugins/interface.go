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

package plugins

import (
	"github.com/bcurnow/zonemgr/schema"
)

type TypeHandler interface {
	// Returns the version of plugin
	PluginVersion() string
	// Returns the set of plugin types that this plugin supports
	PluginTypesSupported() []PluginType
	// Allows for configuration of the plugin, this will be called once for each zone in the file
	Configure(config schema.Config) error
	// Allows for validation and normalization/defaulting for the resource record
	Normalize(identifier string, rr schema.ResourceRecord) (schema.ResourceRecord, error)
	// Allows for validation of the entire normalized zone
	// This enables checks such as all CNAME records properly referencing a defined A record
	// This allows validation only, no defaulting
	ValidateZone(name string, zone schema.Zone) error
	// Converts the resource record into a string to be writting out to a file
	Render(identifier string, rr schema.ResourceRecord) (string, error)
}
