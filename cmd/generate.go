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
	"os"
	"path/filepath"

	"github.com/bcurnow/zonemgr/ctx"
	"github.com/bcurnow/zonemgr/dns"
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
			outputDir = toAbsoluteFilePath(outputDir, "output directory")
			inputFile = toAbsoluteFilePath(inputFile, "input file")

			zoneFileGenerator = dns.PluginZoneFileGenerator(pluginManager.Plugins(), pluginManager.Metadata())
			normalizer = dns.PluginNormalizer(pluginManager)
			zoneYamlParser = dns.YamlZoneParser(normalizer)
			return nil
		},
	}

	inputFile                  string
	outputDir                  string
	generateReverseLookupZones bool
	generateSerial             bool
	serialChangeIndexDirectory string
	zoneReverser               dns.ZoneReverser = &dns.StandardZoneReverser{}
	zoneFileGenerator          dns.ZoneFileGenerator
	zoneYamlParser             dns.ZoneParser
	normalizer                 dns.Normalizer
)

func generateZoneFile() error {
	hclog.L().Info("Generating BIND zone file(s)", "outputDir", outputDir, "inputFile", inputFile)
	zones, err := zoneYamlParser.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse input file %s: %w", inputFile, err)

	}

	for name, zone := range zones {
		if err := zoneFileGenerator.GenerateZone(name, zone, outputDir); err != nil {
			return err
		}

		if zone.Config.GenerateReverseLookupZones {
			hclog.L().Debug("Zone has generate reverse lookup zones turned on", "zone", name)
			reverseLookupZones := zoneReverser.ReverseZone(name, zone)
			if err := normalizer.Normalize(reverseLookupZones); err != nil {
				return err
			}

			for name, zone := range reverseLookupZones {
				if err := zoneFileGenerator.GenerateZone(name, zone, outputDir); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func init() {
	generateCmd.Flags().StringVarP(&inputFile, "input-file`", "i", "zones.yaml", "Input YAML file")
	generateCmd.MarkFlagRequired("input")
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Directory to output the BIND zone file(s) to")
	generateCmd.Flags().BoolVarP(&generateReverseLookupZones, ctx.FlagGenerateReverseLookupZone, "r", false, "If true, reverse lookup zones will be generated as well")
	generateCmd.Flags().BoolVarP(&generateSerial, ctx.FlagGenerateSerial, "s", false, "If true, the serial number on the SOA record will be automatically generated")
	generateCmd.Flags().StringVarP(&serialChangeIndexDirectory, ctx.FlagSerialChangeIndexDirectory, "", "~/.local/share/zonemgr/serial", "The directory to write the serial change index files to, these files keep track of the index portion of the serial number")

	rootCmd.AddCommand(generateCmd)

}

func toAbsoluteFilePath(path string, name string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		hclog.L().Error("Could not convert %s value '%s' into an absolute path", name, path)
		os.Exit(1)
	}
	return absPath
}
