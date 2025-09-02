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

	"github.com/bcurnow/zonemgr/internal/mocks"
	models "github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/golang/mock/gomock"
)

var (
	mockController  *gomock.Controller
	mockAPlugin     *mocks.MockZoneMgrPlugin
	mockCNAMEPlugin *mocks.MockZoneMgrPlugin
	mockPlugins     map[plugins.PluginType]plugins.ZoneMgrPlugin
	mockMetadata    map[plugins.PluginType]*plugins.PluginMetadata
	mockFs          *mocks.MockFileSystem
	testZone        *models.Zone
	testZones       map[string]*models.Zone
	globalConfig    = &models.Config{SerialChangeIndexDirectory: "global-serial-change-index-directory"}
)

func dnsSetup(t *testing.T) {
	mockController = gomock.NewController(t)
	mockAPlugin = mocks.NewMockZoneMgrPlugin(mockController)
	mockCNAMEPlugin = mocks.NewMockZoneMgrPlugin(mockController)

	mockPlugins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
	mockPlugins[plugins.A] = mockAPlugin
	mockPlugins[plugins.CNAME] = mockCNAMEPlugin

	mockMetadata = make(map[plugins.PluginType]*plugins.PluginMetadata)
	mockMetadata[plugins.A] = &plugins.PluginMetadata{Name: string(plugins.A), Command: "none", BuiltIn: true}
	mockMetadata[plugins.CNAME] = &plugins.PluginMetadata{Name: string(plugins.CNAME), Command: "none", BuiltIn: true}

	mockFs = mocks.NewMockFileSystem(mockController)
	fs = mockFs

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
			Value:   toInt32Ptr(30),
			Comment: "testZone-TTL",
		},
	}
	testZones = make(map[string]*models.Zone)
	testZones["zone1"] = testZone
	testZones["zone2"] = testZone

	globalConfig = &models.Config{}
}

func dnsTeardown(_ *testing.T) {
	mockController.Finish()
}

func toInt32Ptr(i int32) *int32 {
	return &i
}
