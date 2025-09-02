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

package plugin_manager

import (
	"io"
	"os"
	"testing"
)

func TestPluginLogger(t *testing.T) {
	testCases := []struct {
		pluginDebug bool
	}{
		{true},
		{false},
	}

	for _, tc := range testCases {
		pluginDebug = tc.pluginDebug

		logger := PluginLogger()
		wanted := "plugin"
		if pluginDebug {
			//Logger will be a sub-logger
			wanted = "zonemgr.plugin"
		}

		if logger.Name() != wanted {
			t.Errorf("incorrect logger name %s, wanted %s", logger.Name(), wanted)
		}
	}
}

func TestPluginStdout(t *testing.T) {
	testCases := []struct {
		want        io.Writer
		pluginDebug bool
	}{
		{os.Stdout, true},
		{io.Discard, false},
	}
	for _, tc := range testCases {
		pluginDebug = tc.pluginDebug
		result := PluginStdout()
		if result != tc.want {
			t.Errorf("pluginStdout=%s, want %s", result, tc.want)
		}
	}
}

func TestPluginStderr(t *testing.T) {
	testCases := []struct {
		want        io.Writer
		pluginDebug bool
	}{
		{os.Stderr, true},
		{io.Discard, false},
	}
	for _, tc := range testCases {
		pluginDebug = tc.pluginDebug
		result := PluginStderr()
		if result != tc.want {
			t.Errorf("pluginStderr=%s, want %s", result, tc.want)
		}
	}
}
