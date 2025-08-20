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

	"github.com/bcurnow/zonemgr/utils"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zonemgr",
	Short: "Converts YAML files to BIND zone files.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cleanup()
	},
}

var logLevel string
var logJsonFormat bool
var logTime bool
var logColor bool
var pluginDebug bool
var pluginsDir string

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
	rootCmd.PersistentFlags().StringVarP(&pluginsDir, "plugins-dir", "p", utils.PluginsDirectory.Value, "The directory to find Zonemgr plugins")
}

func setupLogging() {
	if pluginDebug {
		utils.EnablePluginDebug()
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

func toAbsoluteFilePath(path string, name string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve %s %s: %v\n", name, path, err)
		os.Exit(1)
	}
	return absPath
}

// Ensures that any created plugin clients are properly cleaned up
func cleanup() {
	hclog.L().Trace("Cleaning up the clients...")
	goplugin.CleanupClients()
}
