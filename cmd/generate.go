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
	"strconv"

	"github.com/bcurnow/zonemgr/env"
	"github.com/bcurnow/zonemgr/normalize"
	"github.com/bcurnow/zonemgr/parse"
	"github.com/bcurnow/zonemgr/plugins/manager"
	"github.com/bcurnow/zonemgr/zonefile"
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
			serialChangeIndexDirectory = toAbsoluteFilePath(serialChangeIndexDirectory, "serial change index directory")

			//Override the environment variables with any command line variables
			if cmd.Flags().Changed("generate-reverse-lookup-zones") {
				env.GenerateReverseLookupZones.Value = strconv.FormatBool(generateReverseLookupZones)
			}

			if cmd.Flags().Changed("generate-serial") {
				env.GenerateSerial.Value = strconv.FormatBool(generateSerial)
			}

			if cmd.Flags().Changed("serial-change-index-directory") {
				env.SerialChangeIndexDirectory.Value = serialChangeIndexDirectory
			}

			// Make sure we load up all the plugins at the start
			if _, err := manager.Default().Plugins(); err != nil {
				return err
			}
			return nil
		},
	}

	inputFile                  string
	outputDir                  string
	generateReverseLookupZones bool
	generateSerial             bool
	serialChangeIndexDirectory string
	zoneReverser               = zonefile.Reverser()
	zoneFileGenerator          = zonefile.Generator()
	zoneYamlParser             = parse.Parser()
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

		if zone.Config.GenerateReverseLookupZones != nil && *zone.Config.GenerateReverseLookupZones {
			hclog.L().Debug("Zone has generate reverse lookup zones turned on", "zone", name)
			reverseLookupZones := zoneReverser.ReverseZone(name, zone)
			if err := normalize.Default().Normalize(reverseLookupZones); err != nil {
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

	generateReverseLookupZonesEnvValue, err := strconv.ParseBool(env.GenerateReverseLookupZones.Value)
	if err != nil {
		generateReverseLookupZonesEnvValue = false
	}
	generateCmd.Flags().BoolVarP(&generateReverseLookupZones, "generate-reverse-lookup-zones", "r", generateReverseLookupZonesEnvValue, "If true, reverse lookup zones will be generated as well")
	generateSerialEnvValue, err := strconv.ParseBool(env.GenerateSerial.Value)
	if err != nil {
		generateSerialEnvValue = false
	}
	generateCmd.Flags().BoolVarP(&generateSerial, "generate-serial", "s", generateSerialEnvValue, "If true, the serial number on the SOA record will be automatically generated")
	generateCmd.Flags().StringVarP(&serialChangeIndexDirectory, "serial-change-index-directory", "", env.SerialChangeIndexDirectory.Value, "The directory to write the serial change index files to, these files keep track of the index portion of the serial number")

	rootCmd.AddCommand(generateCmd)

}
