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

var _ plugins.ZoneMgrPlugin = &CNAMEPlugin{}

var cnameSupportedPluginTypes = []plugins.PluginType{plugins.RecordCNAME}

type CNAMEPlugin struct {
}

func (p *CNAMEPlugin) PluginVersion() (string, error) {
	return version.Version(), nil
}

func (p *CNAMEPlugin) PluginTypes() ([]plugins.PluginType, error) {
	return cnameSupportedPluginTypes, nil
}

func (p *CNAMEPlugin) Configure(config *schema.Config) error {
	// no config
	return nil
}

func (p *CNAMEPlugin) Normalize(identifier string, rr *schema.ResourceRecord) error {
	if err := plugins.StandardValidations(identifier, rr, cnameSupportedPluginTypes); err != nil {
		return err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := plugins.IsValidNameOrWildcard(rr.Name, identifier, rr); err != nil {
		return err
	}

	// Make sure the name isn't an IP
	if net.ParseIP(rr.Name) != nil {
		return fmt.Errorf("CNAME record invalid, '%s' cannot be an IP address, identifier: '%s'", rr.Name, identifier)
	}

	value, err := rr.RetrieveSingleValue(identifier)
	if err != nil {
		return err
	}
	rr.Value = value

	// Make sure the value isn't an IP
	if net.ParseIP(value) != nil {
		return fmt.Errorf("CNAME record invalid, '%s' must be a valid IP address, identifier: '%s'", rr.Value, identifier)
	}

	return nil
}

func (p *CNAMEPlugin) ValidateZone(name string, zone *schema.Zone) error {
	resourceRecordsByType := zone.ResourceRecordsByType()
	aRecords := resourceRecordsByType[schema.A]
	cnameRecords := resourceRecordsByType[schema.CNAME]

	if len(cnameRecords) > 0 && len(aRecords) == 0 {
		return fmt.Errorf("found CNAME records but there are no A records present, all CNAMES must reference an A record name, zone: '%s'", name)
	}

	for _, cnameRecord := range cnameRecords {
		_, ok := aRecords[cnameRecord.Value]
		if !ok {
			return fmt.Errorf("CNAME record '%s' has a value of '%s' which does not match any defined A record name, zone: '%s'", cnameRecord.Name, cnameRecord.Value, name)
		}
	}

	return nil
}

func (p *CNAMEPlugin) Render(identifier string, rr *schema.ResourceRecord) (string, error) {
	if err := plugins.IsSupportedPluginType(identifier, rr, cnameSupportedPluginTypes); err != nil {
		return "", err
	}
	return plugins.RenderSingleValueResource(rr), nil
}

func init() {
	registerBuiltIn(plugins.RecordCNAME, &CNAMEPlugin{})
}
