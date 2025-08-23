/**
 * Copyright (C) 2025 bcurnow
 *
 * This file is part of yamlconv.
 *
 * yamlconv is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * yamlconv is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with yamlconv.  If not, see <https://www.gnu.org/licenses/>.
 */
package cmd

import (
	"fmt"

	"github.com/bcurnow/zonemgr/ctx"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Prints the environment variables used (or defaulted)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s=\"%t\"\n", ctx.GenerateReverseLookupZonesEnvName, ctx.C().GenerateReverseLookupZones())
		fmt.Printf("%s=\"%t\"\n", ctx.GenerateSerialEnvName, ctx.C().GenerateSerial())
		fmt.Printf("%s=\"%s\"\n", ctx.PluginsDirectoryEnvName, ctx.C().PluginsDirectory())
		fmt.Printf("%s=\"%s\"\n", ctx.SerialChangeIndexDirectoryEnvName, ctx.C().SerialChangeIndexDirectory())
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
