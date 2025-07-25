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

	"github.com/bcurnow/zonemgr/parse"
	"github.com/bcurnow/zonemgr/templates"
	"github.com/bcurnow/zonemgr/zonefile"

	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates a BIND zone file from YAML input",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateZoneFile()
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			outputDir = toAbsoluteFilePath(outputDir, "output directory")
			inputFile = toAbsoluteFilePath(inputFile, "input file")
			zonefileTemplate = templateContent(zonefileTemplate, "zonefile template", templates.DefaultZoneFileTemplate)
			reverseZonefileTemplate = templateContent(reverseZonefileTemplate, "reverse zonefile template", templates.DefaultReverseZoneFileTemplate)
		},
	}

	inputFile               string
	outputDir               string
	zonefileTemplate        string
	reverseZonefileTemplate string
)

func generateZoneFile() error {
	fmt.Printf("Generating BIND zone file(s) to directory %s using %s\n", outputDir, inputFile)
	zones, err := parse.ToZones(inputFile)
	if err != nil {
		return fmt.Errorf("Failed to parse input file %s: %w", inputFile, err)

	}

	err = zonefile.ToZoneFiles(zones, outputDir, zonefileTemplate, reverseZonefileTemplate)
	if err != nil {
		return fmt.Errorf("Failed to generate zone files: %w", err)
	}

	return nil
}

func init() {
	generateCmd.Flags().StringVarP(&inputFile, "input", "i", "zones.yaml", "Input YAML file")
	generateCmd.MarkFlagRequired("input")
	generateCmd.Flags().StringVarP(&outputDir, "outputDir", "d", ".", "Directory to output the BIND zone file(s) to")
	generateCmd.MarkFlagRequired("outputDir")
	generateCmd.Flags().StringVarP(&zonefileTemplate, "zonefileTemplate", "z", "", "The go-lang template file to use to generate a zonefile, if unset will use a default template")
	generateCmd.Flags().StringVarP(&reverseZonefileTemplate, "reverseZonefileTemplate", "r", "", "The go-lang template file to use to generate a reverse zonefile, if unset will use a default template")
	rootCmd.AddCommand(generateCmd)
}
