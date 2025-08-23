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
	"slices"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/spf13/cobra"
)

var (
	pluginsCmd = &cobra.Command{
		Use:   "plugins",
		Short: "Prints information about the current plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			registeredPlugins := pluginManager.Plugins()

			// Get all the plugin types so we sort and print the plugins in order
			pluginTypes := make([]plugins.PluginType, 0, len(registeredPlugins))
			for pluginType := range registeredPlugins {
				pluginTypes = append(pluginTypes, pluginType)
			}
			slices.Sort(pluginTypes)

			formatString := "%-6s %-60s %-20s %s\n"
			// Turn on underline mode
			fmt.Println("\033[4m")
			fmt.Printf(formatString, "Type", "Name", "Version", "Plugin Command")
			// Turn off underline mode
			fmt.Print("\033[24m")

			for _, pluginType := range pluginTypes {
				plugin := registeredPlugins[pluginType]
				pluginMetadata := pluginManager.Metadata()[pluginType]

				pluginVersion, err := plugin.PluginVersion()
				if err != nil {
					return err
				}
				fmt.Printf(formatString, pluginType, pluginMetadata.Name, pluginVersion, pluginMetadata.Command)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(pluginsCmd)
}
