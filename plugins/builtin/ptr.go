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

var _ plugins.ZoneMgrPlugin = &PTRPlugin{}

type PTRPlugin struct {
	plugins.ZoneMgrPlugin
}

func (p *PTRPlugin) PluginVersion() (string, error) {
	return version.Version(), nil
}

func (p *PTRPlugin) PluginTypes() ([]plugins.PluginType, error) {
	return plugins.PluginTypes(plugins.PTR), nil
}

func (p *PTRPlugin) Configure(config *schema.Config) error {
	return nil
}

func (p *PTRPlugin) Normalize(identifier string, rr *schema.ResourceRecord) error {
	if err := validations.StandardValidations(identifier, rr, plugins.PTR); err != nil {
		return err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := validations.IsValidNameOrWildcard(identifier, rr.Name, rr.Type); err != nil {
		return err
	}

	if err := validations.IsFullyQualified(identifier, rr.RetrieveSingleValue(), rr.Type); err != nil {
		return err
	}

	return nil
}

func (p *PTRPlugin) ValidateZone(name string, zone *schema.Zone) error {
	return nil
}

func (p *PTRPlugin) Render(identifier string, rr *schema.ResourceRecord) (string, error) {
	if err := validations.IsSupportedPluginType(identifier, rr.Type, plugins.PTR); err != nil {
		return "", err
	}

	return rr.RenderSingleValueResource(), nil
}

func init() {
	registerBuiltIn(plugins.PTR, &PTRPlugin{})
}
