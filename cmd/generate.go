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

			return nil
		},
	}

	inputFile         string
	outputDir         string
	zoneReverser      dns.ZoneReverser = dns.Reverser()
	zoneFileGenerator dns.ZoneFileGenerator
	normalizer        dns.Normalizer
)

func generateZoneFile() error {
	hclog.L().Info("Generating BIND zone file(s)", "outputDir", outputDir, "inputFile", inputFile)
	zones, err := parser.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse input file %s: %w", inputFile, err)
	}

	if err := models.WithSortedZones(zones, func(name string, zone *models.Zone) error {
		if err := zoneFileGenerator.GenerateZone(name, zone, outputDir); err != nil {
			return err
		}
		if err := generateReverseLookupZones(name, zone); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func generateReverseLookupZones(name string, zone *models.Zone) error {
	if !zone.Config.GenerateReverseLookupZones {
		return nil
	}

	if zone.Config.GenerateReverseLookupZones {
		hclog.L().Debug("Zone has generate reverse lookup zones turned on", "zone", name)
		reverseLookupZones := zoneReverser.ReverseZone(name, zone)
		if err := normalizer.Normalize(reverseLookupZones); err != nil {
			return err
		}

		if err := models.WithSortedZones(reverseLookupZones, func(name string, zone *models.Zone) error {
			if err := zoneFileGenerator.GenerateZone(name, zone, outputDir); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	generateCmd.Flags().StringVar(&inputFile, "input-file", "zones.yaml", "Input YAML file")
	if err := generateCmd.MarkFlagRequired("input-file"); err != nil {
		panic(err) // This should never happen in init
	}
	generateCmd.Flags().StringVar(&outputDir, "output-dir", ".", "Directory to output the BIND zone file(s) to")

	rootCmd.AddCommand(generateCmd)

}
