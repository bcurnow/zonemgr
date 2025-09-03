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
	"github.com/bcurnow/zonemgr/plugins"
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

func TestSOANoramlize2(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &BuiltinPluginSOA{}

	testCases := []struct {
		identifier string
		testConfig *soaNormalization
		err        error
	}{
		{identifier: "valid-dont-generate-serial", testConfig: &soaNormalization{}},
		{identifier: "valid-generate-serial", testConfig: &soaNormalization{generateSerial: true}},
		{identifier: "common-validations-error", testConfig: &soaNormalization{commonValidationsErr: true}, err: testingError},
		{identifier: "identifier-as-name", testConfig: &soaNormalization{identifierAsName: true}},
		{identifier: "is-fully-qualified-name-error", testConfig: &soaNormalization{isFullyQualifiedNameErr: true}, err: testingError},
		{identifier: "value-used-error", testConfig: &soaNormalization{valueUsedErr: true}, err: errors.New("value field cannot be used on SOA records, please use the values field, identifier: 'value-used-error'")},
		{identifier: "comment-used-error", testConfig: &soaNormalization{commentUsedErr: true}, err: errors.New("comment field cannot be used on SOA records, please use the values field, identifier: 'comment-used-error'")},
		{identifier: "generate-serial-error", testConfig: &soaNormalization{generateSerial: true, generateSerialErr: true}, err: testingError},
		{identifier: "normalizer-error", testConfig: &soaNormalization{soaValuesNormalizerErr: true}, err: testingError},
	}

	for _, tc := range testCases {
		rr := testSOA(*tc.testConfig)
		expects := normalizeExpects_SOAPlugin(tc.testConfig)
		expects(tc.identifier, rr, tc.err != nil)
		config := &models.Config{
			GenerateSerial: tc.testConfig.generateSerial,
		}
		// Make sure we call configure because we do use this in the SOA plugin
		plugin.Configure(config)
		// Because we called config, several the objects will now have non-mock values
		// Replace those with mocks
		serialIndexManager = mockSerialIndexManager
		soaValuesNormalizer = mockSoaValuesNormalizer
		if err := plugin.Normalize(tc.identifier, rr); err != nil {
			handleCustomError(t, err, tc.err)
		} else {
			if tc.err != nil {
				t.Errorf("expected error, did not get one")
			} else {
				if tc.testConfig.identifierAsName {
					if rr.Name != tc.identifier {
						t.Errorf("incorrect name: %s, expected %s", rr.Name, tc.identifier)
					}
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
		{zone: &models.Zone{
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
			} else {
				if tc.err.Error() != err.Error() {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		}
	}
}

func TestSOARender(t *testing.T) {
	setup(t)
	defer teardown(t)
	rr := &models.ResourceRecord{
		Type: models.SOA,
		Name: "render",
	}
	plugin := &BuiltinPluginSOA{}
	pluginType := plugins.SOA
	testRender(t, testConfig{
		plugin:     plugin,
		pluginType: pluginType,
		rrType:     rr.Type,
		expects: func(identifier string, rr *models.ResourceRecord, err bool) {
			call := mockValidator.EXPECT().EnsureSupportedPluginType(identifier, rr.Type, pluginType)
			if err {
				call.Return(testingError)
			}
		},
	}, rr)
	//Render uses the standard method so we're going to cheat
	mockValidator.EXPECT().EnsureSupportedPluginType("testing", rr.Type, pluginType)
	_, err := plugin.Render("testing", rr)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

// This feels a bit excessing and like
func normalizeExpects_SOAPlugin(sn *soaNormalization) func(identifier string, rr *models.ResourceRecord, err bool) {
	return func(identifier string, rr *models.ResourceRecord, err bool) {
		call := mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.SOA)
		if sn.commonValidationsErr {
			call.Return(testingError)
			return
		}
		if sn.identifierAsName {
			call = mockValidator.EXPECT().EnsureFullyQualified(identifier, identifier, rr.Type)
		} else {
			call = mockValidator.EXPECT().EnsureFullyQualified(identifier, rr.Name, rr.Type)
		}
		if sn.isFullyQualifiedNameErr {
			call.Return(testingError)
			return
		}
		if sn.valueUsedErr {
			return
		}
		if sn.commentUsedErr {
			return
		}
		serial := ""
		if sn.generateSerial {
			serial = testingSerial
			if sn.identifierAsName {
				call = mockSerialIndexManager.EXPECT().Next(identifier)
			} else {
				call = mockSerialIndexManager.EXPECT().Next(rr.Name)
			}

			if sn.generateSerialErr {
				call.Return("", testingError)
				return
			} else {
				call.Return(serial, nil)
			}
		}
		call = mockSoaValuesNormalizer.EXPECT().Normalize(identifier, rr, mockValidator, sn.generateSerial, serial)
		if sn.soaValuesNormalizerErr {
			call.Return(testingError)
			return
		}
	}
}

func testSOA(sn soaNormalization) *models.ResourceRecord {
	soa := &models.ResourceRecord{Type: models.SOA}

	if !sn.identifierAsName {
		soa.Name = "soa.example.com."
	}

	if sn.hasSerialInValues {
		// We need to have all 7 values populated
		soa.Values = []*models.ResourceRecordValue{
			{Value: "ns1.example.com."},
			{Value: "admin@example.com"},
			{Value: testingSerial},
			{Value: "refresh"},
			{Value: "retry"},
			{Value: "expire"},
			{Value: "ncache"},
		}
	} else {
		// We just need 6 values
		soa.Values = []*models.ResourceRecordValue{
			{Value: "ns1.example.com."},
			{Value: "admin@example.com"},
			{Value: "refresh"},
			{Value: "retry"},
			{Value: "expire"},
			{Value: "ncache"},
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
