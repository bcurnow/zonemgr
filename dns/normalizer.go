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
	"fmt"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/hashicorp/go-hclog"
)

type Normalizer interface {
	Normalize(zones map[string]*models.Zone, globalConfig *models.Config) error
}

type pluginNormalizer struct {
	Normalizer
	plugins  map[plugins.PluginType]plugins.ZoneMgrPlugin
	metadata map[plugins.PluginType]*plugins.Metadata
}

func PluginNormalizer(plugins map[plugins.PluginType]plugins.ZoneMgrPlugin, metadata map[plugins.PluginType]*plugins.Metadata) Normalizer {
	return &pluginNormalizer{plugins: plugins, metadata: metadata}
}

func (n *pluginNormalizer) Normalize(zones map[string]*models.Zone, globalConfig *models.Config) error {
	hclog.L().Trace("Normalizing zones", "count", len(zones))
	if len(zones) == 0 {
		return fmt.Errorf("no zones found")
	}

	return models.WithSortedZones(zones, func(name string, zone *models.Zone) error {
		// Normalize the config if necessary
		if err := n.normalizeConfig(name, zone, globalConfig); err != nil {
			return err
		}

		// Configure each of the plugins for this specific zone
		// We need to do multiple loops over the plugins because we need all the plugins configured
		// Then all the normalization done
		// Then all the zone validation
		// If we do this in a single loop, we'd end up calling ValidateZone before all the normalization for the zone is complete
		if err := plugins.WithSortedPlugins(n.plugins, n.metadata, func(pluginType plugins.PluginType, p plugins.ZoneMgrPlugin, metadata *plugins.Metadata) error {
			hclog.L().Debug("Calling Configure", "zoneName", name, "pluginName", metadata.Name)
			p.Configure(zone.Config)
			return nil
		}); err != nil {
			return err
		}

		if err := n.normalizeZone(name, zone); err != nil {
			return err
		}

		// Now perform any validations on the zone itself
		if err := plugins.WithSortedPlugins(n.plugins, n.metadata, func(pluginType plugins.PluginType, p plugins.ZoneMgrPlugin, metadata *plugins.Metadata) error {
			hclog.L().Debug("Calling Validatemodels.Zone", "zoneName", name, "pluginName", metadata.Name)
			if err := p.ValidateZone(name, zone); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
}

func (n *pluginNormalizer) normalizeConfig(name string, zone *models.Zone, globalConfig *models.Config) error {
	if nil == zone.Config {
		hclog.L().Debug("Zone missing config, setting to global config", "zoneName", name, "globalConfig", globalConfig)
		zone.Config = globalConfig
	}

	// Ensure that the serial change index directory is an absolute path
	hclog.L().Trace("Ensuring serial-change-index-directory is an absolute path", "serialChangeIndexDirectory", zone.Config.SerialChangeIndexDirectory)
	absSerialChangeIndexDirectory, err := fs.ToAbsoluteFilePath(zone.Config.SerialChangeIndexDirectory)
	if err != nil {
		return err
	}
	zone.Config.SerialChangeIndexDirectory = absSerialChangeIndexDirectory
	return nil
}

func (n *pluginNormalizer) normalizeZone(name string, zone *models.Zone) error {
	hclog.L().Debug("Normalizing zone", "name", name)
	if err := zone.WithSortedResourceRecords(func(identifier string, rr *models.ResourceRecord) error {
		hclog.L().Trace("Normalizing record", "identifier", identifier, "zoneName", name)
		// We only call normalize on the resource record types we have plugins for, no need to loop
		plugin := n.plugins[plugins.PluginType(rr.Type)]
		if nil == plugin {
			return fmt.Errorf("unable to normalize zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
		}
		hclog.L().Trace("Calling Normalize on plugin", "identifier", identifier, "resourceRecordType", rr.Type, "zoneName", name, "plugin", plugin)
		err := plugin.Normalize(identifier, rr)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
