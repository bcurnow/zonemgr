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

	"github.com/bcurnow/zonemgr/dns/serial"
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
)

var (
	_                   plugins.ZoneMgrPlugin = &BuiltinPluginSOA{}
	serialIndexManager  serial.SerialManager
	soaValuesNormalizer plugins.SOAValuesNormalizer
)

type BuiltinPluginSOA struct {
	plugins.ZoneMgrPlugin
	config *models.Config
}

func (p *BuiltinPluginSOA) PluginVersion() (string, error) {
	return utils.Version(), nil
}

func (p *BuiltinPluginSOA) PluginTypes() ([]plugins.PluginType, error) {
	return plugins.PluginTypes(plugins.SOA), nil
}

func (p *BuiltinPluginSOA) Configure(config *models.Config) error {
	p.config = config
	serialIndexManager = serial.FileSerialManager(p.config.SerialChangeIndexDirectory)
	soaValuesNormalizer = plugins.SVN()
	return nil
}

func (p *BuiltinPluginSOA) Normalize(identifier string, rr *models.ResourceRecord) error {
	if err := validations.CommonValidations(identifier, rr, plugins.SOA); err != nil {
		return err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := validations.IsFullyQualified(identifier, rr.Name, rr.Type); err != nil {
		return err
	}

	// Validate the the Value and Comment fields are empty, there are no shortcuts for SOA records
	if rr.Value != "" {
		return fmt.Errorf("value field cannot be used on SOA records, please use the values field, identifier: '%s'", identifier)
	}

	if rr.Comment != "" {
		return fmt.Errorf("comment field cannot be used on SOA records, please use the values field, identifier: '%s'", identifier)
	}

	// Only generate the next serial number if the config option is set
	generatedSerial := ""
	if p.config.GenerateSerial {
		// The name of the SOA record is also the name of the zone
		nextSerial, err := serialIndexManager.Next(rr.Name)
		if err != nil {
			return err
		}
		generatedSerial = nextSerial
	}

	if err := soaValuesNormalizer.Normalize(identifier, rr, validations, p.config.GenerateSerial, generatedSerial); err != nil {
		return err
	}

	return nil
}

func (p *BuiltinPluginSOA) ValidateZone(name string, zone *models.Zone) error {
	resourceRecords := zone.ResourceRecordsByType()
	soaRecords, ok := resourceRecords[models.SOA]
	if !ok {
		return fmt.Errorf("invalid zone, missing SOA record, zone=%s", name)
	}
	if len(soaRecords) > 1 {
		return fmt.Errorf("more than one SOA record found, only one SOA record is allowed, zone=%s", name)
	}

	return nil
}

func (p *BuiltinPluginSOA) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	if err := validations.IsSupportedPluginType(identifier, rr.Type, plugins.SOA); err != nil {
		return "", err
	}

	return rr.RenderMultivalueResource(), nil
}

func init() {
	registerBuiltIn(plugins.SOA, &BuiltinPluginSOA{})
}
