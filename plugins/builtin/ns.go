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
	"github.com/bcurnow/zonemgr/plugins"

	"github.com/bcurnow/zonemgr/schema"
	"github.com/bcurnow/zonemgr/version"
)

// Make sure we're correctly implementing the ZonmgrPlugin interface
var _ plugins.ZoneMgrPlugin = &NSPlugin{}

type NSPlugin struct {
	plugins.ZoneMgrPlugin
}

func (p *NSPlugin) PluginVersion() (string, error) {
	return version.Version(), nil
}

func (p *NSPlugin) PluginTypes() ([]plugins.PluginType, error) {
	return plugins.PluginTypes(plugins.NS), nil
}

func (p *NSPlugin) Configure(config *schema.Config) error {
	// We don't need to do anything with the configuration
	return nil
}

func (p *NSPlugin) Normalize(identifier string, rr *schema.ResourceRecord) error {
	if err := validations.StandardValidations(identifier, rr, plugins.NS); err != nil {
		return err
	}

	// Empty names are allowed in NS record but if set, must be valid or a wild card (e.g. @)
	if rr.Name == "" {
		// Default to the wildcard
		rr.Name = "@"
	}

	if err := validations.IsValidNameOrWildcard(identifier, rr.Name, rr.Type); err != nil {
		return err
	}

	if rr.RetrieveSingleValue() == "" {
		rr.Value = identifier
	}

	// Check if the value is a valid name (not an IP address)
	if err := validations.EnsureIP(identifier, rr.RetrieveSingleValue(), rr.Type); err != nil {
		return err
	}

	if err := validations.IsFullyQualified(identifier, rr.RetrieveSingleValue(), rr.Type); err != nil {
		return err
	}

	return nil
}

func (p *NSPlugin) ValidateZone(name string, zone *schema.Zone) error {
	//no-op
	return nil
}

func (p *NSPlugin) Render(identifier string, rr *schema.ResourceRecord) (string, error) {
	if err := validations.IsSupportedPluginType(identifier, rr.Type, plugins.NS); err != nil {
		return "", err
	}

	return rr.RenderSingleValueResource(), nil
}

func init() {
	registerBuiltIn(plugins.NS, &NSPlugin{})
}
