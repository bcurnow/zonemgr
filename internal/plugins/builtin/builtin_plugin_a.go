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
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
)

var _ plugins.ZoneMgrPlugin = &BuiltinPluginA{}

type BuiltinPluginA struct {
	plugins.ZoneMgrPlugin
}

func (p *BuiltinPluginA) PluginVersion() (string, error) {
	return utils.Version(), nil
}

func (p *BuiltinPluginA) PluginTypes() ([]plugins.PluginType, error) {
	return plugins.PluginTypes(plugins.A), nil
}

func (p *BuiltinPluginA) Configure(config *models.Config) error {
	// no config
	return nil
}

func (p *BuiltinPluginA) Normalize(identifier string, rr *models.ResourceRecord) error {
	if err := validations.CommonValidations(identifier, rr, plugins.A); err != nil {
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

func (p *BuiltinPluginA) ValidateZone(name string, zone *models.Zone) error {
	//no-op
	return nil
}

func (p *BuiltinPluginA) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	if err := validations.IsSupportedPluginType(identifier, rr.Type, plugins.A); err != nil {
		return "", err
	}

	return rr.RenderSingleValueResource(), nil
}

func init() {
	registerBuiltIn(plugins.A, &BuiltinPluginA{})
}
