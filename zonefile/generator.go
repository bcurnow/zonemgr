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
)

func generateZone(name string, zone *schema.Zone, outputDir string) error {
	outputFile, err := os.Create(filepath.Join(outputDir, name))
	if err != nil {
		return fmt.Errorf("Failed to create output file for zone %s: %w", name, err)
	}
	defer outputFile.Close()

	logger.Info("Generating zone file", "outputFile", outputFile.Name(), "zone", name)

	// Write out the origin
	fmt.Fprintf(outputFile, "$ORIGIN %s\n", name)

	// Write out the TTL
	if zone.TTL.Value != nil {
		comment := zone.TTL.Comment
		if comment != "" {
			comment = " ; " + comment
		}
		fmt.Fprintf(outputFile, "$TTL %d%s\n", *zone.TTL.Value, comment)
	}

	registeredPlugins, err := manager.Plugins()
	if err != nil {
		return err
	}

	for identifier, rr := range zone.ResourceRecords {
		plugin := registeredPlugins[plugins.PluginType(rr.Type)]
		if nil == plugin {
			return fmt.Errorf("Unable to write zone '%s', no plugin for resource record type '%s', identifier: '%s'", name, rr.Type, identifier)
		}
		logger.Trace("Resource Record to render", "resourceRecord", rr)
		renderedRecord, err := plugin.Render(identifier, rr)
		logger.Trace("Rendered resource record", "string", renderedRecord)
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
