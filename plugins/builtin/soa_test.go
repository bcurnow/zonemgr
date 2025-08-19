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

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
)

func TestSOAPluginVersion(t *testing.T) {
	testPluginVersion(t, &SOAPlugin{})
}

func TestSOAPluginTypes(t *testing.T) {
	testPluginTypes(t, &SOAPlugin{}, plugins.SOA)
}

func TestSOAConfigure(t *testing.T) {
	testConfigure(t, &SOAPlugin{})
}

func TestSOANormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &SOAPlugin{}
	testCases := []struct {
		rr         *schema.ResourceRecord
		identifier string
		name       string
		err        string
		config     *schema.Config
	}{
		{
			identifier: "ValidRecordWithOutName",
			rr: &schema.ResourceRecord{
				Type: schema.SOA,
				Values: []*schema.ResourceRecordValue{
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
			config: &schema.Config{GenerateSerial: false},
		},
		{
			identifier: "Valid record with a name",
			rr: &schema.ResourceRecord{
				Type: schema.SOA,
				Values: []*schema.ResourceRecordValue{
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
			config: &schema.Config{GenerateSerial: false},
		},
	}

	for _, tc := range testCases {
		mockValidations.EXPECT().StandardValidations(tc.identifier, tc.rr, plugins.SOA)
		mockValidations.EXPECT().IsFullyQualified(tc.identifier, tc.name, tc.rr.Type)
		mockValidations.EXPECT().IsFullyQualified(tc.identifier, tc.rr.Values[0].Value, tc.rr.Type)
		mockValidations.EXPECT().FormatEmail(tc.identifier, tc.rr.Values[1].Value, tc.rr.Type)
		mockValidations.EXPECT().IsPositive(tc.identifier, tc.rr.Values[3].Value, "REFRESH", tc.rr.Type)
		mockValidations.EXPECT().IsPositive(tc.identifier, tc.rr.Values[4].Value, "RETRY", tc.rr.Type)
		mockValidations.EXPECT().IsPositive(tc.identifier, tc.rr.Values[5].Value, "EXPIRE", tc.rr.Type)
		mockValidations.EXPECT().IsPositive(tc.identifier, tc.rr.Values[6].Value, "NCACHE", tc.rr.Type)
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
	plugin := &SOAPlugin{}
	zone := &schema.Zone{
		ResourceRecords: map[string]*schema.ResourceRecord{
			"example.com.": {
				Type: schema.SOA,
				Values: []*schema.ResourceRecordValue{
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
	mockValidations.EXPECT().IsSupportedPluginType("testing", schema.SOA, plugins.SOA)
	plugin := &SOAPlugin{}
	_, err := plugin.Render("testing", &schema.ResourceRecord{Type: schema.SOA})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
