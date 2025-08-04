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

import goplugin "github.com/hashicorp/go-plugin"

// This is the goplugin.Plugin implementation
type Plugin struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl TypeHandler
}

// Defines which plugins we can dispense, since we only have one plugin type, here it is.
var PluginMap = map[string]goplugin.Plugin{
	"zonemgrplugin": &Plugin{},
}
