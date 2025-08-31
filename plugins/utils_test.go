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
	"testing"
)

func TestPluginTypes(t *testing.T) {
	pluginTypes := PluginTypes()

	if len(pluginTypes) != 0 {
		t.Errorf("incorrect number of plugin types: %d, want 0", len(pluginTypes))
	}

	pluginTypes = PluginTypes(A)
	if len(pluginTypes) != 1 {
		t.Errorf("incorrect number of plugin types: %d, want 1", len(pluginTypes))
	}

	pluginTypes = PluginTypes(A, CNAME, NS, PTR, SOA)
	if len(pluginTypes) != 5 {
		t.Errorf("incorrect number of plugin types: %d, want 5", len(pluginTypes))
	}
}
