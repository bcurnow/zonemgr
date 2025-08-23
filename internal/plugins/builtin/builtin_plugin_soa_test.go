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
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
)

func TestSOAPluginVersion(t *testing.T) {
	testPluginVersion(t, &BuiltinPluginSOA{})
}

func TestSOAPluginTypes(t *testing.T) {
	testPluginTypes(t, &BuiltinPluginSOA{}, plugins.SOA)
}

func TestSOAConfigure(t *testing.T) {
	testConfigure(t, &BuiltinPluginSOA{})
}

func TestSOANormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &BuiltinPluginSOA{}
	testCases := []struct {
		rr         *models.ResourceRecord
		identifier string
		name       string
		err        string
		config     *models.Config
	}{
		{
			identifier: "ValidRecordWithOutName",
			rr: &models.ResourceRecord{
				Type: models.SOA,
				Values: []*models.ResourceRecordValue{
					{Value: "ns.example.com."},
					{Value: "admin@example.com"},
					{Value: "99"},
					{Value: "1"},
					{Value: "2"},
					{Value: "3"},
					{Value: "4"},
				},
			},
			name:   "ValidRecordWithOutName",
			config: &models.Config{GenerateSerial: false},
		},
		{
			identifier: "Valid record with a name",
			rr: &models.ResourceRecord{
				Type: models.SOA,
				Values: []*models.ResourceRecordValue{
					{Value: "ns.example.com."},
					{Value: "admin@example.com"},
					{Value: "99"},
					{Value: "1"},
					{Value: "2"},
					{Value: "3"},
					{Value: "4"},
				},
				Name: "name",
			},
			name:   "name",
			config: &models.Config{GenerateSerial: false},
		},
	}

	for _, tc := range testCases {
		mockValidator.EXPECT().StandardValidations(tc.identifier, tc.rr, plugins.SOA)
		mockValidator.EXPECT().IsFullyQualified(tc.identifier, tc.name, tc.rr.Type)
		mockValidator.EXPECT().IsFullyQualified(tc.identifier, tc.rr.Values[0].Value, tc.rr.Type)
		mockValidator.EXPECT().FormatEmail(tc.identifier, tc.rr.Values[1].Value, tc.rr.Type)
		mockValidator.EXPECT().IsPositive(tc.identifier, tc.rr.Values[3].Value, "REFRESH", tc.rr.Type)
		mockValidator.EXPECT().IsPositive(tc.identifier, tc.rr.Values[4].Value, "RETRY", tc.rr.Type)
		mockValidator.EXPECT().IsPositive(tc.identifier, tc.rr.Values[5].Value, "EXPIRE", tc.rr.Type)
		mockValidator.EXPECT().IsPositive(tc.identifier, tc.rr.Values[6].Value, "NCACHE", tc.rr.Type)
		plugin.Configure(tc.config)

		err := plugin.Normalize(tc.identifier, tc.rr)
		if err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("unexpected error: '%s', want '%s'", err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("expected error")
			}
		}
	}
}

func TestSOAValidateZone(t *testing.T) {
	plugin := &BuiltinPluginSOA{}
	zone := &models.Zone{
		ResourceRecords: map[string]*models.ResourceRecord{
			"example.com.": {
				Type: models.SOA,
				Values: []*models.ResourceRecordValue{
					{Value: "ns.example.com."},
					{Value: "admin@example.com"},
					{Value: "99"},
					{Value: "1"},
					{Value: "2"},
					{Value: "3"},
					{Value: "4"},
				},
			},
		},
	}

	if err := plugin.ValidateZone("example.com.", zone); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestSOARender(t *testing.T) {
	setup(t)
	defer teardown(t)
	//Render uses the standard method so we're going to cheat
	mockValidator.EXPECT().IsSupportedPluginType("testing", models.SOA, plugins.SOA)
	plugin := &BuiltinPluginSOA{}
	_, err := plugin.Render("testing", &models.ResourceRecord{Type: models.SOA})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
