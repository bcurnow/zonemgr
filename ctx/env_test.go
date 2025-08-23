/**
 * Copyright (C) 2025 Brian Curnow
 *
 * This file is part of zonemgr.
 *
 * zonemgr is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * zonemgr is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with zonemgr.  If not, see <https://www.gnu.org/licenses/>.
 */

package ctx

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"testing"
)

func TestEnvValueOrDefault(t *testing.T) {
	testCases := []struct {
		want         string
		defaultValue string
		value        string
		setEnv       bool
	}{
		{"default", "default", "", false},
		{"env", "default", "env", true},
	}

	for _, tc := range testCases {
		envName := "TestEnvValueOrDefault"
		if tc.setEnv {
			if err := os.Setenv(envName, tc.value); err != nil {
				t.Errorf("Unable to set %s=%s", envName, tc.value)
			}
		} else {
			if err := os.Unsetenv(envName); err != nil {
				t.Errorf("Unable to unset %s", envName)
			}
		}
		env := environment{}
		result := env.envValueOrDefault(envName, tc.defaultValue)
		if result != tc.want {
			t.Errorf("defaultEnv=%s, want %s", result, tc.want)
		}
	}
}

func TestLoadValues_GenerateReverseLookupZones(t *testing.T) {
	testCases := []struct {
		setEnv bool
		want   bool
	}{
		{false, false},
		{true, false},
		{true, true},
	}

	for _, tc := range testCases {
		if tc.setEnv {
			if err := os.Setenv(GenerateReverseLookupZonesEnvName, strconv.FormatBool(tc.want)); err != nil {
				t.Errorf("Unable to set %s=%s", GenerateReverseLookupZonesEnvName, strconv.FormatBool(tc.want))
			}
		} else {
			if err := os.Unsetenv(GenerateReverseLookupZonesEnvName); err != nil {
				t.Errorf("Unable to unset %s", GenerateReverseLookupZonesEnvName)
			}
		}
		e := environment{}
		if err := e.LoadValues(); err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if tc.want != e.generateReverseLookupZones {
			t.Errorf("%s=%t, want %t", GenerateReverseLookupZonesEnvName, e.generateReverseLookupZones, tc.want)
		}
	}
}

func TestLoadValues_GenerateSerial(t *testing.T) {
	testCases := []struct {
		setEnv bool
		want   bool
	}{
		{false, false},
		{true, false},
		{true, true},
	}

	for _, tc := range testCases {
		if tc.setEnv {
			if err := os.Setenv(GenerateSerialEnvName, strconv.FormatBool(tc.want)); err != nil {
				t.Errorf("Unable to set %s=%s", GenerateSerialEnvName, strconv.FormatBool(tc.want))
			}
		} else {
			if err := os.Unsetenv(GenerateSerialEnvName); err != nil {
				t.Errorf("Unable to unset %s", GenerateSerialEnvName)
			}
		}
		e := environment{}
		if err := e.LoadValues(); err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if tc.want != e.generateSerial {
			t.Errorf("%s=%t, want %t", GenerateSerialEnvName, e.generateSerial, tc.want)
		}
	}
}
func TestLoadValues_PluginsDirectory(t *testing.T) {
	user, err := user.Current()
	if nil != err {
		t.Errorf("Could not determine current user, what did you do?")
	}

	testCases := []struct {
		setEnv bool
		want   string
	}{
		{true, "/somedir"},
		{false, filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "plugins")},
	}

	for _, tc := range testCases {
		if tc.setEnv {
			if err := os.Setenv(PluginsDirectoryEnvName, tc.want); err != nil {
				t.Errorf("Unable to set %s=%s", PluginsDirectoryEnvName, tc.want)
			}
		} else {
			if err := os.Unsetenv(PluginsDirectoryEnvName); err != nil {
				t.Errorf("Unable to unset %s", PluginsDirectoryEnvName)
			}
		}
		e := environment{}
		e.LoadValues()
		if tc.want != e.pluginsDirectory {
			t.Errorf("%s=%s, want %s", PluginsDirectoryEnvName, e.pluginsDirectory, tc.want)
		}
	}
}

func TestLoadValues_SerialChangeIndexDirectory(t *testing.T) {
	user, err := user.Current()
	if nil != err {
		t.Errorf("Could not determine current user, what did you do?")
	}

	testCases := []struct {
		setEnv bool
		want   string
	}{
		{true, "/somedir"},
		{false, filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "serial")},
	}

	for _, tc := range testCases {
		if tc.setEnv {
			if err := os.Setenv(SerialChangeIndexDirectoryEnvName, tc.want); err != nil {
				t.Errorf("Unable to set %s=%s", SerialChangeIndexDirectoryEnvName, tc.want)
			}
		} else {
			if err := os.Unsetenv(SerialChangeIndexDirectoryEnvName); err != nil {
				t.Errorf("Unable to unset %s", SerialChangeIndexDirectoryEnvName)
			}
		}
		e := environment{}
		e.LoadValues()
		if tc.want != e.serialChangeIndexDirectory {
			t.Errorf("%s=%s, want %s", SerialChangeIndexDirectoryEnvName, e.serialChangeIndexDirectory, tc.want)
		}
	}
}
