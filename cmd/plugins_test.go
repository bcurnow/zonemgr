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
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/bcurnow/zonemgr/plugins"
)

func TestRunE_Plugins(t *testing.T) {
	setup(t)
	defer teardown(t)

	// The entire function is essentially logging so capture stdout (this doesn't log to stderr)
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("could not create pipe to capture stdout")
	}
	os.Stdout = w

	testPlugin.EXPECT().PluginVersion().Return("test", nil).Times(2)
	mockPluginManager.EXPECT().Plugins().Return(testPlugins)
	mockPluginManager.EXPECT().Metadata().Return(testMetadata)

	pluginsCmd.RunE(pluginsCmd, []string{})
	w.Close()

	wanted := fmt.Sprintf(pluginHeaderFormatString, pluginHeaders...) +
		fmt.Sprintf(pluginFormatString, "A", "A", "test", "A") +
		fmt.Sprintf(pluginFormatString, "NS", "NS", "test", "NS")

	var buf bytes.Buffer
	buf.ReadFrom(r)
	r.Close()

	if buf.String() != wanted {
		t.Errorf("Invalid plugin output:\n%s\nwanted:\n%s", buf.String(), wanted)
	}
}

func TestRunE_Plugins_PluginVersionError(t *testing.T) {
	setup(t)
	defer teardown(t)

	testPlugin := plugins.NewMockZoneMgrPlugin(mockController)
	// var testPlugin *plugins.MockZoneMgrPlugin
	testPlugins := make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
	testMetadata := make(map[plugins.PluginType]*plugins.Metadata)

	testPlugins[plugins.A] = testPlugin
	testPlugins[plugins.NS] = testPlugin
	testPlugin.EXPECT().PluginVersion().Return("test", nil)
	testPlugin.EXPECT().PluginVersion().Return("", errors.New("pluginVersionErr"))
	testMetadata[plugins.A] = &plugins.Metadata{Name: "A", Command: "A"}
	testMetadata[plugins.NS] = &plugins.Metadata{Name: "NS", Command: "NS"}

	mockPluginManager.EXPECT().Plugins().Return(testPlugins)
	mockPluginManager.EXPECT().Metadata().Return(testMetadata)

	if err := pluginsCmd.RunE(pluginsCmd, []string{}); err != nil {
		if err.Error() != "pluginVersionErr" {
			t.Errorf("incorrect error: '%s', want: '%s'", err, "pluginVersionErr")
		}
	} else {
		t.Error("expected an error, found none")
	}
}
