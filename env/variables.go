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

package env

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const EnvPrefix = "ZONEMGR_"

// Represents a value retrieve from the environment (or defaulted)
type Env struct {
	Value   string
	EnvName string
}

var (
	GenerateReverseLookupZones = &Env{EnvName: EnvPrefix + "GENERATE_REVERSE_LOOKUP_ZONES"}
	GenerateSerial             = &Env{EnvName: EnvPrefix + "GENERATE_SERIAL"}
	PluginsDirectory           = &Env{EnvName: EnvPrefix + "PLUGINS_DIR"}
	SerialChangeIndexDirectory = &Env{EnvName: EnvPrefix + "SERIAL_INDEX_DIR"}
)

func init() {
	defaultValues()
}

func defaultValues() {
	// Get the current user
	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine the current user, can not continue")
		os.Exit(1)
	}

	GenerateReverseLookupZones.Value = defaultValue(GenerateReverseLookupZones, "false")
	GenerateSerial.Value = defaultValue(GenerateSerial, "false")
	PluginsDirectory.Value = defaultValue(PluginsDirectory, filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "plugins"))
	SerialChangeIndexDirectory.Value = defaultValue(SerialChangeIndexDirectory, filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "serial"))

}

func defaultValue(e *Env, defaultValue string) string {
	value := os.Getenv(e.EnvName)

	if value == "" {
		value = defaultValue
	}

	return value
}
