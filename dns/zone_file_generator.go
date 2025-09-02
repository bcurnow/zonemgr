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
	"bytes"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/hashicorp/go-hclog"
)

type ZoneFileGenerator interface {
	GenerateZone(name string, zone *models.Zone, outputDir string) error
}
type pluginZoneFileGenerator struct {
	ZoneFileGenerator
	plugins map[plugins.PluginType]plugins.ZoneMgrPlugin
}

func PluginZoneFileGenerator(plugins map[plugins.PluginType]plugins.ZoneMgrPlugin) ZoneFileGenerator {
	return &pluginZoneFileGenerator{plugins: plugins}
}

func (zfg *pluginZoneFileGenerator) GenerateZone(name string, zone *models.Zone, outputDir string) error {
	outputFileName := filepath.Join(outputDir, name)
	return fs.CreateFile(outputFileName, 0755, func() ([]byte, error) {
		hclog.L().Info("Generating zone file", "outputFile", outputFileName, "zone", name)
		return zfg.generate(name, zone)
	})
}

func (zfg *pluginZoneFileGenerator) generate(name string, zone *models.Zone) ([]byte, error) {
	var content bytes.Buffer
	// Write out the origin
	content.WriteString(fmt.Sprintf("$ORIGIN %s\n", name))

	if zone.TTL != nil {
		content.WriteString(zone.TTL.Render())
		content.WriteString("\n")
	}

	registeredPlugins := zfg.plugins

	// Configure each of the plugins for this specific zone
	for _, plugin := range registeredPlugins {
		plugin.Configure(zone.Config)
	}

	// We're going to sort the record by identifier to make the output deterministic
	identifiers := make([]string, 0, len(zone.ResourceRecords))
	for k := range zone.ResourceRecords {
		identifiers = append(identifiers, k)
	}
	sort.Strings(identifiers)

	for _, identifier := range identifiers {
		// We're takiing advantage of the fact that we have plugin types that match standard resource record types
		// so we can cast directly
		rr := zone.ResourceRecords[identifier]
		plugin := registeredPlugins[plugins.PluginType(rr.Type)]
		if nil == plugin {
			return nil, fmt.Errorf("unable to write zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
		}
		renderedRecord, err := plugin.Render(identifier, rr)
		if err != nil {
			return nil, err
		}
		// It is possible that the plugin determines that this resource should not be rendered, don't add to the file
		if renderedRecord != "" {
			content.WriteString(renderedRecord)
			content.WriteString("\n")
		}
	}

	return content.Bytes(), nil
}
