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
	"errors"
	"testing"

	"github.com/bcurnow/zonemgr/models"
)

const testingSerial = "testserial"

type soaNormalization struct {
	identifierAsName        bool
	hasSerialInValues       bool
	commonValidationsErr    bool
	isFullyQualifiedNameErr bool
	valueUsedErr            bool
	commentUsedErr          bool
	generateSerial          bool
	generateSerialErr       bool
	soaValuesNormalizerErr  bool
}

func TestSOANormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &BuiltinPluginSOA{}

	testCases := []struct {
		identifier string
		testConfig *soaNormalization
		err        error
	}{
		{identifier: "valid-dont-generate-serial", testConfig: &soaNormalization{hasSerialInValues: true}},
		{identifier: "valid-generate-serial", testConfig: &soaNormalization{generateSerial: true}},
		{identifier: "common-validations-error", testConfig: &soaNormalization{commonValidationsErr: true}, err: errors.New("this plugin does not handle resource records of type 'A' only '[SOA]', identifier: 'common-validations-error'")},
		{identifier: "soa.example.com.", testConfig: &soaNormalization{identifierAsName: true, hasSerialInValues: true}},
		{identifier: "is-fully-qualified-name-error", testConfig: &soaNormalization{isFullyQualifiedNameErr: true}, err: errors.New("invalid SOA record, must end with a trailing dot: 'soa-not-fqdn', identifier: 'is-fully-qualified-name-error'")},
		{identifier: "value-used-error", testConfig: &soaNormalization{valueUsedErr: true}, err: errors.New("invalid SOA record, both value and values are set, identifier: 'value-used-error'")},
		{identifier: "comment-used-error", testConfig: &soaNormalization{commentUsedErr: true}, err: errors.New("invalid SOA record, both comment and values are set, identifier: 'comment-used-error'")},
		{identifier: "generate-serial-error", testConfig: &soaNormalization{generateSerial: true, generateSerialErr: true}, err: errTesting},
		{identifier: "normalizer-error", testConfig: &soaNormalization{soaValuesNormalizerErr: true, hasSerialInValues: true}, err: errors.New("REFRESH must not be less than 0 on a SOA record, was '-1', identifier: 'normalizer-error'")},
	}

	for _, tc := range testCases {
		rr := testSOA(*tc.testConfig)
		setupSerialExpects(t, tc.identifier, tc.testConfig, rr)
		config := &models.Config{GenerateSerial: tc.testConfig.generateSerial}
		if err := plugin.Configure(config); err != nil {
			t.Fatalf("Configure failed: %v", err)
		}
		// Replace serialIndexManager with mock after Configure sets the real one
		serialIndexManager = mockSerialIndexManager

		if err := plugin.Normalize(tc.identifier, rr); err != nil {
			if tc.err == nil {
				t.Errorf("%s - unexpected error: %v", tc.identifier, err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("%s - got error %q, want %q", tc.identifier, err.Error(), tc.err.Error())
			}
		} else {
			if tc.err != nil {
				t.Errorf("%s - expected error %q, got nil", tc.identifier, tc.err)
			} else if tc.testConfig.identifierAsName {
				if rr.Name != tc.identifier {
					t.Errorf("%s - incorrect name: %s, expected %s", tc.identifier, rr.Name, tc.identifier)
				}
			}
		}
	}
}

func TestSOAValidateZone(t *testing.T) {
	testCases := []struct {
		zone *models.Zone
		err  error
	}{
		{zone: &models.Zone{}, err: errors.New("invalid zone, missing SOA record, zone=testing")},
		{zone: &models.Zone{ResourceRecords: map[string]*models.ResourceRecord{"example.com.": {Type: models.SOA}}}},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"example.com.":     {Type: models.SOA},
					"two.example.com.": {Type: models.SOA},
				},
			},
			err: errors.New("more than one SOA record found, only one SOA record is allowed, zone=testing"),
		},
	}

	plugin := &BuiltinPluginSOA{}
	for _, tc := range testCases {
		if err := plugin.ValidateZone("testing", tc.zone); err != nil {
			if tc.err == nil {
				t.Errorf("unexpected error: %s", err)
			} else if tc.err.Error() != err.Error() {
				t.Errorf("incorrect error: %s, want %s", err, tc.err)
			}
		} else if tc.err != nil {
			t.Errorf("expected error %q, got nil", tc.err)
		}
	}
}

func TestSOARender(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.SOA, Name: "render"},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.A, Name: "render"},
			wantErr:    "this plugin does not handle resource records of type 'A' only '[SOA]', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := (&BuiltinPluginSOA{}).Render(tc.identifier, tc.rr)
			checkErr(t, err, tc.wantErr)
		})
	}
}

func setupSerialExpects(t *testing.T, identifier string, sn *soaNormalization, rr *models.ResourceRecord) {
	t.Helper()
	if !sn.generateSerial {
		return
	}
	name := rr.Name
	if sn.identifierAsName {
		name = identifier
	}
	call := mockSerialIndexManager.EXPECT().Next(name)
	if sn.generateSerialErr {
		call.Return("", errTesting)
	} else {
		call.Return(testingSerial, nil)
	}
}

func testSOA(sn soaNormalization) *models.ResourceRecord {
	soa := &models.ResourceRecord{Type: models.SOA}

	if sn.commonValidationsErr {
		soa.Type = models.A
	}

	if !sn.identifierAsName {
		soa.Name = "soa.example.com."
	}

	if sn.isFullyQualifiedNameErr {
		soa.Name = "soa-not-fqdn"
	}

	timerValues := []string{"3600", "900", "604800", "300"}
	if sn.soaValuesNormalizerErr {
		timerValues = []string{"-1", "900", "604800", "300"}
	}

	if sn.hasSerialInValues {
		soa.Values = []*models.ResourceRecordValue{
			{Value: "ns1.example.com."},
			{Value: "admin@example.com"},
			{Value: testingSerial},
			{Value: timerValues[0]},
			{Value: timerValues[1]},
			{Value: timerValues[2]},
			{Value: timerValues[3]},
		}
	} else {
		soa.Values = []*models.ResourceRecordValue{
			{Value: "ns1.example.com."},
			{Value: "admin@example.com"},
			{Value: timerValues[0]},
			{Value: timerValues[1]},
			{Value: timerValues[2]},
			{Value: timerValues[3]},
		}
	}

	if sn.valueUsedErr {
		soa.Value = "i-should-not-be-here"
	}
	if sn.commentUsedErr {
		soa.Comment = "i-should-not-be-here"
	}

	return soa
}
