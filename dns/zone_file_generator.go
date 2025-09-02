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

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/hashicorp/go-hclog"
)

type ZoneFileGenerator interface {
	GenerateZone(name string, zone *models.Zone, outputDir string) error
}
type pluginZoneFileGenerator struct {
	ZoneFileGenerator
	plugins  map[plugins.PluginType]plugins.ZoneMgrPlugin
	metadata map[plugins.PluginType]*plugins.PluginMetadata
}

func PluginZoneFileGenerator(plugins map[plugins.PluginType]plugins.ZoneMgrPlugin, metadata map[plugins.PluginType]*plugins.PluginMetadata) ZoneFileGenerator {
	return &pluginZoneFileGenerator{plugins: plugins, metadata: metadata}
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

	if err := plugins.WithSortedPlugins(zfg.plugins, zfg.metadata, func(pluginType plugins.PluginType, p plugins.ZoneMgrPlugin, metadata *plugins.PluginMetadata) error {
		p.Configure(zone.Config)
		return nil
	}); err != nil {
		return nil, err
	}

	if err := zone.WithSortedResourceRecords(func(identifier string, rr *models.ResourceRecord) error {
		// We're takiing advantage of the fact that we have plugin types that match standard resource record types
		// so we can cast directly
		plugin := zfg.plugins[plugins.PluginType(rr.Type)]
		if nil == plugin {
			return fmt.Errorf("unable to write zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
		}
		renderedRecord, err := plugin.Render(identifier, rr)
		if err != nil {
			return err
		}
		// It is possible that the plugin determines that this resource should not be rendered, don't add to the file
		if renderedRecord != "" {
			content.WriteString(renderedRecord)
			content.WriteString("\n")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return content.Bytes(), nil
}
