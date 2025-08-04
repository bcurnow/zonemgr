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

package manager

import (
	"maps"

	"github.com/bcurnow/zonemgr/logging"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/builtin"
)

var (
	logger      = logging.Logger().Named("plugin-manager")
	allPlugins  = make(map[plugins.PluginType]plugins.TypeHandler)
	initialized = false
)

func Plugins() (map[plugins.PluginType]plugins.TypeHandler, error) {
	if err := initializePlugins(); err != nil {
		return nil, err
	}

	return allPlugins, nil
}

func initializePlugins() error {
	if initialized {
		return nil
	}

	maps.Copy(allPlugins, builtin.BuiltinPlugins())
	externalPlugins, err := plugins.RegisterPlugins()
	if err != nil {
		return err
	}

	// We could just copy the externalPlugins map ove the allPlugins map and everything would be fine
	// However, iterating gives us better diagnostic logs
	for name, externalPlugin := range externalPlugins {
		// Get the list of resource record types the external plugin supports
		supportedTypes := externalPlugin.PluginTypesSupported()

		for _, rrType := range supportedTypes {
			// Check to see if we already have a plugin for this ResourceRecord Type
			_, ok := allPlugins[rrType]
			if ok {
				logger.Debug("Replacing existing plugin", "resourceRecordType", rrType, "newPluginName", name)
				allPlugins[rrType] = externalPlugin
			}
		}
	}
	initialized = true
	return nil
}
