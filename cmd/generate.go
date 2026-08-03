/*
Copyright © 2025 Brian Curnow

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"

	"github.com/bcurnow/zonemgr/dns"
	"github.com/bcurnow/zonemgr/models"
	"github.com/hashicorp/go-hclog"

	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates a BIND zone file from YAML input",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateZoneFile()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			absOutputDir, err := fs.ToAbsoluteFilePath(outputDir)
			if err != nil {
				return err
			}
			outputDir = absOutputDir

			absInputFile, err := fs.ToAbsoluteFilePath(inputFile)
			if err != nil {
				return err
			}
			inputFile = absInputFile

			zoneFileGenerator = dns.PluginZoneFileGenerator(pluginManager.Plugins(), pluginManager.Metadata())
			normalizer = dns.PluginNormalizer(pluginManager.Plugins(), pluginManager.Metadata())
			parser = dns.YamlZoneParser(normalizer)
			catalogGenerator = dns.PluginCatalogGenerator(pluginManager.Plugins(), pluginManager.Metadata())

			return nil
		},
	}

	inputFile         string
	outputDir         string
	zoneReverser      dns.ZoneReverser = dns.Reverser()
	zoneFileGenerator dns.ZoneFileGenerator
	normalizer        dns.Normalizer
	catalogGenerator  dns.CatalogGenerator
)

func generateZoneFile() error {
	hclog.L().Info("generating BIND zone file(s)", "outputDir", outputDir, "inputFile", inputFile)
	zones, err := parser.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse input file %s: %w", inputFile, err)
	}

	var memberZoneNames []string
	catalogZones := make(map[string]*models.Zone)
	reverseZones := make(map[string]*models.Zone)

	// Pass 1: compute the full set of forward, reverse and catalog zones without writing anything.
	// A catalog zone needs to know about every other zone before its file can be written, so nothing
	// is written until this pass completes.
	if err := models.WithSortedZones(zones, func(name string, zone *models.Zone) error {
		if zone.Config.IsCatalog {
			// Catalog zones are never members of any catalog, not even themselves, and have no
			// reverse-lookup zones of their own.
			catalogZones[name] = zone
			return nil
		}

		memberZoneNames = append(memberZoneNames, name)

		if !zone.Config.GenerateReverseLookupZones {
			return nil
		}

		hclog.L().Debug("zone has generate reverse lookup zones turned on", "zone", name)
		zoneReverseZones, err := zoneReverser.ReverseZone(name, zone)
		if err != nil {
			return err
		}
		return mergeReverseZones(reverseZones, zoneReverseZones)
	}); err != nil {
		return err
	}

	if len(reverseZones) > 0 {
		if err := normalizer.Normalize(reverseZones); err != nil {
			return err
		}
	}

	if err := populateCatalogZones(catalogZones, memberZoneNames, reverseZones); err != nil {
		return err
	}

	// Pass 2: write everything now that every zone is fully populated.
	if err := models.WithSortedZones(zones, func(name string, zone *models.Zone) error {
		if zone.Config.IsCatalog {
			return nil
		}
		return zoneFileGenerator.GenerateZone(name, zone, outputDir)
	}); err != nil {
		return err
	}

	if err := models.WithSortedZones(reverseZones, func(name string, zone *models.Zone) error {
		return zoneFileGenerator.GenerateZone(name, zone, outputDir)
	}); err != nil {
		return err
	}

	return models.WithSortedZones(catalogZones, func(name string, zone *models.Zone) error {
		return zoneFileGenerator.GenerateZone(name, zone, outputDir)
	})
}

// mergeReverseZones merges newZones into accumulated. The first source zone (in processing order) to
// produce a given reverse zone name establishes that zone's Config, TTL and SOA record; later source
// zones producing the same reverse zone (e.g. two forward zones with hosts in the same subnet)
// contribute only their PTR records, rather than overwriting the earlier zone entirely.
func mergeReverseZones(accumulated, newZones map[string]*models.Zone) error {
	for zoneName, newZone := range newZones {
		existing, ok := accumulated[zoneName]
		if !ok {
			accumulated[zoneName] = newZone
			continue
		}

		for identifier, rr := range newZone.ResourceRecords {
			if rr.Type == models.SOA {
				// The zone that first established this reverse zone owns its SOA record.
				continue
			}
			if existingRR, ok := existing.ResourceRecords[identifier]; ok {
				if existingRR.Value != rr.Value {
					return fmt.Errorf("conflicting PTR record for '%s' in reverse zone '%s': '%s' vs '%s'", identifier, zoneName, existingRR.Value, rr.Value)
				}
				continue
			}
			existing.ResourceRecords[identifier] = rr
		}
	}
	return nil
}

// populateCatalogZones injects the RFC 9432 catalog records into each catalog zone found during pass 1.
func populateCatalogZones(catalogZones map[string]*models.Zone, memberZoneNames []string, reverseZones map[string]*models.Zone) error {
	reverseZoneNames := make([]string, 0, len(reverseZones))
	for name := range reverseZones {
		reverseZoneNames = append(reverseZoneNames, name)
	}

	return models.WithSortedZones(catalogZones, func(name string, zone *models.Zone) error {
		members := memberZoneNames
		if zone.Config.CatalogIncludeReverseZones {
			members = append(append([]string{}, memberZoneNames...), reverseZoneNames...)
		}
		return catalogGenerator.AddCatalogRecords(name, zone, members)
	})
}

func init() {
	generateCmd.Flags().StringVar(&inputFile, "input-file", "zones.yaml", "Input YAML file")
	cobra.CheckErr(generateCmd.MarkFlagRequired("input-file"))
	generateCmd.Flags().StringVar(&outputDir, "output-dir", ".", "Directory to output the BIND zone file(s) to")

	rootCmd.AddCommand(generateCmd)

}
