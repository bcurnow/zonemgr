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
package cmd

import (
	"testing"

	"github.com/bcurnow/zonemgr/dns"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/plugin_manager"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/golang/mock/gomock"
)

var (
	mockController        *gomock.Controller
	mockFs                *utils.MockFileSystemOperations
	mockPluginManager     *plugin_manager.MockPluginManager
	mockParser            *dns.MockZoneParser
	mockZoneReverser      *dns.MockZoneReverser
	mockZoneFileGenerator *dns.MockZoneFileGenerator
	mockNormalizer        *dns.MockNormalizer
	testPlugin            *plugins.MockZoneMgrPlugin
	testPlugins           map[plugins.PluginType]plugins.ZoneMgrPlugin
	testMetadata          map[plugins.PluginType]*plugins.Metadata
)

func setup(t *testing.T) {
	mockController = gomock.NewController(t)
	mockFs = utils.NewMockFileSystemOperations(mockController)
	fs = mockFs

	mockPluginManager = plugin_manager.NewMockPluginManager(mockController)
	pluginManager = mockPluginManager

	mockParser = dns.NewMockZoneParser(mockController)
	parser = mockParser

	mockZoneReverser = dns.NewMockZoneReverser(mockController)
	zoneReverser = mockZoneReverser

	mockZoneFileGenerator = dns.NewMockZoneFileGenerator(mockController)
	zoneFileGenerator = mockZoneFileGenerator

	mockNormalizer = dns.NewMockNormalizer(mockController)
	normalizer = mockNormalizer

	testPlugin = plugins.NewMockZoneMgrPlugin(mockController)
	testPlugins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
	testMetadata = make(map[plugins.PluginType]*plugins.Metadata)

	testPlugins[plugins.A] = testPlugin
	testPlugins[plugins.NS] = testPlugin
	testMetadata[plugins.A] = &plugins.Metadata{Name: "A", Command: "A"}
	testMetadata[plugins.NS] = &plugins.Metadata{Name: "NS", Command: "NS"}
}

func teardown(_ *testing.T) {
	defer func() { fs = &utils.FileSystem{} }()
	defer func() { pluginManager = plugin_manager.Manager() }()
	mockController.Finish()
}
