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

package cmd

import (
	"errors"
	"testing"

	"github.com/bcurnow/zonemgr/plugins/plugin_manager"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/viper"
)

func TestPersistentPostRun_Root(t *testing.T) {
	defer func() { cleanupClients = goplugin.CleanupClients }()
	called := false
	cleanupClients = func() {
		called = true
	}

	rootCmd.PersistentPostRun(rootCmd, []string{})

	if !called {
		t.Errorf("expected cleanupClients to be called, it was not")
	}
}

func TestSetupLoging(t *testing.T) {
	testCases := []struct {
		logLevelString string
		wantedLogLevel hclog.Level
		logJson        bool
		logTime        bool
		logColor       bool
	}{
		{logLevelString: "bogus", wantedLogLevel: hclog.Info},
		{logLevelString: "warn", wantedLogLevel: hclog.Warn},
		{wantedLogLevel: hclog.Info, logJson: true},
		{wantedLogLevel: hclog.Info, logTime: true},
		{wantedLogLevel: hclog.Info, logColor: true},
		{wantedLogLevel: hclog.Info, logJson: true, logTime: true, logColor: true},
	}

	for _, tc := range testCases {
		v = viper.New()
		v.Set("log-level", tc.logLevelString)
		v.Set("log-json", tc.logJson)
		v.Set("log-time", tc.logTime)
		v.Set("log-color", tc.logColor)

		setupLogging()

		if hclog.DefaultOptions.Level != tc.wantedLogLevel {
			t.Errorf("incorrect log level: '%s', want: '%s'", hclog.DefaultOptions.Level, tc.wantedLogLevel)
		}

		if hclog.DefaultOptions.JSONFormat != tc.logJson {
			t.Errorf("incorrect JSONFormat: '%v', want: '%v'", hclog.DefaultOptions.JSONFormat, tc.logJson)
		}

		want := hclog.ColorOff
		if tc.logColor {
			want = hclog.AutoColor
		}
		if hclog.DefaultOptions.Color != want {
			t.Errorf("incorrect Color: '%v', want: '%v'", hclog.DefaultOptions.Color, want)
		}

		if hclog.DefaultOptions.DisableTime != !tc.logTime {
			t.Errorf("incorrect DisableTime: '%v', want: '%v'", hclog.DefaultOptions.DisableTime, !tc.logTime)
		}
	}
}

func TestPersistentPreRunE_Root(t *testing.T) {
	setup(t)
	defer teardown(t)
	testCases := []struct {
		initErr          bool
		pluginManagerErr bool
		pluginDebug      bool
	}{
		{},
		{pluginDebug: true},
		{initErr: true},
		{pluginManagerErr: true},
	}

	for _, tc := range testCases {
		call := mockFs.EXPECT().ToAbsoluteFilePath("testing")
		if tc.initErr {
			call.Return("", errors.New("initErr"))
		} else {
			call.Return("testing", nil)
			call = mockPluginManager.EXPECT().LoadPlugins("testing")
			if tc.pluginManagerErr {
				call.Return(errors.New("pluginManagerErr"))
			} else {
				call.Return(nil)
			}
		}

		args := make([]string, 2)
		args[0] = "--plugin-dir"
		args[1] = "testing"
		if tc.pluginDebug {
			args = append(args, "--plugin-debug")
		}
		rootCmd.ParseFlags(args)
		if err := rootCmd.PersistentPreRunE(rootCmd, []string{}); err != nil {
			want := ""
			if tc.initErr {
				want = "initErr"
			} else if tc.pluginManagerErr {
				want = "pluginManagerErr"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.initErr || tc.pluginManagerErr {
				t.Error("expected an error, found none")
			} else {

				if v.GetEnvPrefix() != "zonemgr" {
					t.Errorf("incorrect viper env prefix: '%s', want: '%s'", v.GetEnvPrefix(), "zonemgr")
				}

				if v.GetString("plugin-dir") != "testing" {
					t.Errorf("incorrect plugin-dir: '%s', want: '%s'", v.GetString("plugin-dir"), "testing")
				}

				logger := plugin_manager.PluginLogger()
				loggerName := "plugin"
				if tc.pluginDebug {
					loggerName = "zonemgr.plugin"
				}

				if logger.Name() != loggerName {
					t.Errorf("plugin debug not set to the correct value: %v, want: %v", logger.Name(), loggerName)
				}
			}
		}
	}
}

func TestExecute(t *testing.T) {
	// This test is purely for coverage
	// This should simply print the help and exit
	Execute()
}
