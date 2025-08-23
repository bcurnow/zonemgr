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

package ctx

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

const (
	envPrefix                         = "ZONEMGR_"
	GenerateReverseLookupZonesEnvName = envPrefix + "GENERATE_REVERSE_LOOKUP_ZONES"
	GenerateSerialEnvName             = envPrefix + "GENERATE_SERIAL"
	PluginsDirectoryEnvName           = envPrefix + "PLUGINS_DIR"
	SerialChangeIndexDirectoryEnvName = envPrefix + "SERIAL_INDEX_DIR"
)

type Environment interface {
	GenerateReverseLookupZones() bool
	GenerateSerial() bool
	PluginsDirectory() string
	SerialChangeIndexDirectory() string
	LoadValues() error
}

type environment struct {
	Environment
	generateReverseLookupZones bool
	generateSerial             bool
	pluginsDirectory           string
	serialChangeIndexDirectory string
}

func (e *environment) GenerateReverseLookupZones() bool {
	return e.generateReverseLookupZones
}

func (e *environment) GenerateSerial() bool {
	return e.generateSerial
}

func (e *environment) PluginsDirectory() string {
	return e.pluginsDirectory
}

func (e *environment) SerialChangeIndexDirectory() string {
	return e.serialChangeIndexDirectory
}

func (e *environment) LoadValues() error {
	// Get the current user
	user, err := user.Current()
	if err != nil {
		return err
	}

	generateReverseLookupZones, err := strconv.ParseBool(e.envValueOrDefault(GenerateReverseLookupZonesEnvName, "false"))
	if err != nil {
		return err
	}
	e.generateReverseLookupZones = generateReverseLookupZones

	generateSerial, err := strconv.ParseBool(e.envValueOrDefault(GenerateSerialEnvName, "false"))
	if err != nil {
		return err
	}
	e.generateSerial = generateSerial

	defaultPluginDir := filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "plugins")
	pluginDir, err := e.toAbsoluteFilePath(e.envValueOrDefault(PluginsDirectoryEnvName, defaultPluginDir))
	if err != nil {
		return err
	}
	e.pluginsDirectory = pluginDir

	defaultSerialChangeIndexDirectory := filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "serial")
	serialChangeIndexDirectory, err := e.toAbsoluteFilePath(e.envValueOrDefault(SerialChangeIndexDirectoryEnvName, defaultSerialChangeIndexDirectory))
	if err != nil {
		return err
	}
	e.serialChangeIndexDirectory = serialChangeIndexDirectory

	return nil
}

func (e *environment) envValueOrDefault(envName string, defaultValue string) string {
	value := os.Getenv(envName)

	if value == "" {
		value = defaultValue
	}

	return value
}

func (e *environment) toAbsoluteFilePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
