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

package parse

import (
	"fmt"

	"github.com/bcurnow/zonemgr/logging"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/manager"
	"github.com/bcurnow/zonemgr/schema"
)

var logger = logging.Logger().Named("normalize")

func normalize(zones map[string]*schema.Zone) (map[string]*schema.Zone, error) {
	if len(zones) == 0 {
		return nil, fmt.Errorf("No zones found")
	}

	registeredPlugins, err := manager.Plugins()
	if err != nil {
		return nil, err
	}

	for name, zone := range zones {
		for identifier, rr := range zone.ResourceRecords {
			plugin := registeredPlugins[plugins.PluginType(rr.Type)]
			if nil == plugin {
				return nil, fmt.Errorf("Unable to normalize zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
			}

			logger.Trace("Resource record to normalize", "resourceRecord", rr)
			normalizedRR, err := plugin.Normalize(identifier, rr)
			if err != nil {
				return nil, err
			}
			logger.Trace("Normalized resource record", "resourceRecord", normalizedRR)
			zone.ResourceRecords[identifier] = normalizedRR
		}
		// Now perform any validations on the zone itself
		for _, plugin := range registeredPlugins {
			if err := plugin.ValidateZone(name, *zone); err != nil {
				return nil, err
			}
		}
	}

	return zones, nil
}
