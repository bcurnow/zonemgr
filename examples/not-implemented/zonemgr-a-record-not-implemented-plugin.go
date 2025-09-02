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

package main

import (
	"errors"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	goplugin "github.com/hashicorp/go-plugin"
)

// Concrete implementation of the plugin
type Plugin struct {
}

var _ plugins.ZoneMgrPlugin = &Plugin{}

func (th *Plugin) PluginVersion() (string, error) {
	return "", errors.New("testing Plugin - Not Implemented")
}

func (th *Plugin) PluginTypes() ([]plugins.PluginType, error) {
	return nil, errors.New("testing Plugin - Not Implemented")
}

func (th *Plugin) Configure(config *models.Config) error {
	return errors.New("testing Plugin - Not Implemented")
}

func (th *Plugin) Normalize(identifier string, rr *models.ResourceRecord) error {
	return errors.New("testing Plugin - Not Implemented")
}

func (th *Plugin) ValidateZone(name string, zone *models.Zone) error {
	return errors.New("testing Plugin - Not Implemented")
}

func (th *Plugin) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	return "", errors.New("testing Plugin - Not Implemented")
}

func main() {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"zonemgr-a-record-comment-override-plugin": &plugins.GRPCPlugin{Impl: &Plugin{}},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
	})
}
