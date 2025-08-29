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

package dns

import (
	"testing"

	models "github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/models/testingutils"
	"github.com/bcurnow/zonemgr/plugins"
)

var (
	mockAPlugin     *plugins.MockZoneMgrPlugin
	mockCNAMEPlugin *plugins.MockZoneMgrPlugin
	mockPlugins     map[plugins.PluginType]plugins.ZoneMgrPlugin
	mockMetadata    map[plugins.PluginType]*plugins.PluginMetadata
	testZone        *models.Zone
	testZones       map[string]*models.Zone
)

func dnsSetup(t *testing.T) {
	testingutils.Setup(t)
	mockAPlugin = plugins.NewMockZoneMgrPlugin(testingutils.MockController)
	mockCNAMEPlugin = plugins.NewMockZoneMgrPlugin(testingutils.MockController)

	mockPlugins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
	mockPlugins[plugins.A] = mockAPlugin
	mockPlugins[plugins.CNAME] = mockCNAMEPlugin

	mockMetadata = make(map[plugins.PluginType]*plugins.PluginMetadata)
	mockMetadata[plugins.A] = &plugins.PluginMetadata{Name: string(plugins.A), Command: "none", BuiltIn: true}
	mockMetadata[plugins.CNAME] = &plugins.PluginMetadata{Name: string(plugins.CNAME), Command: "none", BuiltIn: true}

	testZone = &models.Zone{
		Config: &models.Config{
			GenerateSerial:             true,
			SerialChangeIndexDirectory: "testZone-scid",
			GenerateReverseLookupZones: true,
		},
		ResourceRecords: map[string]*models.ResourceRecord{
			"record1": {Type: models.A},
			"record2": {Type: models.CNAME},
		},
		TTL: &models.TTL{
			Value:   testingutils.ToInt32Ptr(30),
			Comment: "testZone-TTL",
		},
	}
	testZones = make(map[string]*models.Zone)
	testZones["zone 1"] = testZone
	testZones["zone 2"] = testZone
}
