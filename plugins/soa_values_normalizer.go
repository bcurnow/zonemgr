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
	"github.com/hashicorp/go-hclog"
)

const generatedSerialNumberComment = "Zonemgr generated serial number"

type SOAValuesNormalizer interface {
	// Will take a SOA record and normalize the values
	// generateSerial indicates if we should be generating the serial number
	// serial is the next serial number to use, will be "" if generateSerial is false
	Normalize(identifer string, rr *models.ResourceRecord, validations Validator, generateSerial bool, serial string) error
}

var _ SOAValuesNormalizer = &soaValuesNormalizer{}

func SVN() SOAValuesNormalizer {
	return &soaValuesNormalizer{}
}

type soaValuesNormalizer struct {
}

func (svn *soaValuesNormalizer) Normalize(identifier string, rr *models.ResourceRecord, validations Validator, generateSerial bool, serial string) error {
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
		svn.insertSerial(rr, serial)
	}

	if numValues == 7 && generateSerial {
		// Replace element 2
		rr.Values[2] = &models.ResourceRecordValue{Value: serial, Comment: generatedSerialNumberComment}
	}

	// The only other option is that we have 7 values and we don't generate the serial number so there's nothing to do

	// Now, check the values, first the nameserver
	if err := validations.EnsureFullyQualified(identifier, rr.Values[0].Value, rr.Type); err != nil {
		return err
	}

	email, err := validations.FormatEmail(identifier, rr.Values[1].Value, rr.Type)
	if err != nil {
		return err
	}

	rr.Values[1].Value = email

	//Make sure none of the other values are < 0
	if err := validations.EnsurePositive(identifier, rr.Values[3].Value, "REFRESH", rr.Type); err != nil {
		return err
	}

	if err := validations.EnsurePositive(identifier, rr.Values[4].Value, "RETRY", rr.Type); err != nil {
		return err
	}

	if err := validations.EnsurePositive(identifier, rr.Values[5].Value, "EXPIRE", rr.Type); err != nil {
		return err
	}

	if err := validations.EnsurePositive(identifier, rr.Values[6].Value, "NCACHE", rr.Type); err != nil {
		return err
	}

	return nil
}

func (svn *soaValuesNormalizer) insertSerial(rr *models.ResourceRecord, serial string) {
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
