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
	"maps"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bcurnow/zonemgr/internal/plugins/builtin"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
)

type PluginManager interface {
	Plugins() map[plugins.PluginType]plugins.ZoneMgrPlugin
	Metadata() map[plugins.PluginType]*plugins.PluginMetadata
	LoadPlugins(pluginDir string) error
}

type pluginManager struct {
	PluginManager
	plugins  map[plugins.PluginType]plugins.ZoneMgrPlugin
	metadata map[plugins.PluginType]*plugins.PluginMetadata
}

var instance = &pluginManager{plugins: make(map[plugins.PluginType]plugins.ZoneMgrPlugin), metadata: make(map[plugins.PluginType]*plugins.PluginMetadata)}

func Manager() PluginManager {
	return instance
}

func (pm *pluginManager) Plugins() map[plugins.PluginType]plugins.ZoneMgrPlugin {
	return pm.plugins
}

func (pm *pluginManager) Metadata() map[plugins.PluginType]*plugins.PluginMetadata {
	return pm.metadata
}

func (pm *pluginManager) LoadPlugins(pluginDir string) error {
	maps.Copy(pm.plugins, builtin.BuiltinPlugins())
	maps.Copy(pm.metadata, builtin.BuiltinMetadata())

	if err := pm.loadExternalPlugins(pluginDir); err != nil {
		return err
	}

	return nil
}

func (pm *pluginManager) loadExternalPlugins(pluginDir string) error {
	hclog.L().Debug("Loading plugins", "pluginDir", pluginDir)
	executables, err := pm.discoverPlugins(pluginDir)
	if err != nil {
		hclog.L().Error("Error discovering plugins", "pluginDir", pluginDir, "err", err)
		return err
	}

	for pluginName, pluginCmd := range executables {
		client := pm.buildClient(pluginName, pluginCmd)
		zonemgrPlugin, err := pm.pluginInstance(pluginName, client)
		if err != nil {
			return err
		}
		supportedTypes, err := zonemgrPlugin.PluginTypes()
		if err != nil {
			return err
		}

		for _, pluginType := range supportedTypes {
			existingMetadata, ok := pm.metadata[pluginType]
			newMetadata := &plugins.PluginMetadata{Name: pluginName, Command: pluginCmd, BuiltIn: false}

			if ok {
				// We've already got on, and it's very nice
				pm.handleOverride(pluginType, existingMetadata, newMetadata)
			}

			pm.plugins[pluginType] = zonemgrPlugin
			pm.metadata[pluginType] = newMetadata
		}
	}

	return nil
}

func (pm *pluginManager) handleOverride(pluginType plugins.PluginType, existingMetadata *plugins.PluginMetadata, newMetadata *plugins.PluginMetadata) {
	// Check to see if we already have a plugin for this ResourceRecord Type
	_, ok := pm.plugins[pluginType]
	if ok {
		// If the plugin already exists then we are overriding. If what we're overriding isn't the default
		// then there are multiple plugins in the path which support the same resource record types and we should warn the user
		if existingMetadata.BuiltIn {
			hclog.L().Debug("Replacing default plugin", "pluginType", pluginType, "oldPluginName", existingMetadata.Name, "newPluginName", newMetadata.Name)
		} else {
			hclog.L().Warn("Replacing non-default plugin", "pluginType", pluginType, "oldPluginName", existingMetadata.Name, "newPluginName", newMetadata.Name)
		}
	}
}

// Walks the specified directory looking for plugins, returns an array of all the executables found
// Until goplugin.Discover is updated to check for the executable bit, this is our own implementation
func (pm *pluginManager) discoverPlugins(pluginDir string) (map[string]string, error) {
	var executables = make(map[string]string)

	hclog.L().Trace("Walking plugins dir", "dir", pluginDir)
	err := filepath.WalkDir(pluginDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				hclog.L().Trace("Could not find plugin directory", "dir", pluginDir)
				return nil
			}
			return walkErr
		}
		hclog.L().Trace("Processing path", "path", path, "dir", pluginDir)

		// Don't traverse sub-directories, this is arbitrary but we are keeping it simple
		if d.IsDir() && path != pluginDir {
			hclog.L().Trace("Subdirectories are not supported, skipping", "path", path, "dir", pluginDir)
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
				hclog.L().Trace("Skipping non-executable file", "path", path, "dir", pluginDir)
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

func (pm *pluginManager) buildClient(pluginName string, pluginCmd string) *goplugin.Client {
	hclog.L().Debug("Building a plugin client", "pluginName", pluginName, "pluginCmd", pluginCmd)

	clientConfig := &goplugin.ClientConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			pluginName: &plugins.GRPCPlugin{},
		},
		Managed:          true,                                       // Allow the plugin runtime to manage this plugin
		SyncStdout:       PluginStdout(),                             // Print any extra output to Stdout from the plugin to the host processes Stdout
		SyncStderr:       PluginStderr(),                             // Discard any any extra output to Stderr from the plugin
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC}, // We only support plugins of type grpc
		Logger:           PluginLogger(),
		SkipHostEnv:      true, // Don't pass the host environment to the plugin to avoid security issues
		AutoMTLS:         true, // Ensure that we're using MTLS for communication between the plugin and the host
		Cmd:              exec.Command(pluginCmd),
	}

	return goplugin.NewClient(clientConfig)
}

func (pm *pluginManager) pluginInstance(pluginName string, client *goplugin.Client) (plugins.ZoneMgrPlugin, error) {
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
