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
	"errors"
	"sort"
)

func PluginTypes(pluginTypes ...PluginType) []PluginType {
	return pluginTypes
}

func WithSortedPlugins(p map[PluginType]ZoneMgrPlugin, pluginMetadata map[PluginType]*PluginMetadata, fn func(pluginType PluginType, p ZoneMgrPlugin, metadata *PluginMetadata) error) error {
	for _, pluginType := range sortedPluginKeys(p) {
		metadata, ok := pluginMetadata[pluginType]
		if !ok {
			return errors.New("could not find plugin metadata for plugin type: " + string(pluginType))
		}
		if err := fn(pluginType, p[pluginType], metadata); err != nil {
			return err
		}
	}
	return nil
}

func sortedPluginKeys(p map[PluginType]ZoneMgrPlugin) []PluginType {
	keys := make([]PluginType, 0, len(p))
	for pluginType := range p {
		keys = append(keys, pluginType)
	}

	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})

	return keys
}
