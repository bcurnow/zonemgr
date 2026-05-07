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

func TestPTRNormalize(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.PTR, Name: "ptr-record", Value: "host.example.com."},
		},
		{
			name:       "name-defaulting",
			identifier: "ptr-record",
			rr:         &models.ResourceRecord{Type: models.PTR, Value: "host.example.com."},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "ptr-record", Value: "host.example.com."},
			wantErr:    "this plugin does not handle resource records of type 'CNAME' only '[PTR]', identifier: 'record1'",
		},
		{
			name:       "invalid-name",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.PTR, Name: "-invalid", Value: "host.example.com."},
			wantErr:    "invalid PTR record, cannot start or end with a hyphen (-): '-invalid', identifier: 'record1'",
		},
		{
			name:       "value-not-fqdn",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.PTR, Name: "ptr-record", Value: "host.example.com"},
			wantErr:    "invalid PTR record, must end with a trailing dot: 'host.example.com', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkErr(t, (&BuiltinPluginPTR{}).Normalize(tc.identifier, tc.rr), tc.wantErr)
		})
	}
}

func TestPTRRender(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.PTR, Name: "ptr-record", Value: "host.example.com."},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "ptr-record"},
			wantErr:    "this plugin does not handle resource records of type 'CNAME' only '[PTR]', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := (&BuiltinPluginPTR{}).Render(tc.identifier, tc.rr)
			checkErr(t, err, tc.wantErr)
		})
	}
}
