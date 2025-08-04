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
)

var _ plugins.TypeHandler = &SOAPlugin{}

var soaSupportedPluginTypes = []plugins.PluginType{plugins.RecordSOA}

type SOAPlugin struct {
	config schema.Config
}

func (p *SOAPlugin) PluginVersion() string {
	return version.Version
}

func (p *SOAPlugin) PluginTypesSupported() []plugins.PluginType {
	return soaSupportedPluginTypes
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

	if err := plugins.IsFullyQualified(rr.Name); err != nil {
		return plugins.NilResourceRecord(), err
	}

	// Validate the the Value and Comment fields are empty, there are no shortcuts for SOA records
	if rr.Value != "" {
		return plugins.NilResourceRecord(), fmt.Errorf("The value field can not be used on SOA records, please use the values field, identifier: '%s'", identifier)
	}

	if rr.Comment != "" {
		return plugins.NilResourceRecord(), fmt.Errorf("The comment field can not be used on SOA records, please use the values field, identifier: '%s'", identifier)
	}

	if err := normalizeValues(&rr, p.config.GenerateSerial, p.config.SerialChangeIndex); err != nil {
		return plugins.NilResourceRecord(), err
	}

	return rr, nil
}

func (p *SOAPlugin) ValidateZone(name string, zone schema.Zone) error {
	//no-op
	return nil
}

func (p *SOAPlugin) Render(identifier string, rr schema.ResourceRecord) (string, error) {
	if err := plugins.IsSupportedPluginType(identifier, &rr, soaSupportedPluginTypes); err != nil {
		return "", err
	}

	return plugins.RenderMultivalueResource(&rr), nil
}

func normalizeValues(rr *schema.ResourceRecord, generateSerial bool, serialChangeIndex uint32) error {
	numValues := len(rr.Values)
	switch numValues {
	case 6:
		// No serial number present, this is an error unless generateSerial is true
		if !generateSerial {
			return fmt.Errorf("Must specify a serial number when generate serial is false, found only 6 values when 7 are required, name: '%s'", rr.Name)
		}

		if err := validateWithNoSerial(rr); err != nil {
			return err
		}
	case 7:
		// There is a serial number provided
		if err := validateWithSerial(rr); err != nil {
			return err
		}
	default:
		return fmt.Errorf("SOA records must have either 6 (no serial number) or 7 values, found %d values, name: '%s'", numValues, rr.Name)
	}

	if generateSerial {
		serial, err := serial(serialChangeIndex)
		if err != nil {
			return err
		}
		rr.Values[2].Value = strconv.Itoa(int(serial))
	} else {
		// Make sure the provided number fitx in an uint32
		if _, err := strconv.ParseUint(rr.Values[2].Value, 10, 32); err != nil {
			return fmt.Errorf("Serial number must be an unsigned 32 bit integer, found '%s', error: %w, name: '%s'", rr.Values[2].Value, err, rr.Name)
		}
	}

	return nil
}

// Validates the values on the
func validateWithNoSerial(rr *schema.ResourceRecord) error {
	// Convert the array to individual variables, they are required to be in a specific order
	// This method is used when there is no serial number field present (6 values total)
	primaryNameServer := rr.Values[0].Value
	administrator := rr.Values[1].Value
	refresh := rr.Values[2].Value
	retry := rr.Values[3].Value
	expire := rr.Values[4].Value
	negativeCache := rr.Values[5].Value

	if err := plugins.IsFullyQualified(primaryNameServer); err != nil {
		return err
	}

	email, err := plugins.FormatEmail(administrator)
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

func validateWithSerial(rr *schema.ResourceRecord) error {
	// Convert the array to individual variables, they are required to be in a specific order
	// This method is used when there is a serial number field present (7 values total)
	primaryNameServer := rr.Values[0].Value
	administrator := rr.Values[1].Value
	// The serial number is Values[2], this is handled elsewhere
	refresh := rr.Values[3].Value
	retry := rr.Values[4].Value
	expire := rr.Values[5].Value
	negativeCache := rr.Values[6].Value

	if err := plugins.IsFullyQualified(primaryNameServer); err != nil {
		return err
	}

	email, err := plugins.FormatEmail(administrator)
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
