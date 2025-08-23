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
	"os"

	"github.com/bcurnow/zonemgr/ctx"
	"github.com/bcurnow/zonemgr/plugin_manager"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "zonemgr",
		Short: "Converts YAML files to BIND zone files.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			setupLogging()
			// Always load the plugin context at the start of every command
			if err := ctx.InitPluginContext(cmd.Flags()); err != nil {
				return err
			}

			// Always load the plugins as the start of every command
			if err := pluginManager.LoadPlugins(ctx.C().PluginsDirectory()); err != nil {
				return err
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cleanup()
		},
	}

	logLevel      string
	logJsonFormat bool
	logTime       bool
	logColor      bool
	pluginDebug   bool
	pluginsDir    string
	pluginManager = plugin_manager.Manager()
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "The log level (trace, debug, info, warn, error, fatal), not case sensitive")
	rootCmd.PersistentFlags().BoolVarP(&logJsonFormat, "log-json", "", false, "If set, enables JSON loggiing output")
	rootCmd.PersistentFlags().BoolVarP(&logTime, "log-time", "", false, "If set, prints the time on all the log messages")
	rootCmd.PersistentFlags().BoolVarP(&logColor, "log-color", "", true, "If set, prints the log messages in color where possible")
	rootCmd.PersistentFlags().BoolVarP(&pluginDebug, "plugin-debug", "", false, "If set, will including plugin stdout/stderr in the log messages")
	rootCmd.PersistentFlags().StringVarP(&pluginsDir, ctx.FlagPluginsDirectory, "p", "~/.local/share/zonemgr/plugins", "The directory to find Zonemgr plugins")
}

func setupLogging() {
	if pluginDebug {
		plugin_manager.EnablePluginDebug()
	}

	level := hclog.LevelFromString(logLevel)

	if level == hclog.NoLevel {
		// Default to Warn
		level = hclog.Info
		hclog.L().Error("Invalid log level specified, defaulting to Info", "level", logLevel)
	}
	utils.ConfigureLogging(level, logJsonFormat, !logTime, logColor)
	hclog.L().Trace("Log level set", "level", logLevel)

}

// Ensures that any created plugin clients are properly cleaned up
func cleanup() {
	hclog.L().Trace("Cleaning up the clients...")
	goplugin.CleanupClients()
}
