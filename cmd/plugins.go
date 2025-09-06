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

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/spf13/cobra"
)

const (
	pluginFormatBase         = "%-6s %-60s %-20s %s"
	pluginFormatString       = pluginFormatBase + "\n"
	pluginHeaderFormatString = "\033[4m" + pluginFormatBase + "\033[24m\n"
)

var (
	pluginHeaders = []any{"Type", "Name", "Version", "Command"}
	pluginsCmd    = &cobra.Command{
		Use:   "plugins",
		Short: "Prints information about the current plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf(pluginHeaderFormatString, pluginHeaders...)

			if err := plugins.WithSortedPlugins(pluginManager.Plugins(), pluginManager.Metadata(), func(pluginType plugins.Type, p plugins.ZoneMgrPlugin, metadata *plugins.Metadata) error {
				pluginVersion, err := p.PluginVersion()
				if err != nil {
					return err
				}
				fmt.Printf(pluginFormatString, pluginType, metadata.Name, pluginVersion, metadata.Command)
				return nil
			}); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(pluginsCmd)
}
