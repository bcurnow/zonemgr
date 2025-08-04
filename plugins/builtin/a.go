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

package builtin

import (
	"fmt"
	"net"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
	"github.com/bcurnow/zonemgr/version"
)

var _ plugins.TypeHandler = &APlugin{}
var aSupportedPluginTypes = []plugins.PluginType{plugins.RecordA}

type APlugin struct {
}

func (p *APlugin) PluginVersion() string {
	return version.Version
}

func (p *APlugin) PluginTypesSupported() []plugins.PluginType {
	return aSupportedPluginTypes
}

func (p *APlugin) Configure(config schema.Config) error {
	// no config
	return nil
}

func (p *APlugin) Normalize(identifier string, rr schema.ResourceRecord) (schema.ResourceRecord, error) {
	if err := plugins.StandardValidations(identifier, &rr, aSupportedPluginTypes); err != nil {
		return plugins.NilResourceRecord(), err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := plugins.IsValidNameOrWildcard(rr.Name); err != nil {
		return plugins.NilResourceRecord(), err
	}

	// Make sure the name isn't an IP
	if net.ParseIP(rr.Name) != nil {
		return plugins.NilResourceRecord(), fmt.Errorf("A record invalid, '%s' cannot be an IP address, identifier: '%s'", rr.Name, identifier)
	}

	value, err := plugins.RetrieveSingleValue(identifier, &rr)
	if err != nil {
		return plugins.NilResourceRecord(), err
	}
	rr.Value = value

	// Make sure the value IS an IP
	if net.ParseIP(value) == nil {
		return plugins.NilResourceRecord(), fmt.Errorf("A record invalid, '%s' must be a valid IP address, identifier: '%s'", rr.Value, identifier)
	}

	return rr, nil
}

func (p *APlugin) ValidateZone(name string, zone schema.Zone) error {
	//no-op
	return nil
}

func (p *APlugin) Render(identifier string, rr schema.ResourceRecord) (string, error) {
	if err := plugins.IsSupportedPluginType(identifier, &rr, aSupportedPluginTypes); err != nil {
		return "", err
	}

	return plugins.RenderSingleValueResource(&rr), nil
}

func init() {
	registerBuiltIn(plugins.RecordA, &APlugin{})
}
