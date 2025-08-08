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
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
	goplugin "github.com/hashicorp/go-plugin"
)

var pluginTypes = []plugins.PluginType{plugins.RecordA}

// Concrete implementation of the TypeHandler
type TypeHandler struct {
}

func (th TypeHandler) PluginVersion() (string, error) {
	return "1.0.0", nil
}

func (th TypeHandler) PluginTypes() ([]plugins.PluginType, error) {
	return pluginTypes, nil
}

func (th TypeHandler) Configure(config schema.Config) error {
	// no config
	return nil
}

func (th TypeHandler) Normalize(identifier string, rr schema.ResourceRecord) (schema.ResourceRecord, error) {
	if err := plugins.StandardValidations(identifier, &rr, pluginTypes); err != nil {
		return plugins.NilResourceRecord(), err
	}

	rr.Comment = "This is an overridden comment value, it doesn't matter what is in the input file, this is the comment that will print out"

	// Make sure that there's nothing in the Values array
	if len(rr.Values) > 0 {
		// Ignore anything but the first value, move it into the Value field
		rr.Value = rr.Values[0].Value
	}

	return rr, nil
}

func (th TypeHandler) ValidateZone(name string, zone schema.Zone) error {
	// no-op
	return nil
}

func (th TypeHandler) Render(identifier string, rr schema.ResourceRecord) (string, error) {
	// Leverage the standard rendering
	return plugins.RenderSingleValueResource(&rr), nil
}

func main() {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"zonemgr-a-record-comment-override-plugin": &plugins.Plugin{Impl: &TypeHandler{}},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
	})
}
