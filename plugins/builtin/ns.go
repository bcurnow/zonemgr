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

// Make sure we're correctly implementing the ZonmgrPlugin interface
var _ plugins.TypeHandler = &NSPlugin{}

var nsSupportedPluginTypes = []plugins.PluginType{plugins.RecordNS}

type NSPlugin struct{}

func (p *NSPlugin) PluginVersion() (string, error) {
	return version.Version(), nil
}

func (p *NSPlugin) PluginTypes() ([]plugins.PluginType, error) {
	return nsSupportedPluginTypes, nil
}

func (p *NSPlugin) Configure(config schema.Config) error {
	// We don't need to do anything with the configuration
	return nil
}

func (p *NSPlugin) Normalize(identifier string, rr schema.ResourceRecord) (schema.ResourceRecord, error) {
	if err := plugins.StandardValidations(identifier, &rr, nsSupportedPluginTypes); err != nil {
		return plugins.NilResourceRecord(), err
	}

	// Empty names are allowed in NS record but if set, must be valid or a wild card (e.g. @)
	if rr.Name != "" {
		if err := plugins.IsValidNameOrWildcard(rr.Name); err != nil {
			return plugins.NilResourceRecord(), err
		}
	} else {
		// Default to the wildcard
		rr.Name = "@"
	}

	value, err := plugins.RetrieveSingleValue(identifier, &rr)
	if err != nil {
		return plugins.NilResourceRecord(), err
	}
	rr.Value = value

	// Check if the value is a valid name (not an IP address)
	if net.ParseIP(value) != nil {
		return plugins.NilResourceRecord(), fmt.Errorf("NS record invalid, '%s' cannot be an IP address, identifier: '%s'", value, identifier)
	}

	err = plugins.IsFullyQualified(value)
	if err != nil {
		return plugins.NilResourceRecord(), err
	}

	//Check the comment
	comment, err := plugins.RetrieveSingleComment(identifier, &rr)
	if err != nil {
		return plugins.NilResourceRecord(), err
	}
	rr.Comment = comment

	return rr, nil
}

func (p *NSPlugin) ValidateZone(name string, zone schema.Zone) error {
	//no-op
	return nil
}

func (p *NSPlugin) Render(identifier string, rr schema.ResourceRecord) (string, error) {
	if err := plugins.IsSupportedPluginType(identifier, &rr, nsSupportedPluginTypes); err != nil {
		return "", err
	}

	return plugins.RenderSingleValueResource(&rr), nil
}

func init() {
	registerBuiltIn(plugins.RecordNS, &NSPlugin{})
}
