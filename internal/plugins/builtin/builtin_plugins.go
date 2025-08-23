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

package builtin

import "github.com/bcurnow/zonemgr/plugins"

var builtins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
var metadata = make(map[plugins.PluginType]*plugins.PluginMetadata)

func BuiltinPlugins() map[plugins.PluginType]plugins.ZoneMgrPlugin {
	return builtins
}

func BuiltinMetadata() map[plugins.PluginType]*plugins.PluginMetadata {
	return metadata
}

func registerBuiltIn(pluginType plugins.PluginType, plugin plugins.ZoneMgrPlugin) {
	builtins[pluginType] = plugin
	metadata[pluginType] = &plugins.PluginMetadata{Name: string(pluginType), Command: "Built In", BuiltIn: true}
}
