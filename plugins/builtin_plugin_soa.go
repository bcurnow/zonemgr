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
	"github.com/hashicorp/go-hclog"
)

const generatedSerialNumberComment = "Zonemgr generated serial number"

var (
	_                  ZoneMgrPlugin            = &BuiltinPluginSOA{}
	serialIndexManager utils.SerialIndexManager = &utils.FileSerialIndexManager{}
)

type BuiltinPluginSOA struct {
	ZoneMgrPlugin

	config *models.Config
}

func (p *BuiltinPluginSOA) PluginVersion() (string, error) {
	return utils.Version(), nil
}

func (p *BuiltinPluginSOA) PluginTypes() ([]PluginType, error) {
	return PluginTypes(SOA), nil
}

func (p *BuiltinPluginSOA) Configure(config *models.Config) error {
	p.config = config
	return nil
}

func (p *BuiltinPluginSOA) Normalize(identifier string, rr *models.ResourceRecord) error {
	if err := validations.StandardValidations(identifier, rr, SOA); err != nil {
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
		nextSerial, err := serialIndexManager.GetNext(rr.Name)
		if err != nil {
			return err
		}
		generatedSerial = nextSerial
	}

	if err := normalizeValues(identifier, rr, p.config.GenerateSerial, generatedSerial); err != nil {
		return err
	}

	return nil
}

func (p *BuiltinPluginSOA) ValidateZone(name string, zone *models.Zone) error {
	hasSOA := false
	for identifier, rr := range zone.ResourceRecords {
		// Track the SOA records, there can be only one
		if rr.Type == models.SOA {
			if hasSOA {
				return fmt.Errorf("more than one SOA record found, only one SOA record is allowed, identifier:'%s'", identifier)
			}
			hasSOA = true
		}
	}

	if !hasSOA {
		return fmt.Errorf("invalid zone, missing SOA record, zone=%s", name)
	}

	return nil
}

func (p *BuiltinPluginSOA) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	if err := validations.IsSupportedPluginType(identifier, rr.Type, SOA); err != nil {
		return "", err
	}

	return rr.RenderMultivalueResource(), nil
}

// Checks for the various required configurations and then normalizes the values to include a serial number (if necessary) in proper location
func normalizeValues(identifier string, rr *models.ResourceRecord, generateSerial bool, serial string) error {
	numValues := len(rr.Values)

	if numValues != 6 && numValues != 7 {
		return fmt.Errorf("incorrect number of values for the SOA record, expected 6 (no serial) or 7, found %d", numValues)
	}

	if generateSerial && numValues != 6 {
		hclog.L().Debug("Ignoring serial number of SOA record, using generated one", "identifier", identifier, "serialNumber", rr.Values[2].Value, "generateSerial", generateSerial, "generatedSerialNumber", serial)
	}

	if !generateSerial && numValues == 6 {
		return fmt.Errorf("must specify a serial number when generate serial is false, found only 6 values when 7 are required, identifier: '%s', name: '%s'", identifier, rr.Name)
	}

	if numValues == 6 {
		insertSerial(rr, serial)
	}

	if numValues == 7 && generateSerial {
		// Replace element 2
		rr.Values[2] = &models.ResourceRecordValue{Value: serial, Comment: generatedSerialNumberComment}
	}

	// The only other option is that we have 7 values and we don't generate the serial number so there's nothing to do

	// Now, check the values
	if err := validations.IsFullyQualified(identifier, rr.Values[0].Value, rr.Type); err != nil {
		return err
	}

	email, err := validations.FormatEmail(identifier, rr.Values[1].Value, rr.Type)
	if err != nil {
		return err
	}

	rr.Values[1].Value = email

	//Make sure none of the other values are < 0
	if err := validations.IsPositive(identifier, rr.Values[3].Value, "REFRESH", rr.Type); err != nil {
		return err
	}

	if err := validations.IsPositive(identifier, rr.Values[4].Value, "RETRY", rr.Type); err != nil {
		return err
	}

	if err := validations.IsPositive(identifier, rr.Values[5].Value, "EXPIRE", rr.Type); err != nil {
		return err
	}

	if err := validations.IsPositive(identifier, rr.Values[6].Value, "NCACHE", rr.Type); err != nil {
		return err
	}

	return nil
}

func insertSerial(rr *models.ResourceRecord, serial string) {
	// We don't have a serial number so use the generated one, we'll need to rearrange the values
	// serial number should be at position two, create a new array with seven elements
	withSerial := make([]*models.ResourceRecordValue, 7)

	//Copy indicies 0 and 1 into the array
	copy(withSerial[:2], rr.Values[:2])

	// insert the serial number
	withSerial[2] = &models.ResourceRecordValue{Value: serial, Comment: generatedSerialNumberComment}

	// Copy the remaining values
	copy(withSerial[3:], rr.Values[2:])

	//We now have a seven element array with the serial number in the 2nd element
	rr.Values = withSerial
}

func init() {
	registerBuiltIn(SOA, &BuiltinPluginSOA{})
}
