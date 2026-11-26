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
)

func TestNormalize_CNAMEPlugin(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "alias.example.com", Value: "target.example.com"},
		},
		{
			name:       "name-defaulting",
			identifier: "alias.example.com",
			rr:         &models.ResourceRecord{Type: models.CNAME, Value: "target.example.com"},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.A, Name: "alias.example.com", Value: "target.example.com"},
			wantErr:    "this plugin does not handle resource records of type 'A' only '[CNAME]', identifier: 'record1'",
		},
		{
			name:       "invalid-name",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "-invalid", Value: "target.example.com"},
			wantErr:    "invalid CNAME record, cannot start or end with a hyphen (-): '-invalid', identifier: 'record1'",
		},
		{
			name:       "value-is-ip",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "alias.example.com", Value: "1.2.3.4"},
			wantErr:    "invalid CNAME record, '1.2.3.4' must not be an IP address, identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkErr(t, (&BuiltinPluginCNAME{}).Normalize(tc.identifier, tc.rr), tc.wantErr)
		})
	}
}

func TestValidateZone_CNAMEPlugin(t *testing.T) {
	plugin := &BuiltinPluginCNAME{}
	record := func(name string, recordType models.ResourceRecordType, value string) *models.ResourceRecord {
		return &models.ResourceRecord{Name: name, Type: recordType, Value: value}
	}
	testCases := []struct {
		name    string
		records map[string]*models.ResourceRecord
		wantErr string
	}{
		{
			name:    "empty-zone",
			records: map[string]*models.ResourceRecord{},
		},
		{
			name: "target-is-a-record",
			records: map[string]*models.ResourceRecord{
				"arecord": record("arecord", models.A, "1.2.3.4"),
				"cname":   record("cname", models.CNAME, "arecord"),
			},
		},
		{
			name: "target-is-aaaa-record-only",
			records: map[string]*models.ResourceRecord{
				"arecord": record("arecord", models.AAAA, "::1"),
				"cname":   record("cname", models.CNAME, "arecord"),
			},
		},
		{
			name: "target-does-not-exist-in-zone",
			records: map[string]*models.ResourceRecord{
				"cname": record("cname", models.CNAME, "notgoingtofindme"),
			},
		},
		{
			name: "target-in-another-zone",
			records: map[string]*models.ResourceRecord{
				"cname": record("cname", models.CNAME, "host.other.org."),
			},
		},
		{
			name: "multiple-names-alias-one-target",
			records: map[string]*models.ResourceRecord{
				"grafana":    record("grafana", models.AAAA, "::1"),
				"dashboards": record("dashboards", models.CNAME, "grafana"),
				"monitoring": record("monitoring", models.CNAME, "grafana"),
			},
		},
		{
			name: "cname-chain",
			records: map[string]*models.ResourceRecord{
				"first":  record("first", models.CNAME, "second"),
				"second": record("second", models.CNAME, "third"),
				"third":  record("third", models.A, "1.2.3.4"),
			},
		},
		{
			name: "target-from-values",
			records: map[string]*models.ResourceRecord{
				"cname": {Name: "cname", Type: models.CNAME, Values: []*models.ResourceRecordValue{{Value: "elsewhere"}}},
			},
		},
		{
			name: "relative-and-absolute-names-differ-by-owner",
			records: map[string]*models.ResourceRecord{
				"cname": record("www", models.CNAME, "host"),
				"other": record("www.other.org.", models.A, "1.2.3.4"),
			},
		},
		{
			name: "cname-with-a-record",
			records: map[string]*models.ResourceRecord{
				"cname": record("www", models.CNAME, "host"),
				"arec":  record("www", models.A, "1.2.3.4"),
			},
			wantErr: "invalid CNAME record, 'cname' has the name 'www' which also has a record of type 'A', a name with a CNAME cannot have any other record types, zone: 'example.com'",
		},
		{
			name: "cname-with-other-record-relative-vs-absolute",
			records: map[string]*models.ResourceRecord{
				"cname": record("www", models.CNAME, "host"),
				"txt":   record("WWW.example.com.", models.TXT, "hello"),
			},
			wantErr: "invalid CNAME record, 'cname' has the name 'www' which also has a record of type 'TXT', a name with a CNAME cannot have any other record types, zone: 'example.com'",
		},
		{
			name: "cname-at-apex",
			records: map[string]*models.ResourceRecord{
				"soa":   record("example.com.", models.SOA, ""),
				"cname": record("@", models.CNAME, "host"),
			},
			wantErr: "invalid CNAME record, 'cname' has the name '@' which also has a record of type 'SOA', a name with a CNAME cannot have any other record types, zone: 'example.com'",
		},
		{
			name: "duplicate-cname",
			records: map[string]*models.ResourceRecord{
				"one": record("www", models.CNAME, "a"),
				"two": record("www", models.CNAME, "b"),
			},
			wantErr: "invalid CNAME record, 'one' and 'two' both define a CNAME for the name 'www', a name can only have one CNAME, zone: 'example.com'",
		},
		{
			name: "self-loop",
			records: map[string]*models.ResourceRecord{
				"loop": record("loop", models.CNAME, "loop"),
			},
			wantErr: "invalid CNAME record, 'loop' is part of a CNAME loop, zone: 'example.com'",
		},
		{
			name: "multi-step-loop",
			records: map[string]*models.ResourceRecord{
				"a": record("a", models.CNAME, "b"),
				"b": record("b", models.CNAME, "c"),
				"c": record("c", models.CNAME, "a.example.com."),
			},
			wantErr: "invalid CNAME record, 'a' is part of a CNAME loop, zone: 'example.com'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkErr(t, plugin.ValidateZone("example.com", &models.Zone{ResourceRecords: tc.records}), tc.wantErr)
		})
	}
}

func TestRender_CNAMEPlugin(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "alias.example.com", Value: "target.example.com"},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.A, Name: "alias.example.com"},
			wantErr:    "this plugin does not handle resource records of type 'A' only '[CNAME]', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := (&BuiltinPluginCNAME{}).Render(tc.identifier, tc.rr)
			checkErr(t, err, tc.wantErr)
		})
	}
}
