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
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/bcurnow/zonemgr/internal/mocks"
	"github.com/bcurnow/zonemgr/internal/plugins/builtin"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
)

var (
	pluginDir      string
	mockController *gomock.Controller
	mockFs         *mocks.MockFileSystem
)

func pluginManagerSetup(t *testing.T) {
	mockController = gomock.NewController(t)
	mockFs = mocks.NewMockFileSystem(mockController)
	fs = mockFs

	// Make sure that the metadata and plugins maps are empty before we start
	instance.plugins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
	instance.metadata = make(map[plugins.PluginType]*plugins.PluginMetadata)

	//Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Errorf("unable to create temp directory for testing: %s", err)
	}
	pluginDir = tempDir
}

func pluginManagerTearDown(t *testing.T) {
	if err := os.RemoveAll(pluginDir); err != nil {
		t.Errorf("unable to cleanup, unexpected error: %s", err)
	}
	mockController.Finish()
}
func TestManager(t *testing.T) {
	pluginManagerSetup(t)
	defer pluginManagerTearDown(t)
	one := Manager()
	two := Manager()

	if one != instance {
		t.Errorf("expected result to be the instance, was: %v", one)
	}

	if two != instance {
		t.Errorf("expected result to be the instance, was: %v", one)
	}

	if one != two {
		t.Errorf("expected a single instance to be return, two calls returned different instances")
	}
}

func TestPlugins(t *testing.T) {
	pluginManagerSetup(t)
	defer pluginManagerTearDown(t)
	pm := &pluginManager{}

	pm.plugins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
	wanted := &builtin.BuiltinPluginPTR{}
	pm.plugins["custom key"] = wanted

	pluginsFromPM := pm.Plugins()

	if len(pluginsFromPM) != 1 {
		t.Errorf("expected as single value, found %d", len(pluginsFromPM))
	}

	p, ok := pluginsFromPM["custom key"]

	if !ok {
		t.Errorf("expected key %s, did not find", "custom key")
	}

	if p != wanted {
		t.Errorf("incorrect value for key %s, found %v, wanted %v", "custom key", p, wanted)
	}
}

func TestMetadata(t *testing.T) {
	pluginManagerSetup(t)
	defer pluginManagerTearDown(t)

	wanted := &plugins.PluginMetadata{Name: "testing", Command: "testing command", BuiltIn: false}
	instance.metadata["custom key"] = wanted

	metadataFromPM := instance.Metadata()

	if len(metadataFromPM) != 1 {
		t.Errorf("expected as single value, found %d", len(metadataFromPM))
	}

	p, ok := metadataFromPM["custom key"]

	if !ok {
		t.Errorf("expected key %s, did not find", "custom key")
	}

	if p != wanted {
		t.Errorf("incorrect value for key %s, found %v, wanted %v", "custom key", p, wanted)
	}
}

func TestLoadPlugins(t *testing.T) {
	testCases := []struct {
		walkExecutablesErr bool
	}{
		{walkExecutablesErr: true},
		{walkExecutablesErr: false},
	}
	pluginManagerSetup(t)
	defer pluginManagerTearDown(t)

	for _, tc := range testCases {
		call := mockFs.EXPECT().WalkExecutables(t.Name(), false)
		if tc.walkExecutablesErr {
			call.Return(nil, errors.New("testing"))
		} else {
			call.Return(map[string]string{}, nil)
		}
		if err := Manager().LoadPlugins(t.Name()); err != nil {
			if tc.walkExecutablesErr {
				if err.Error() != "testing" {
					t.Errorf("incorrect error: '%s', want 'testing'", err)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if tc.walkExecutablesErr {
				t.Errorf("expected error, found none")
			} else {
				// Make sure that we only have the builtins in the maps
				if len(Manager().Plugins()) != 5 {
					t.Errorf("expected only builtins to be loaded, found %d plugins", len(Manager().Plugins()))
				}

				if len(Manager().Metadata()) != 5 {
					t.Errorf("expected only builtins to be loaded, found %d metadata", len(Manager().Plugins()))
				}
			}
		}
	}

}

func TestLoadExternalPlugins(t *testing.T) {
	testCases := []struct {
		pluginDir             string
		walkExecutablesErr    bool
		walkExecutablesResult map[string]string
		wantErr               string
		realFs                bool
		wantedPluginCount     int
	}{
		{pluginDir: "../examples/bin/comment-override", realFs: true, wantedPluginCount: 1},
		{pluginDir: "walk-executable-error", walkExecutablesErr: true},
		{pluginDir: "plugin-instance-error", walkExecutablesResult: map[string]string{"does-not-exist": "does-not-exist"}, wantErr: "exec: \"does-not-exist\": executable file not found in $PATH"},
		{pluginDir: "../examples/bin/not-implemented", realFs: true, wantErr: "rpc error: code = Unknown desc = testing Plugin - Not Implemented"},
	}
	pluginManagerSetup(t)
	defer pluginManagerTearDown(t)

	for _, tc := range testCases {
		if tc.realFs {
			fs = utils.FS()
		} else {
			fs = mockFs
			call := mockFs.EXPECT().WalkExecutables(tc.pluginDir, false)
			if tc.walkExecutablesErr {
				call.Return(nil, errors.New("testing"))
			} else {
				call.Return(tc.walkExecutablesResult, nil)
			}
		}

		if err := instance.loadExternalPlugins(tc.pluginDir); err != nil {
			wantErr := tc.wantErr
			if tc.walkExecutablesErr {
				wantErr = "testing"
			}
			if wantErr != "" {
				if err.Error() != wantErr {
					t.Errorf("incorrect error: '%s', want '%s'", err, wantErr)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if len(instance.plugins) != tc.wantedPluginCount {
				t.Errorf("expected %d plugin(s), found %d", tc.wantedPluginCount, len(instance.plugins))
			}
		}
	}
}

func TestHandleOverride(t *testing.T) {
	testCases := []struct {
		havePlugin bool
		builtin    bool
		want       string
	}{
		{havePlugin: false},
		{havePlugin: true, builtin: false, want: "Replacing non-default plugin"},
		{havePlugin: true, builtin: true, want: "Replacing default plugin"},
	}

	for _, tc := range testCases {
		//Create a pipe for logging, all we're actually testing here is the logging
		r, w, originalLogger := captureLogging(t)
		// Setup the plugins map
		instance.plugins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
		if tc.havePlugin {
			// This doesn't matter except to have value
			instance.plugins[plugins.A] = &builtin.BuiltinPluginA{}
		}

		existingMetadata := &plugins.PluginMetadata{Name: "Testing", Command: "Testing Command", BuiltIn: tc.builtin}
		newMetadata := &plugins.PluginMetadata{Name: "New", Command: "New Command", BuiltIn: false}

		// Turn up the log level
		instance.handleOverride(plugins.A, existingMetadata, newMetadata)

		// Get the output
		w.Close()
		restoreLogging(originalLogger)
		output := getOutput(r)

		if tc.want != "" {
			if !strings.Contains(output, tc.want) {
				t.Errorf("incorrect output:\n%s\nwant it to contain: '%s'", output, tc.want)
			}
		} else {
			if output != "" {
				t.Errorf("unexpected output:\n%s\nwanted none", output)
			}
		}
	}
}

func TestPluginInstance_DispenseError(t *testing.T) {
	pluginManagerSetup(t)
	defer pluginManagerTearDown(t)

	// To get a dispense error, we need to setup a client such that
	// the plugin we request doesn't exist
	clientConfig := &goplugin.ClientConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"does-not-exist": &plugins.GRPCPlugin{},
		},
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC}, // We only support plugins of type grpc
		Cmd:              exec.Command("../examples/bin/comment-override/zonemgr-a-record-comment-override-plugin"),
	}
	client := goplugin.NewClient(clientConfig)

	if _, err := instance.pluginInstance("plugin-that-should-exist", client); err != nil {
		wanted := "unknown plugin type: plugin-that-should-exist"
		if err.Error() != wanted {
			t.Errorf("incorrect error: '%s', wanted: '%s'", err, wanted)
		}
	} else {
		t.Error("expected error, found none")
	}
}

func captureLogging(t *testing.T) (*os.File, *os.File, hclog.Logger) {
	//Create a pipe for logging, all we're actually testing here is the logging
	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("could not create pipe to capture stdout/stderr")
	}

	// Override the default logger
	originalDefault := hclog.Default()
	hclog.SetDefault(hclog.New(&hclog.LoggerOptions{
		Output: w,
		Level:  hclog.Trace,
	}))

	return r, w, originalDefault
}

func getOutput(r *os.File) string {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func restoreLogging(l hclog.Logger) {
	hclog.SetDefault(l)
}
