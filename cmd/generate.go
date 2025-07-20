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

	parse "github.com/bcurnow/zonemgr/sourceyaml"
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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			absOutputDir, err := filepath.Abs(outputDir)
			if err != nil {
				fmt.Printf("Failed to resolve output directory %s: %v\n", outputDir, err)
				os.Exit(1)
			}
			outputDir = absOutputDir

			absInput, err := filepath.Abs(inputFile)
			if err != nil {
				fmt.Printf("Failed to resolve input file %s: %v\n", inputFile, err)
				os.Exit(1)
			}
			inputFile = absInput
		},
	}

	inputFile string
	outputDir string
)

func generateZoneFile() error {
	fmt.Printf("Generating BIND zone file(s) to directory %s using %s\n", outputDir, inputFile)
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("Failed to open input file %s: %w", inputFile, err)
	}

	zones, err := parse.ToZones(inputBytes)
	if err != nil {
		return fmt.Errorf("Failed to parse input file %s: %w", inputFile, err)

	}

	err = zonefile.ToZoneFiles(zones, outputDir)
	if err != nil {
		return fmt.Errorf("Failed to generate zone files: %w", err)
	}

	return nil
}

func init() {
	generateCmd.Flags().StringVarP(&inputFile, "input", "i", "zones.yaml", "Input YAML file")
	generateCmd.Flags().StringVarP(&outputDir, "outputDir", "d", ".", "Directory to output the BIND zone file(s) to")
	rootCmd.AddCommand(generateCmd)
}
