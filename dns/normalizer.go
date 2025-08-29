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
	"sort"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/hashicorp/go-hclog"
)

type Normalizer interface {
	Normalize(zones map[string]*models.Zone) error
}

type pluginNormalizer struct {
	Normalizer
	plugins  map[plugins.PluginType]plugins.ZoneMgrPlugin
	metadata map[plugins.PluginType]*plugins.PluginMetadata
}

func PluginNormalizer(plugins map[plugins.PluginType]plugins.ZoneMgrPlugin, metadata map[plugins.PluginType]*plugins.PluginMetadata) Normalizer {
	return &pluginNormalizer{plugins: plugins, metadata: metadata}
}

func (n *pluginNormalizer) Normalize(zones map[string]*models.Zone) error {
	hclog.L().Trace("Normalizing zones", "count", len(zones))
	if len(zones) == 0 {
		return fmt.Errorf("no zones found")
	}

	// Get the zone names so we can sort them and provide a stable iteration order, this will make the log output look better as well has help with testing
	zoneNames := make([]string, 0, len(zones))
	for k := range zones {
		zoneNames = append(zoneNames, k)
	}
	sort.Strings(zoneNames)

	registeredPlugins := n.plugins
	for _, name := range zoneNames {
		zone := zones[name]

		// Configure each of the plugins for this specific zone
		for pluginType, plugin := range registeredPlugins {
			pluginMetadata := n.metadata[pluginType]
			hclog.L().Debug("Calling Configure", "zoneName", name, "pluginName", pluginMetadata.Name)
			if nil == zone.Config {
				// This shouldn't happen unless there's an accident in the code (like perhaps creating new models.Zone for a reverse lookup zone and forgetting to populate Config)
				return fmt.Errorf("zone is missing config, zoneName=%s", name)
			}
			plugin.Configure(zone.Config)
		}

		if err := n.normalizeZone(name, zone); err != nil {
			return err
		}

		// Now perform any validations on the zone itself
		for pluginType, plugin := range registeredPlugins {
			pluginMetadata := n.metadata[pluginType]
			hclog.L().Debug("Calling Validatemodels.Zone", "zoneName", name, "pluginName", pluginMetadata.Name)
			if err := plugin.ValidateZone(name, zone); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *pluginNormalizer) normalizeZone(name string, zone *models.Zone) error {
	hclog.L().Debug("Normalizing zone", "name", name)
	for _, identifier := range zone.SortedResourceRecordKeys() {
		rr := zone.ResourceRecords[identifier]
		hclog.L().Trace("Normalizing record", "identifier", identifier, "zoneName", name)
		plugin := n.plugins[plugins.PluginType(rr.Type)]
		if nil == plugin {
			return fmt.Errorf("unable to normalize zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
		}
		hclog.L().Trace("Calling Normalize on plugin", "identifier", identifier, "resourceRecordType", rr.Type, "zoneName", name, "plugin", plugin)
		err := plugin.Normalize(identifier, rr)
		if err != nil {
			return err
		}
	}
	return nil
}
