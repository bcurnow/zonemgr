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

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/manager"
	"github.com/bcurnow/zonemgr/schema"
	"github.com/hashicorp/go-hclog"
)

func NormalizeZones(zones map[string]*schema.Zone) error {
	hclog.L().Trace("Normalizing zones", "count", len(zones))
	if len(zones) == 0 {
		return fmt.Errorf("no zones found")
	}

	registeredPlugins, err := manager.Plugins()
	if err != nil {
		return err
	}

	for name, zone := range zones {
		// Configure each of the plugins for this specific zone
		for _, plugin := range registeredPlugins {
			hclog.L().Trace("Calling Configure", "zoneName", name, "pluginName", plugin.PluginName)
			if nil == zone.Config {
				// This shouldn't happen unless there's an accident in the code (like perhaps creating new Zone for a reverse lookup zone and forgetting to populate Config)
				return fmt.Errorf("zone is missing config, zoneName=%s", name)
			}
			plugin.Plugin.Configure(zone.Config)
		}

		if err := normalizeZone(name, zone, registeredPlugins); err != nil {
			return err
		}

		// Now perform any validations on the zone itself
		for _, plugin := range registeredPlugins {
			if err := plugin.Plugin.ValidateZone(name, zone); err != nil {
				return err
			}
		}
	}

	return nil
}

func normalizeZone(name string, zone *schema.Zone, registeredPlugins map[plugins.PluginType]*plugins.Plugin) error {
	hclog.L().Trace("Normalizing zone", "name", name)
	for identifier, rr := range zone.ResourceRecords {
		hclog.L().Trace("Normalizing record", "identifier", identifier, "zoneName", name)
		plugin := registeredPlugins[plugins.PluginType(rr.Type)]
		if nil == plugin {
			return fmt.Errorf("unable to normalize zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
		}
		hclog.L().Trace("Calling Normalize on plugin", "identifier", identifier, "resourceRecordType", rr.Type, "zoneName", name, "plugin", plugin)
		err := plugin.Plugin.Normalize(identifier, rr)
		if err != nil {
			return err
		}
	}
	return nil
}
