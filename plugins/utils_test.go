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
	"testing"

	"github.com/golang/mock/gomock"
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

func TestWithSortedPlugins(t *testing.T) {
	testCases := []struct {
		missingMetadata bool
		functionErr     bool
	}{
		{},
		{missingMetadata: true},
		{functionErr: true},
	}

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mockZoneMgrPlugin := NewMockZoneMgrPlugin(mockController)
	mockZoneMgrPluginMetadata := &Metadata{Name: "mock", Command: "testing", BuiltIn: false}

	for _, tc := range testCases {
		mockPlugins := make(map[PluginType]ZoneMgrPlugin)
		mockPlugins[A] = mockZoneMgrPlugin
		mockPlugins[CNAME] = mockZoneMgrPlugin
		mockMetadata := make(map[PluginType]*Metadata)
		mockMetadata[A] = mockZoneMgrPluginMetadata
		mockMetadata[CNAME] = mockZoneMgrPluginMetadata

		first := true
		testFn := func(pluginType PluginType, p ZoneMgrPlugin, metadata *Metadata) error {
			if first {
				first = false
				if pluginType != A {
					t.Errorf("incorrect first plugin type: '%s', want: '%s'", pluginType, A)
				}
			}
			return nil
		}

		if tc.functionErr {
			testFn = func(pluginType PluginType, p ZoneMgrPlugin, metadata *Metadata) error {
				return errors.New("functionErr")
			}
		}

		if tc.missingMetadata {
			mockMetadata = make(map[PluginType]*Metadata)
		}

		if err := WithSortedPlugins(mockPlugins, mockMetadata, testFn); err != nil {
			want := ""
			if tc.functionErr {
				want = "functionErr"
			} else if tc.missingMetadata {
				want = "could not find plugin metadata for plugin type: A"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.functionErr || tc.missingMetadata {
				t.Error("expected an error, found none")
			}
		}
	}
}
