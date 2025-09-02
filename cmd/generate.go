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
	"path/filepath"

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

			zoneFileGenerator = dns.PluginZoneFileGenerator(pluginManager.Plugins())
			normalizer = dns.PluginNormalizer(pluginManager.Plugins(), pluginManager.Metadata())
			zoneYamlParser = dns.YamlZoneParser(normalizer)

			// ensure that the serial-change-index-directory is an absolute file path
			absSerialChangeIndexDirectory, err := fs.ToAbsoluteFilePath(v.GetString("serial-change-index-directory"))
			if err != nil {
				return err
			}
			v.Set("serial-change-index-directory", absSerialChangeIndexDirectory)

			globalConfig = &models.Config{}
			globalConfig.GenerateReverseLookupZones = v.GetBool("generate-reverse-lookup-zones")
			globalConfig.GenerateSerial = v.GetBool("generate-serial")
			globalConfig.SerialChangeIndexDirectory = v.GetString("serial-change-index-directory")
			return nil
		},
	}

	inputFile                  string
	outputDir                  string
	generateReverseLookupZones bool
	generateSerial             bool
	serialChangeIndexDirectory string
	zoneReverser               dns.ZoneReverser = dns.Reverser()
	zoneFileGenerator          dns.ZoneFileGenerator
	zoneYamlParser             dns.ZoneParser
	normalizer                 dns.Normalizer
	globalConfig               *models.Config
)

func generateZoneFile() error {
	hclog.L().Info("Generating BIND zone file(s)", "outputDir", outputDir, "inputFile", inputFile)
	zones, err := zoneYamlParser.Parse(inputFile, globalConfig)
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
	generateCmd.Flags().StringVarP(&inputFile, "input-file", "", "zones.yaml", "Input YAML file")
	generateCmd.MarkFlagRequired("input")
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "", ".", "Directory to output the BIND zone file(s) to")
	generateCmd.Flags().BoolVarP(&generateReverseLookupZones, "generate-reverse-lookup-zones", "", false, "If true, reverse lookup zones will be generated as well")
	generateCmd.Flags().BoolVarP(&generateSerial, "generate-serial", "", false, "If true, the serial number on the SOA record will be automatically generated")
	generateCmd.Flags().StringVarP(&serialChangeIndexDirectory, "serial-change-index-directory", "", filepath.Join(homeDir, ".local", "share", "zonemgr", "serial"), "The directory to write the serial change index files to, these files keep track of the index portion of the serial number")

	rootCmd.AddCommand(generateCmd)

}
