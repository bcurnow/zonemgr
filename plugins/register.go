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

package plugins

import (
	"io"
	"os"
	"os/exec"

	"github.com/bcurnow/zonemgr/env"
	goplugin "github.com/hashicorp/go-plugin"
)

func RegisterPlugins() (map[string]TypeHandler, error) {
	logger.Debug("Loading plugins", "pluginDir", env.PLUGINS.Value)
	typeHandlers := make(map[string]TypeHandler)
	executables, err := discoverPlugins(env.PLUGINS.Value)
	if err != nil {
		return nil, err
	}

	for pluginName, pluginCmd := range executables {
		clientConfig := buildClientConfig(pluginName, pluginCmd)
		client := buildClient(pluginName, clientConfig)
		typeHandler, err := pluginInstance(pluginName, client)
		if err != nil {
			return nil, err
		}
		typeHandlers[pluginName] = typeHandler
	}

	return typeHandlers, nil
}

func buildClientConfig(pluginName string, pluginCmd string) *goplugin.ClientConfig {
	logger.Debug("Building a plugin client", "pluginName", pluginName, "pluginCmd", pluginCmd)

	return &goplugin.ClientConfig{
		HandshakeConfig:  HandshakeConfig,
		Plugins:          PluginMap,
		Managed:          true,                                       // Allow the plugin runtime to manage this plugin
		SyncStdout:       os.Stdout,                                  // Print any extra output to Stdout from the plugin to the host processes Stdout
		SyncStderr:       io.Discard,                                 // Discard any any extra output to Stderr from the plugin
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC}, // We only support plugins of type grpc
		Logger:           logger,
		SkipHostEnv:      true, // Don't pass the host environment to the plugin to avoid security issues
		AutoMTLS:         true, // Ensure that we're using MTLS for communication between the plugin and the host
		Cmd:              exec.Command(pluginCmd),
	}
}

func buildClient(pluginName string, config *goplugin.ClientConfig) *goplugin.Client {
	logger.Trace("Creating new client", "pluginName", pluginName)
	client := goplugin.NewClient(config)
	logger.Debug("Plugin client created", "Protocol", client.Protocol())
	return client
}

func pluginInstance(pluginName string, client *goplugin.Client) (TypeHandler, error) {
	logger.Trace("Getting the ClientProtocol from the client", "pluginName`", pluginName)
	// Get the RPC Client from the plugin definition
	clientProtocol, err := client.Client()
	if err != nil {
		return nil, err
	}

	logger.Trace("Dispensing plugin", "pluginName", pluginName)
	// Get the actual client so we can use it
	raw, err := clientProtocol.Dispense(pluginName)
	if err != nil {
		return nil, err
	}

	// Cast the raw plugin to the TypeHandler interface so we have access to the methods
	return raw.(TypeHandler), nil
}
