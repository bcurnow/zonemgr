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
	testCases := []struct {
		zone *models.Zone
		err  error
	}{
		{
			zone: &models.Zone{ResourceRecords: map[string]*models.ResourceRecord{}},
		},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"lonecname": {
						Name: "lonecname",
						Type: models.CNAME,
					},
				},
			},
			err: errors.New("found CNAME records but there are no A records present, all CNAMES must reference an A record name, zone: 'testing'"),
		},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"cname": {
						Name:  "badcname",
						Type:  models.CNAME,
						Value: "notgoingtofindme",
					},
					"arecord": {
						Name:  "nottherightone",
						Type:  models.A,
						Value: "1.2.3.4",
					},
				},
			},
			err: errors.New("invalid CNAME record, 'cname' has a value of 'notgoingtofindme' which does not match any defined A record name, zone: 'testing'"),
		},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"arecord": {
						Name:  "arecord",
						Type:  models.A,
						Value: "1.2.3.4",
					},
					"cname": {
						Name:  "cname",
						Type:  models.CNAME,
						Value: "arecord",
					},
				},
			},
		},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"differentidentifier": {
						Name:  "arecord",
						Type:  models.A,
						Value: "1.2.3.4",
					},
					"cname": {
						Name:  "cname",
						Type:  models.CNAME,
						Value: "arecord",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		if err := plugin.ValidateZone("testing", tc.zone); err != nil {
			if tc.err == nil {
				t.Errorf("unexpected error: %s", err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("incorrect error: %s, want %s", err, tc.err)
			}
		} else if tc.err != nil {
			t.Errorf("expected error %q, got nil", tc.err)
		}
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
