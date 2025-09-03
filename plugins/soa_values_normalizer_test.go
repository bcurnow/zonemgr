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
	"testing"

	"github.com/bcurnow/zonemgr/models"
)

var (
	soaSevenValues models.ResourceRecord
	soaSixValues   models.ResourceRecord
)

func resetRecords() {
	soaSevenValues = *defaultSOA("20250903", "3", "4", "5", "6")
	soaSixValues = *defaultSOA("3", "4", "5", "6")
}
func TestNormalize(t *testing.T) {
	testCases := []struct {
		name           string
		insertSerial   bool
		rr             *models.ResourceRecord
		generateSerial bool
		err            string
	}{
		{name: "success-withserial", rr: &soaSevenValues},
		{name: "success-noserial", rr: &soaSixValues, generateSerial: true},
		{name: "fail-noserial", rr: &soaSixValues, err: "must specify a serial number when generate serial is false, found only 6 values when 7 are required, identifier: 'testing', name: 'example.com.'"},
		{name: "wrong-number-of-values-too-many", rr: defaultSOA("20050903", "3", "4", "5", "6", "7"), err: "incorrect number of values for the SOA record, expected 6 (no serial) or 7, found 8, identifier: 'testing', name: 'example.com.'"},
		{name: "wrong-number-of-values-too-few", rr: defaultSOA(), err: "incorrect number of values for the SOA record, expected 6 (no serial) or 7, found 2, identifier: 'testing', name: 'example.com.'"},
		{name: "seven-values-generate-serial-true", rr: &soaSevenValues, generateSerial: true},
		{name: "not-fully-qualified-nameserver", rr: customSOA("ns", "admin", "serial", "1", "1", "1", "1"), err: "invalid SOA record, must end with a trailing dot: 'ns', identifier: 'testing'"},
		{name: "invalid-email", rr: customSOA("ns.example.com.", "bogus-not-email", "serial", "1", "1", "1", "1"), err: "invalid SOA record, must end with a trailing dot: 'bogus-not-email', identifier: 'testing'"},
		{name: "refresh-non-positive", rr: defaultSOA("serial", "-1", "1", "1", "1"), err: "REFRESH must not be less than 0 on a SOA record, was '-1', identifier: 'testing'"},
		{name: "retry-non-positive", rr: defaultSOA("serial", "1", "-1", "1", "1"), err: "RETRY must not be less than 0 on a SOA record, was '-1', identifier: 'testing'"},
		{name: "expire-non-positive", rr: defaultSOA("serial", "1", "1", "-1", "1"), err: "EXPIRE must not be less than 0 on a SOA record, was '-1', identifier: 'testing'"},
		{name: "ncache-non-positive", rr: defaultSOA("serial", "1", "1", "1", "-1"), err: "NCACHE must not be less than 0 on a SOA record, was '-1', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		resetRecords()
		if err := (&SOAValuesNormalizer{}).Normalize("testing", tc.rr, V(), tc.generateSerial, "testing-serial"); err != nil {
			if tc.err == "" {
				t.Errorf("%s - unexpected error: %s", tc.name, err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("%s - incorrect error: '%s', want: '%s'", tc.name, err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("%s - expected an error, found none", tc.name)
			}

			// Validate that the serial number is valid
			if len(tc.rr.Values) != 7 {
				t.Errorf("%s, Incorrect number of values %d, expected 7", tc.name, len(tc.rr.Values))
			}

			if tc.generateSerial {
				if tc.rr.Values[2].Value != "testing-serial" {
					t.Errorf("%s, did not use testing-serial, found: '%s'", tc.name, tc.rr.Values[2].Value)
				}
			}
		}
	}
}

func customSOA(ns string, admin string, values ...string) *models.ResourceRecord {
	soaValues := make([]*models.ResourceRecordValue, len(values)+2)
	soaValues[0] = &models.ResourceRecordValue{Value: ns}
	soaValues[1] = &models.ResourceRecordValue{Value: admin}

	for i, value := range values {
		soaValues[i+2] = &models.ResourceRecordValue{
			Value: value,
		}
	}

	soa := &models.ResourceRecord{
		Name:   "example.com.",
		Type:   models.SOA,
		Values: soaValues,
	}

	return soa
}

func defaultSOA(values ...string) *models.ResourceRecord {
	return customSOA("ns.example.com.", "admin@example.com", values...)
}
