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

var _ plugins.ZoneMgrPlugin = &APlugin{}

type APlugin struct {
	plugins.ZoneMgrPlugin
}

func (p *APlugin) PluginVersion() (string, error) {
	return version.Version(), nil
}

func (p *APlugin) PluginTypes() ([]plugins.PluginType, error) {
	return plugins.PluginTypes(plugins.A), nil
}

func (p *APlugin) Configure(config *schema.Config) error {
	// no config
	return nil
}

func (p *APlugin) Normalize(identifier string, rr *schema.ResourceRecord) error {
	if err := validations.StandardValidations(identifier, rr, plugins.A); err != nil {
		return err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := validations.IsValidNameOrWildcard(identifier, rr.Name, rr.Type); err != nil {
		return err
	}

	// Make sure the value IS an IP
	if err := validations.EnsureIP(identifier, rr.RetrieveSingleValue(), rr.Type); err != nil {
		return err
	}

	return nil
}

func (p *APlugin) ValidateZone(name string, zone *schema.Zone) error {
	//no-op
	return nil
}

func (p *APlugin) Render(identifier string, rr *schema.ResourceRecord) (string, error) {
	if err := validations.IsSupportedPluginType(identifier, rr.Type, plugins.A); err != nil {
		return "", err
	}

	return rr.RenderSingleValueResource(), nil
}

func init() {
	registerBuiltIn(plugins.A, &APlugin{})
}
