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
	"fmt"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/utils"
)

var _ ZoneMgrPlugin = &BuiltinPluginCNAME{}

type BuiltinPluginCNAME struct {
	ZoneMgrPlugin
}

func (p *BuiltinPluginCNAME) PluginVersion() (string, error) {
	return utils.Version(), nil
}

func (p *BuiltinPluginCNAME) PluginTypes() ([]PluginType, error) {
	return PluginTypes(CNAME), nil
}

func (p *BuiltinPluginCNAME) Configure(config *models.Config) error {
	// no config
	return nil
}

func (p *BuiltinPluginCNAME) Normalize(identifier string, rr *models.ResourceRecord) error {
	if err := validations.StandardValidations(identifier, rr, CNAME); err != nil {
		return err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := validations.IsValidNameOrWildcard(identifier, rr.Name, rr.Type); err != nil {
		return err
	}

	// Make sure the value isn't an IP
	if err := validations.EnsureNotIP(identifier, rr.RetrieveSingleValue(), rr.Type); err != nil {
		return err
	}

	return nil
}

func (p *BuiltinPluginCNAME) ValidateZone(name string, zone *models.Zone) error {
	resourceRecordsByType := zone.ResourceRecordsByType()
	aRecords := resourceRecordsByType[models.A]
	cnameRecords := resourceRecordsByType[models.CNAME]

	if len(cnameRecords) > 0 && len(aRecords) == 0 {
		return fmt.Errorf("found CNAME records but there are no A records present, all CNAMES must reference an A record name, zone: '%s'", name)
	}

	for _, cnameRecord := range cnameRecords {
		_, ok := aRecords[cnameRecord.Value]
		if !ok {
			return fmt.Errorf("invalid CNAME record, '%s' has a value of '%s' which does not match any defined A record name, zone: '%s'", cnameRecord.Name, cnameRecord.Value, name)
		}
	}

	return nil
}

func (p *BuiltinPluginCNAME) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	if err := validations.IsSupportedPluginType(identifier, rr.Type, CNAME); err != nil {
		return "", err
	}
	return rr.RenderSingleValueResource(), nil
}

func init() {
	registerBuiltIn(CNAME, &BuiltinPluginCNAME{})
}
