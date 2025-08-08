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

var builtins map[plugins.PluginType]*plugins.TypeHandlerPlugin = make(map[plugins.PluginType]*plugins.TypeHandlerPlugin)

func BuiltinPlugins() map[plugins.PluginType]*plugins.TypeHandlerPlugin {
	return builtins
}

func registerBuiltIn(pluginType plugins.PluginType, plugin plugins.TypeHandler) {
	builtins[pluginType] = &plugins.TypeHandlerPlugin{IsBuiltIn: true, PluginName: string(pluginType), Plugin: plugin, PluginCmd: "Built In"}
}
