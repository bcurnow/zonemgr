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
	"reflect"
	"testing"
	"unsafe"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/google/go-cmp/cmp"
)

var allPlugins = map[plugins.PluginType]plugins.ZoneMgrPlugin{
	plugins.A:     &BuiltinPluginA{},
	plugins.CNAME: &BuiltinPluginCNAME{},
	plugins.NS:    &BuiltinPluginNS{},
	plugins.PTR:   &BuiltinPluginPTR{},
	plugins.SOA:   &BuiltinPluginSOA{},
}

// Performs tests across all the builtin plugins where possible to simplify the actual plugin test files
func TestPluginVersion(t *testing.T) {
	for pluginType, plugin := range allPlugins {
		ver, err := plugin.PluginVersion()
		if err != nil {
			t.Errorf("unexpected error: %s, for plugin %s", err, pluginType)
		}

		if ver != utils.Version() {
			t.Errorf("incorrect version ''%s', want '%s', for %s", ver, utils.Version(), pluginType)
		}
	}
}

func TestPluginTypes(t *testing.T) {
	for pluginType, plugin := range allPlugins {
		actual, err := plugin.PluginTypes()
		if err != nil {
			t.Errorf("unexpected error: %s, for plugin %s", err, pluginType)
		}

		want := []plugins.PluginType{pluginType}
		if !cmp.Equal(actual, want) {
			t.Errorf("unexpected plugin types %s, want %s, for plugin %s", actual, want, pluginType)
		}
	}
}

func TestConfigure(t *testing.T) {
	config := &models.Config{}
	var pluginsToTest = map[plugins.PluginType]struct {
		plugin         plugins.ZoneMgrPlugin
		expectedConfig *models.Config
	}{
		plugins.A:     {plugin: &BuiltinPluginA{}, expectedConfig: nil},
		plugins.CNAME: {plugin: &BuiltinPluginCNAME{}, expectedConfig: nil},
		plugins.NS:    {plugin: &BuiltinPluginNS{}, expectedConfig: nil},
		plugins.PTR:   {plugin: &BuiltinPluginPTR{}, expectedConfig: nil},
		plugins.SOA:   {plugin: &BuiltinPluginSOA{}, expectedConfig: config},
	}

	for pluginType, pluginTest := range pluginsToTest {
		if err := pluginTest.plugin.Configure(config); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		// validate that the value of the internal config is the expectedConfig
		// We're going to use some reflection here sinze plugins.ZoneMgrPlugin doesn't expose a getter for Config
		value := reflect.ValueOf(pluginTest.plugin)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		if value.Kind() != reflect.Struct {
			t.Errorf("Unable to use reflection to get config value for %s", pluginType)
		}
		configField := value.FieldByName("config")
		if configField.IsValid() {
			// This plugin stores configuration
			actualConfig := reflect.NewAt(configField.Type(), unsafe.Pointer(configField.UnsafeAddr())).Elem().Interface().(*models.Config)
			if !cmp.Equal(pluginTest.expectedConfig, actualConfig) {
				t.Errorf("Unexpected config value for %s:\n%s", pluginType, cmp.Diff(pluginTest.expectedConfig, actualConfig))
			}
		}
	}
}

func TestValidateZone(t *testing.T) {
	// NOTE: CNAME and SOA are not in this list because they actually have a ValidateZone implementation
	pluginsToTest := map[plugins.PluginType]plugins.ZoneMgrPlugin{
		plugins.A:   &BuiltinPluginA{},
		plugins.NS:  &BuiltinPluginNS{},
		plugins.PTR: &BuiltinPluginPTR{},
	}

	for pluginType, plugin := range pluginsToTest {
		if err := plugin.ValidateZone("testing", &models.Zone{}); err != nil {
			t.Errorf("unexpected error: %s, for %s", err, pluginType)
		}
	}

}
