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
	"text/template"

	"github.com/bcurnow/zonemgr/parse"
	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates the various files used by zonemgr",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			input = toAbsoluteFilePath(input, "input")
		},
	}

	validateYamlCmd = &cobra.Command{
		Use:   "yaml",
		Short: "Validates the YAML input file",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := parse.ToZones(input)
			if err != nil {
				return fmt.Errorf("Failed to parse input file %s: %w", input, err)

			}
			fmt.Printf("%s is valid\n", input)
			return nil
		},
	}

	validateTemplate = &cobra.Command{
		Use:   "template",
		Short: "Validates a go template file",
		RunE: func(cmd *cobra.Command, args []string) error {
			templateContent, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("Unable to read %s: %v\n", input, err)
			}
			_, err = template.New("template").Parse(string(templateContent))
			if err != nil {
				return fmt.Errorf("Failed to parse template: %w", err)
			}
			fmt.Printf("%s is valid\n", input)
			return nil
		},
	}

	input string
)

func init() {
	validateCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "The input file to validate")
	validateCmd.MarkPersistentFlagRequired("input")
	validateCmd.AddCommand(validateYamlCmd)
	validateCmd.AddCommand(validateTemplate)
	rootCmd.AddCommand(validateCmd)
}
