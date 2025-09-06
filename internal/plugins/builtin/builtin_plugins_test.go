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

import (
	"testing"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBuiltinPlugins(t *testing.T) {

	builtins := BuiltinPlugins()

	if len(builtins) != 5 {
		t.Errorf("expected all 5 builtin plugins to have registered")
	}

	testCases := []struct {
		pluginType        plugins.PluginType
		expectedInterface interface{}
	}{
		{pluginType: plugins.A, expectedInterface: &BuiltinPluginA{}},
		{pluginType: plugins.CNAME, expectedInterface: &BuiltinPluginCNAME{}},
		{pluginType: plugins.NS, expectedInterface: &BuiltinPluginNS{}},
		{pluginType: plugins.PTR, expectedInterface: &BuiltinPluginPTR{}},
		{pluginType: plugins.SOA, expectedInterface: &BuiltinPluginSOA{}},
	}

	for _, tc := range testCases {
		p, ok := builtins[tc.pluginType]
		if !ok {
			t.Errorf("expected to find plugin of type %s", tc.pluginType)
		} else {
			if !cmp.Equal(p, tc.expectedInterface, cmpopts.IgnoreUnexported(BuiltinPluginSOA{})) {
				t.Errorf("expected plugin of type %s to implement %T, but was %T instead", tc.pluginType, tc.expectedInterface, p)
			}

		}
	}
}

func TestBuiltinMetadata(t *testing.T) {

	metadata := BuiltinMetadata()

	if len(metadata) != 5 {
		t.Errorf("expected all 5 builtin plugins to have registered")
	}

	testCases := []struct {
		pluginType       plugins.PluginType
		expectedMetadata *plugins.Metadata
	}{
		{pluginType: plugins.A, expectedMetadata: &plugins.Metadata{Name: string(plugins.A), Command: "Built In", BuiltIn: true}},
		{pluginType: plugins.CNAME, expectedMetadata: &plugins.Metadata{Name: string(plugins.CNAME), Command: "Built In", BuiltIn: true}},
		{pluginType: plugins.NS, expectedMetadata: &plugins.Metadata{Name: string(plugins.NS), Command: "Built In", BuiltIn: true}},
		{pluginType: plugins.PTR, expectedMetadata: &plugins.Metadata{Name: string(plugins.PTR), Command: "Built In", BuiltIn: true}},
		{pluginType: plugins.SOA, expectedMetadata: &plugins.Metadata{Name: string(plugins.SOA), Command: "Built In", BuiltIn: true}},
	}

	for _, tc := range testCases {
		m, ok := metadata[tc.pluginType]
		if !ok {
			t.Errorf("expected to find metadata of type %s", tc.pluginType)
		} else {
			if !cmp.Equal(m, tc.expectedMetadata) {
				t.Errorf("incorrect metadata for plugin of type %s: %v, want %v", tc.pluginType, m, tc.expectedMetadata)
			}

		}
	}
}
