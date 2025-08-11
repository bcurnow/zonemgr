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

package zonefile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/manager"
	"github.com/bcurnow/zonemgr/schema"
	"github.com/hashicorp/go-hclog"
)

func writeZoneFile(name string, zone *schema.Zone, outputDir string) error {
	outputFile, err := os.Create(filepath.Join(outputDir, name))
	if err != nil {
		return fmt.Errorf("failed to create output file for zone %s: %w", name, err)
	}
	defer outputFile.Close()

	hclog.L().Info("Generating zone file", "outputFile", outputFile.Name(), "zone", name)

	// Write out the origin
	fmt.Fprintf(outputFile, "$ORIGIN %s\n", name)

	if zone.TTL != nil {
		fmt.Fprintln(outputFile, zone.TTL.Render())
	}

	registeredPlugins, err := manager.Plugins()
	if err != nil {
		return err
	}

	// Configure each of the plugins for this specific zone
	for _, plugin := range registeredPlugins {
		plugin.Plugin.Configure(zone.Config)
	}

	for identifier, rr := range zone.ResourceRecords {
		// We're takiing advantage of the fact that we have plugin types that match standard resource record types
		// so we can cast directly
		plugin := registeredPlugins[plugins.PluginType(rr.Type)]
		if nil == plugin {
			return fmt.Errorf("unable to write zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
		}
		renderedRecord, err := plugin.Plugin.Render(identifier, rr)
		if err != nil {
			return err
		}
		// It is possible that the plugin determines that this resource should not be rendered, don't add to the file
		if renderedRecord != "" {
			fmt.Fprintln(outputFile, renderedRecord)
		}
	}

	return nil
}
