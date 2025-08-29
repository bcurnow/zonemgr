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
	"path/filepath"
	"strings"

	"github.com/bcurnow/zonemgr/plugin_manager"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "zonemgr",
		Short: "Converts YAML files to BIND zone files.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initConfig(cmd); err != nil {
				return err
			}

			if v.GetBool("plugin-debug") {
				plugin_manager.EnablePluginDebug()
			}

			setupLogging()

			// Always load the plugins as the start of every command
			if err := pluginManager.LoadPlugins(v.GetString("plugins-dir")); err != nil {
				return err
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cleanup()
		},
	}

	pluginManager = plugin_manager.Manager()
	v             *viper.Viper
	homeDir       string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	// Only panic if we can't get the home directory.
	cobra.CheckErr(err)

	rootCmd.PersistentFlags().String("log-level", "info", "The log level (trace, debug, info, warn, error, fatal), not case sensitive")
	rootCmd.PersistentFlags().Bool("log-json", false, "If set, enables JSON loggiing output")
	rootCmd.PersistentFlags().Bool("log-time", false, "If set, prints the time on all the log messages")
	rootCmd.PersistentFlags().Bool("log-color", false, "If set, prints the log messages in color where possible")
	rootCmd.PersistentFlags().Bool("plugin-debug", false, "If set, will including plugin stdout/stderr in the log messages")
	rootCmd.PersistentFlags().String("plugin-dir", filepath.Join(homeDir, ".local", "share", "zonemgr", "plugins"), "The directory to find Zonemgr plugins")

}

func initConfig(cmd *cobra.Command) error {
	v = viper.New()
	v.SetEnvPrefix("zonemgr")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Bind all the cobra flags to viper
	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// Normalize the plugin-dir to an absolute path
	v.Set("plugin-dir", toAbsoluteFilePath(v.GetString("plugin-dir"), "plugin-dir"))

	return nil
}

func setupLogging() {
	level := hclog.LevelFromString(v.GetString("log-level"))

	if level == hclog.NoLevel {
		level = hclog.Info
		hclog.L().Error("Invalid log level specified, defaulting to Info", "level", v.GetString("log-level"))
	}
	utils.ConfigureLogging(level, v.GetBool("log-json"), !v.GetBool("log-time"), v.GetBool("log-color"))
	hclog.L().Trace("Log level set", "level", hclog.L().GetLevel())
}

// Ensures that any created plugin clients are properly cleaned up
func cleanup() {
	hclog.L().Trace("Cleaning up the clients...")
	goplugin.CleanupClients()
}
