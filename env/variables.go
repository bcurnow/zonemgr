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
	"path/filepath"

	"github.com/bcurnow/zonemgr/logging"
)

const EnvPrefix = "ZONEMGR_"

// Represents a value retrieve from the environment (or defaulted)
type Env struct {
	Value   string
	EnvName string
}

var (
	PLUGINS *Env = &Env{EnvName: EnvPrefix + "PLUGINS"}
)

var logger = logging.Logger().Named("env")

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to determine user's home directory, can not continue")
		os.Exit(1)
	}

	PLUGINS.Value = defaultEnv(PLUGINS, filepath.Join(homeDir, ".local", "share", "yamlconv", "plugins"))
}

func defaultEnv(e *Env, defaultValue string) string {
	value := os.Getenv(e.EnvName)

	if value == "" {
		value = defaultValue
	}

	return value
}
