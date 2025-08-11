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

package manager

import (
	"maps"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bcurnow/zonemgr/env"
	"github.com/bcurnow/zonemgr/logging"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/builtin"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
)

var (
	allPlugins  = make(map[plugins.PluginType]*plugins.Plugin)
	initialized = false
)

func Plugins() (map[plugins.PluginType]*plugins.Plugin, error) {
	if err := initializePlugins(); err != nil {
		return nil, err
	}

	return allPlugins, nil
}

func initializePlugins() error {
	if initialized {
		return nil
	}

	maps.Copy(allPlugins, builtin.BuiltinPlugins())

	externalPlugins, err := registerPlugins()
	if err != nil {
		return err
	}

	// We could just copy the externalPlugins map ove the allPlugins map and everything would be fine
	// However, iterating gives us better diagnostic logs
	for name, externalPlugin := range externalPlugins {
		// Get the list of resource record types the external plugin supports
		pluginTypes, err := externalPlugin.Plugin.PluginTypes()
		if err != nil {
			return err
		}

		for _, pluginType := range pluginTypes {
			handleOverride(pluginType, name)
			allPlugins[pluginType] = externalPlugin
		}
	}
	initialized = true
	return nil
}

func handleOverride(pluginType plugins.PluginType, pluginName string) {
	// Check to see if we already have a plugin for this ResourceRecord Type
	plugin, ok := allPlugins[pluginType]
	if ok {
		// If the plugin already exists then we are overriding. If what we're overriding isn't the default
		// then there are multiple plugins in the path which support the same resource record types and we should warn the user
		if plugin.IsBuiltIn {
			hclog.L().Debug("Replacing default plugin", "resourceRecordType", pluginType, "newPluginName", pluginName)
		} else {
			hclog.L().Warn("Replacing non-default plugin", "resourceRecordType", pluginType, "oldPluginName", plugin.PluginName, "newPluginName", pluginName, "pluginDir", env.PLUGINS.Value)
		}
	}
}

// Walks the specified directory looking for plugins, returns an array of all the executables found
// Until goplugin.Discover is updated to check for the executable bit, this is our own implementation
func discoverPlugins(dir string) (map[string]string, error) {
	var executables = make(map[string]string)

	hclog.L().Trace("Walking plugins dir", "dir", dir)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				hclog.L().Trace("Could not find plugin directory", "dir", dir)
				return nil
			}
			return walkErr
		}
		hclog.L().Trace("Processing path", "path", path, "dir", dir)

		// Don't traverse sub-directories, this is arbitrary but we are keeping it simple
		if d.IsDir() && path != dir {
			hclog.L().Trace("Subdirectories are not supported, skipping", "path", path, "dir", dir)
			return filepath.SkipDir
		}

		// Because we're using WalkDir, we need to get the FileInfo from the DirEntry
		info, err := d.Info()
		if err != nil {
			return err
		}

		// Check if this is a file and if the file is executable
		if info.Mode().IsRegular() {
			// 0111 checks for the execute bit to be set
			if info.Mode()&0111 == 0 {
				hclog.L().Trace("Skipping non-executable file", "path", path, "dir", dir)
				return nil
			}

			// Get the absolute path of the file so we can provide the best debugging information
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			executables[filepath.Base(path)] = absPath
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return executables, nil
}

func registerPlugins() (map[string]*plugins.Plugin, error) {
	hclog.L().Debug("Loading plugins", "pluginDir", env.PLUGINS.Value)
	typeHandlers := make(map[string]*plugins.Plugin)
	executables, err := discoverPlugins(env.PLUGINS.Value)
	if err != nil {
		hclog.L().Error("Error discovering plugins", "pluginDir", env.PLUGINS.Value, "err", err)
		return nil, err
	}

	for pluginName, pluginCmd := range executables {
		client := buildClient(pluginName, pluginCmd)
		typeHandler, err := pluginInstance(pluginName, client)
		if err != nil {
			return nil, err
		}
		typeHandlers[pluginName] = &plugins.Plugin{IsBuiltIn: false, PluginName: pluginName, PluginCmd: pluginCmd, Plugin: typeHandler}
	}

	return typeHandlers, nil
}

func buildClient(pluginName string, pluginCmd string) *goplugin.Client {
	hclog.L().Debug("Building a plugin client", "pluginName", pluginName, "pluginCmd", pluginCmd)

	clientConfig := &goplugin.ClientConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			pluginName: &plugins.GRPCPlugin{},
		},
		Managed:          true,                                       // Allow the plugin runtime to manage this plugin
		SyncStdout:       logging.PluginStdout(),                     // Print any extra output to Stdout from the plugin to the host processes Stdout
		SyncStderr:       logging.PluginStderr(),                     // Discard any any extra output to Stderr from the plugin
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC}, // We only support plugins of type grpc
		Logger:           logging.PluginLogger(),
		SkipHostEnv:      true, // Don't pass the host environment to the plugin to avoid security issues
		AutoMTLS:         true, // Ensure that we're using MTLS for communication between the plugin and the host
		Cmd:              exec.Command(pluginCmd),
	}

	return goplugin.NewClient(clientConfig)
}

func pluginInstance(pluginName string, client *goplugin.Client) (plugins.ZoneMgrPlugin, error) {
	hclog.L().Trace("Getting the ClientProtocol from the client", "pluginName", pluginName)
	// Get the RPC Client from the plugin definition
	clientProtocol, err := client.Client()
	if err != nil {
		return nil, err
	}

	hclog.L().Trace("Dispensing plugin", "pluginName", pluginName)
	// Get the actual client so we can use it
	raw, err := clientProtocol.Dispense(pluginName)
	if err != nil {
		return nil, err
	}
	hclog.L().Debug("Plugin dispensed", "pluginName", pluginName, "Protocol", client.Protocol())

	// Cast the raw plugin to the TypeHandler interface so we have access to the methods
	return raw.(plugins.ZoneMgrPlugin), nil
}
