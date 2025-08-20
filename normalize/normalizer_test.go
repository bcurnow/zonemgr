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

package normalize

import (
	"fmt"
	"testing"

	plugins "github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
	"github.com/bcurnow/zonemgr/test"
	"github.com/golang/mock/gomock"
)

var (
	mockController    *gomock.Controller
	mockAPlugin       *test.MockZoneMgrPlugin
	mockCNAMEPlugin   *test.MockZoneMgrPlugin
	mockPluginManager *test.MockPluginManager
	mockPlugins       map[plugins.PluginType]*plugins.Plugin

	testZone = &schema.Zone{
		Config: &schema.Config{
			PluginsDirectory:           "testZone-plugins",
			GenerateSerial:             true,
			SerialChangeIndexDirectory: "testZone-scid",
			GenerateReverseLookupZones: true,
		},
		ResourceRecords: map[string]*schema.ResourceRecord{
			"record1": {Type: schema.A},
			"record2": {Type: schema.CNAME},
		},
		TTL: &schema.TTL{
			Value:   test.ToInt32Ptr(30),
			Comment: "testZone-TTL",
		},
	}
	testZones map[string]*schema.Zone
)

func setup(t *testing.T) {
	mockController = gomock.NewController(t)
	mockPluginManager = test.NewMockPluginManager(mockController)

	// Replace the package pluginManager with the mock
	pluginManager = mockPluginManager

	mockAPlugin = test.NewMockZoneMgrPlugin(mockController)
	mockCNAMEPlugin = test.NewMockZoneMgrPlugin(mockController)
	mockPlugins = make(map[plugins.PluginType]*plugins.Plugin)
	mockPlugins[plugins.A] = &plugins.Plugin{PluginName: "Mock A Plugin", Plugin: mockAPlugin}
	mockPlugins[plugins.CNAME] = &plugins.Plugin{PluginName: "Mock CNAME Plugin", Plugin: mockCNAMEPlugin}

	testZones = make(map[string]*schema.Zone)
	testZones["zone 1"] = testZone
	testZones["zone 2"] = testZone
}

func teardown(_ *testing.T) {
	mockController.Finish()
}

func TestNormalizeZones(t *testing.T) {
	setup(t)
	defer teardown(t)

	mockPluginManager.EXPECT().Plugins().Return(mockPlugins, nil)

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(testZone.Config).Times(2)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config).Times(2)

	// Each plugin should have normalize called for each zone
	mockAPlugin.EXPECT().Normalize("record1", testZone.ResourceRecords["record1"]).Times(2)
	mockCNAMEPlugin.EXPECT().Normalize("record2", testZone.ResourceRecords["record2"]).Times(2)

	// Each plugin should have validate called for each zone
	mockAPlugin.EXPECT().ValidateZone("zone 1", testZone)
	mockCNAMEPlugin.EXPECT().ValidateZone("zone 1", testZone)
	mockAPlugin.EXPECT().ValidateZone("zone 2", testZone)
	mockCNAMEPlugin.EXPECT().ValidateZone("zone 2", testZone)

	if err := Default().Normalize(testZones); err != nil {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}

func TestNormalizeZones_NoZones(t *testing.T) {
	setup(t)
	defer teardown(t)

	if err := Default().Normalize(map[string]*schema.Zone{}); err != nil {
		if err.Error() != "no zones found" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}

func TestNormalizeZones_PluginManagerError(t *testing.T) {
	setup(t)
	defer teardown(t)

	mockPluginManager.EXPECT().Plugins().Return(nil, fmt.Errorf("Testing Plugin Manager Error"))

	if err := Default().Normalize(testZones); err != nil {
		if err.Error() != "Testing Plugin Manager Error" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}

func TestNormalizeZone_NilConfig(t *testing.T) {
	setup(t)
	defer teardown(t)

	mockPluginManager.EXPECT().Plugins().Return(mockPlugins, nil)

	if err := Default().Normalize(map[string]*schema.Zone{"nil config zone": {Config: nil}}); err != nil {
		if err.Error() != "zone is missing config, zoneName=nil config zone" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)

	}
}
func TestNormalizeZones_NoPluginForRecordType(t *testing.T) {
	setup(t)
	defer teardown(t)

	mockPluginManager.EXPECT().Plugins().Return(mockPlugins, nil)

	invalidZone := &schema.Zone{
		Config: testZone.Config,
		ResourceRecords: map[string]*schema.ResourceRecord{
			"bad type": {Type: "bogus"},
		},
		TTL: testZone.TTL,
	}

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(invalidZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(invalidZone.Config)

	if err := Default().Normalize(map[string]*schema.Zone{"invalid zone": invalidZone}); err != nil {
		if err.Error() != "unable to normalize zone 'invalid zone', no plugin for resource record type 'bogus', identifier: 'bad type'" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}
func TestNormalizeZones_NormalizeError(t *testing.T) {
	setup(t)
	defer teardown(t)

	mockPluginManager.EXPECT().Plugins().Return(mockPlugins, nil)

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)

	mockAPlugin.EXPECT().Normalize("record1", testZone.ResourceRecords["record1"]).Return(fmt.Errorf("test normalize error"))

	if err := Default().Normalize(testZones); err != nil {
		if err.Error() != "test normalize error" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}

func TestNormalizeZones_ValidateError(t *testing.T) {
	setup(t)
	defer teardown(t)

	mockPluginManager.EXPECT().Plugins().Return(mockPlugins, nil)

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)

	// Each plugin should have normalize called for each zone
	mockAPlugin.EXPECT().Normalize("record1", testZone.ResourceRecords["record1"])
	mockCNAMEPlugin.EXPECT().Normalize("record2", testZone.ResourceRecords["record2"])

	// Each plugin should have validate called for each zone, we're not sure which order they will iterate in so EXPECT for both
	mockAPlugin.EXPECT().ValidateZone("zone 1", testZone).Return(fmt.Errorf("test validate error")).MaxTimes(1)
	mockCNAMEPlugin.EXPECT().ValidateZone("zone 1", testZone).Return(fmt.Errorf("test validate error")).MaxTimes(1)

	if err := Default().Normalize(testZones); err != nil {
		if err.Error() != "test validate error" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}
