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
	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates the various files used by zonemgr",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if rootCmd.PersistentPreRun != nil {
				rootCmd.PersistentPreRun(cmd, args)
			}

			input = toAbsoluteFilePath(input, "input")

			parser = dns.YamlZoneParser(dns.PluginNormalizer(pluginManager))
		},
	}

	validateYamlCmd = &cobra.Command{
		Use:   "yaml",
		Short: "Validates the YAML input file",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := parser.Parse(input)
			if err != nil {
				return fmt.Errorf("failed to parse input file %s: %w", input, err)

			}
			fmt.Printf("%s is valid\n", input)
			return nil
		},
	}

	input  string
	parser dns.ZoneParser
)

func init() {
	validateCmd.PersistentFlags().StringVarP(&input, "input", "", "", "The input file to validate")
	validateCmd.MarkPersistentFlagRequired("input")
	validateCmd.AddCommand(validateYamlCmd)
	rootCmd.AddCommand(validateCmd)
}
