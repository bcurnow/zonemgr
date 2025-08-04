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

	"github.com/bcurnow/zonemgr/logging"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zonemgr",
	Short: "Converts YAML files to BIND zone files.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
}

var logLevel string
var logger = logging.Logger()

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
}

func setupLogging() {
	logger.SetLevel(hclog.LevelFromString(logLevel))

	if logger.GetLevel() == hclog.NoLevel {
		// Default to Warn
		logger.SetLevel(hclog.Info)
		logger.Error("Invalid log level specified, defaulting to Info", "LogLevel", logLevel)
	}
	logger.Trace("Log level set", "level", logLevel)
}

func toAbsoluteFilePath(path string, name string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve %s %s: %v\n", name, path, err)
		os.Exit(1)
	}
	return absPath
}
