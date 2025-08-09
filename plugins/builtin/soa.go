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
	"strconv"
	"time"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
	"github.com/bcurnow/zonemgr/version"
	"github.com/hashicorp/go-hclog"
)

const generatedSerialNumberComment = "Zonemgr generated serial number"

var _ plugins.TypeHandler = &SOAPlugin{}

var soaSupportedPluginTypes = []plugins.PluginType{plugins.RecordSOA}

type SOAPlugin struct {
	config schema.Config
}

func (p *SOAPlugin) PluginVersion() (string, error) {
	return version.Version(), nil
}

func (p *SOAPlugin) PluginTypes() ([]plugins.PluginType, error) {
	return soaSupportedPluginTypes, nil
}

func (p *SOAPlugin) Configure(config schema.Config) error {
	p.config = config
	return nil
}

func (p *SOAPlugin) Normalize(identifier string, rr schema.ResourceRecord) (schema.ResourceRecord, error) {
	if err := plugins.StandardValidations(identifier, &rr, soaSupportedPluginTypes); err != nil {
		return plugins.NilResourceRecord(), err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := plugins.IsFullyQualified(rr.Name, identifier, rr); err != nil {
		return plugins.NilResourceRecord(), err
	}

	// Validate the the Value and Comment fields are empty, there are no shortcuts for SOA records
	if rr.Value != "" {
		return plugins.NilResourceRecord(), fmt.Errorf("The value field can not be used on SOA records, please use the values field, identifier: '%s'", identifier)
	}

	if rr.Comment != "" {
		return plugins.NilResourceRecord(), fmt.Errorf("The comment field can not be used on SOA records, please use the values field, identifier: '%s'", identifier)
	}

	if err := normalizeValues(identifier, &rr, p.config.GenerateSerial, p.config.SerialChangeIndex); err != nil {
		return plugins.NilResourceRecord(), err
	}

	return rr, nil
}

func (p *SOAPlugin) ValidateZone(name string, zone schema.Zone) error {
	hasSOA := false
	for identifier, rr := range zone.ResourceRecords {
		// Track the SOA records, there can be only one
		if rr.Type == string(schema.SOA) {
			if hasSOA {
				return fmt.Errorf("More than one SOA record found, only one SOA record is allowed, identifier=%s", identifier)
			}
			hasSOA = true
		}
	}

	if !hasSOA {
		return fmt.Errorf("Invalid zone, missing SOA record, zone=%s", name)
	}

	return nil
}

func (p *SOAPlugin) Render(identifier string, rr schema.ResourceRecord) (string, error) {
	if err := plugins.IsSupportedPluginType(identifier, &rr, soaSupportedPluginTypes); err != nil {
		return "", err
	}

	return plugins.RenderMultivalueResource(&rr), nil
}

func normalizeValues(identifier string, rr *schema.ResourceRecord, generateSerial bool, serialChangeIndex uint32) error {
	numValues := len(rr.Values)
	serialNumber, err := serial(serialChangeIndex)
	if err != nil {
		return err
	}
	generatedSerialNumber := strconv.Itoa(int(serialNumber))

	switch numValues {
	case 6:
		hclog.L().Debug("No serial number present in SOA record, only have 6 values", "identifier", identifier)
		// No serial number present, this is an error unless generateSerial is true
		if !generateSerial {
			return fmt.Errorf("Must specify a serial number when generate serial is false, found only 6 values when 7 are required, name: '%s'", rr.Name)
		}

		if err := validateWithNoSerial(identifier, rr, generatedSerialNumber); err != nil {
			return err
		}
	case 7:
		hclog.L().Debug("Serial number present in SOA record", "identifier", identifier, "generateSerialNumber", generateSerial)
		// There is a serial number provided
		if err := validateWithSerial(identifier, rr, generateSerial, generatedSerialNumber); err != nil {
			return err
		}
	default:
		return fmt.Errorf("SOA records must have either 6 (no serial number) or 7 values, found %d values, name: '%s'", numValues, rr.Name)
	}

	return nil
}

// Validates the values on the
func validateWithNoSerial(identifier string, rr *schema.ResourceRecord, generatedSerialNumber string) error {
	// Convert the array to individual variables, they are required to be in a specific order
	// This method is used when there is no serial number field present (6 values total)
	primaryNameServer := rr.Values[0].Value
	administrator := rr.Values[1].Value
	refresh := rr.Values[2].Value
	retry := rr.Values[3].Value
	expire := rr.Values[4].Value
	negativeCache := rr.Values[5].Value

	if err := plugins.IsFullyQualified(primaryNameServer, identifier, *rr); err != nil {
		return err
	}

	email, err := plugins.FormatEmail(administrator, identifier, *rr)
	if err != nil {
		return err
	}

	//Make sure none of the other values are < 0
	if err := greaterThanZero(refresh, "REFRESH", rr); err != nil {
		return err
	}

	if err := greaterThanZero(retry, "RETRY", rr); err != nil {
		return err
	}

	if err := greaterThanZero(expire, "EXPIRE", rr); err != nil {
		return err
	}

	if err := greaterThanZero(negativeCache, "NCACHE", rr); err != nil {
		return err
	}

	//Now we need to reset the values to use the generated serial number because it still needs to be at index 2
	newValues := make([]schema.ResourceRecordValue, 7)
	newValues[0] = schema.ResourceRecordValue{Value: primaryNameServer, Comment: rr.Values[0].Comment}
	newValues[1] = schema.ResourceRecordValue{Value: email, Comment: rr.Values[1].Comment}
	newValues[2] = schema.ResourceRecordValue{Value: generatedSerialNumber, Comment: generatedSerialNumberComment}
	// Note that the comment values start at 2 and not 3, that's because there's only 6 records so we need to shift up
	newValues[3] = schema.ResourceRecordValue{Value: refresh, Comment: rr.Values[2].Comment}
	newValues[4] = schema.ResourceRecordValue{Value: retry, Comment: rr.Values[3].Comment}
	newValues[5] = schema.ResourceRecordValue{Value: expire, Comment: rr.Values[4].Comment}
	newValues[6] = schema.ResourceRecordValue{Value: negativeCache, Comment: rr.Values[5].Comment}
	rr.Values = newValues

	return nil
}

func validateWithSerial(identifier string, rr *schema.ResourceRecord, generateSerial bool, generatedSerialNumber string) error {
	// Convert the array to individual variables, they are required to be in a specific order
	// This method is used when there is a serial number field present (7 values total)
	primaryNameServer := rr.Values[0].Value
	administrator := rr.Values[1].Value
	if rr.Values[2].Value != "" && generateSerial {
		hclog.L().Debug("Ignoring serial number of SOA record, using generated one", "identifier", identifier, "serialNumber", rr.Values[2].Value, "generateSerial", generateSerial, "generatedSerialNumber", generatedSerialNumber)
	}
	if generateSerial {
		rr.Values[2].Value = generatedSerialNumber
		hclog.L().Debug("Replacing existing comment due to generated serial number", "oldComment", rr.Values[2].Comment, "newComment", generatedSerialNumberComment)
		rr.Values[2].Comment = generatedSerialNumberComment
	}
	refresh := rr.Values[3].Value
	retry := rr.Values[4].Value
	expire := rr.Values[5].Value
	negativeCache := rr.Values[6].Value

	if err := plugins.IsFullyQualified(primaryNameServer, identifier, *rr); err != nil {
		return err
	}

	email, err := plugins.FormatEmail(administrator, identifier, *rr)
	if err != nil {
		return err
	}

	rr.Values[1].Value = email

	//Make sure none of the other values are < 0
	if err := greaterThanZero(refresh, "REFRESH", rr); err != nil {
		return err
	}

	if err := greaterThanZero(retry, "RETRY", rr); err != nil {
		return err
	}

	if err := greaterThanZero(expire, "EXPIRE", rr); err != nil {
		return err
	}

	if err := greaterThanZero(negativeCache, "NCACHE", rr); err != nil {
		return err
	}

	return nil
}

func greaterThanZero(str string, fieldName string, rr *schema.ResourceRecord) error {
	value, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return err
	}

	if value < 0 {
		return fmt.Errorf("%s must not be less than 0 on a SOA record, was '%s', name: '%s'", fieldName, str, rr.Name)
	}

	return nil
}

// Generates a time-based serial number plus a numeric index
func serial(index uint32) (uint32, error) {
	t := time.Now()
	serialString := fmt.Sprintf("%02d%02d%04d%02d", t.Day(), t.Month(), t.Year(), index)

	parsedSerial, err := strconv.ParseUint(serialString, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Unable to generate a serial number from day: %d, month: %d, year: %d, changeIndex: %d: %w", t.Day(), t.Month(), t.Year(), index, err)
	}

	// Explicitly convert to a uint32
	return uint32(parsedSerial), nil
}

func init() {
	registerBuiltIn(plugins.RecordSOA, &SOAPlugin{})
}
