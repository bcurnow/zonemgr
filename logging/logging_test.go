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

package logging

import (
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
)

func TestEnablePluginDebug(t *testing.T) {
	testCases := []struct {
		initialValue bool
		want         bool
	}{
		{true, true},
		{false, true},
	}

	for _, tc := range testCases {
		pluginDebug = tc.initialValue
		EnablePluginDebug()
		if pluginDebug != tc.want {
			t.Errorf("pluginDebug=%s, want %s", strconv.FormatBool(pluginDebug), strconv.FormatBool(tc.want))
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

func TestConfigureLogging(t *testing.T) {
	testCases := []struct {
		resetTo *hclog.LoggerOptions
		want    *hclog.LoggerOptions
	}{
		{testLoggingOptions(hclog.NoLevel, true, false, false), testLoggingOptions(hclog.Error, false, true, true)},
		{testLoggingOptions(hclog.NoLevel, false, true, true), testLoggingOptions(hclog.Warn, true, false, false)},
	}
	for _, tc := range testCases {
		hclog.DefaultOptions.Name = tc.resetTo.Name
		hclog.DefaultOptions.Output = tc.resetTo.Output
		hclog.DefaultOptions.Color = tc.resetTo.Color
		hclog.DefaultOptions.Level = tc.resetTo.Level
		hclog.DefaultOptions.JSONFormat = tc.resetTo.JSONFormat
		hclog.DefaultOptions.DisableTime = tc.resetTo.DisableTime

		ConfigureLogging(tc.want.Level, tc.want.JSONFormat, tc.want.DisableTime, toLogColor(tc.want.Color))
		if diff := cmp.Diff(hclog.DefaultOptions, tc.want, loggerOptsCompare()); diff != "" {
			t.Errorf("Incorrect DefaultOptions:\n%s", diff)
		}

	}
}

func TestDefaultLogging(t *testing.T) {
	hclog.DefaultOptions.Name = "notDefault"
	hclog.DefaultOptions.Output = os.Stdin
	hclog.DefaultOptions.Color = hclog.ForceColor
	hclog.DefaultOptions.Level = hclog.NoLevel
	hclog.DefaultOptions.JSONFormat = true
	hclog.DefaultOptions.DisableTime = false

	defaultLogging()
	want := testLoggingOptions(hclog.Info, false, true, true)
	want.Name = "zonemgr"
	want.Output = os.Stderr

	if diff := cmp.Diff(hclog.DefaultOptions, want, loggerOptsCompare()); diff != "" {
		t.Errorf("Incorrect DefaultOptions:\n%s", diff)
	}

}

func testLoggingOptions(level hclog.Level, jsonFormat bool, disableTime bool, logColor bool) *hclog.LoggerOptions {
	opts := &hclog.LoggerOptions{
		Name:   "zonemgr-test",
		Output: os.Stderr,
	}

	if logColor {
		opts.Color = hclog.AutoColor
	} else {
		opts.Color = hclog.ColorOff
	}

	opts.Level = level
	opts.JSONFormat = jsonFormat
	opts.DisableTime = disableTime

	return opts
}

func toLogColor(color hclog.ColorOption) bool {
	return color == hclog.AutoColor
}

func loggerOptsCompare() cmp.Option {
	return cmp.Comparer(func(f1, f2 *hclog.LoggerOptions) bool {
		result := true
		result = result && f1.Color == f2.Color
		result = result && f1.Level == f2.Level
		result = result && f1.JSONFormat == f2.JSONFormat
		result = result && f1.DisableTime == f2.DisableTime
		return result
	})
}
