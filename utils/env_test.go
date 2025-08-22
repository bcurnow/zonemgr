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

package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestDefaultValue(t *testing.T) {
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
		testEnv := &Env{Name: "TestDefaultEnv"}
		if tc.setEnv {
			if err := os.Setenv(testEnv.Name, tc.value); err != nil {
				t.Errorf("Unable to set %s=%s", testEnv.Name, tc.value)
			}
		} else {
			if err := os.Unsetenv(testEnv.Name); err != nil {
				t.Errorf("Unable to unset %s", testEnv.Name)
			}
		}
		result := defaultValue(testEnv, tc.defaultValue)
		if result != tc.want {
			t.Errorf("defaultEnv=%s, want %s", result, tc.want)
		}
	}
}

func TestInit(t *testing.T) {
	user, err := user.Current()
	if nil != err {
		t.Errorf("Could not determine current user, what did you do?")
	}

	testCases := []struct {
		e      *Env
		want   string
		setEnv bool
	}{
		{GenerateReverseLookupZones, "false", false},
		{GenerateReverseLookupZones, "TestInit" + GenerateReverseLookupZones.Name, true},
		{GenerateSerial, "false", false},
		{GenerateSerial, "TestInit" + GenerateSerial.Name, true},
		{PluginsDirectory, filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "plugins"), false},
		{PluginsDirectory, "TestInit" + PluginsDirectory.Name, true},
		{SerialChangeIndexDirectory, filepath.Join(user.HomeDir, ".local", "share", "zonemgr", "serial"), false},
		{SerialChangeIndexDirectory, "TestInit" + SerialChangeIndexDirectory.Name, true},
	}

	for _, tc := range testCases {
		if tc.setEnv {
			if err := os.Setenv(tc.e.Name, "TestInit"+tc.e.Name); err != nil {
				t.Errorf("Unable to set %s=%s", tc.e.Name, "TestInit"+tc.e.Name)
			}
		} else {
			if err := os.Unsetenv(tc.e.Name); err != nil {
				t.Errorf("Unable to unset %s", tc.e.Name)
			}
		}
		// Recall init to update the values
		defaultValues()
		if tc.want != tc.e.Value {
			t.Errorf("%s=%s, want %s", tc.e.Name, tc.e.Value, tc.want)
		}
	}
}
