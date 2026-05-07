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

func TestNormalize_AAAAPlugin(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.AAAA, Name: "host.example.com", Value: "2001:db8::1"},
		},
		{
			name:       "name-defaulting",
			identifier: "host.example.com",
			rr:         &models.ResourceRecord{Type: models.AAAA, Value: "2001:db8::1"},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "host.example.com", Value: "2001:db8::1"},
			wantErr:    "this plugin does not handle resource records of type 'CNAME' only '[AAAA]', identifier: 'record1'",
		},
		{
			name:       "invalid-name",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.AAAA, Name: "-invalid", Value: "2001:db8::1"},
			wantErr:    "invalid AAAA record, cannot start or end with a hyphen (-): '-invalid', identifier: 'record1'",
		},
		{
			name:       "not-an-ip",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.AAAA, Name: "host.example.com", Value: "not-an-ip"},
			wantErr:    "invalid AAAA record, 'not-an-ip' must be a valid IP address, identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkErr(t, (&BuiltinPluginAAAA{}).Normalize(tc.identifier, tc.rr), tc.wantErr)
		})
	}
}

func TestRender_AAAAPlugin(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.AAAA, Name: "host.example.com", Value: "2001:db8::1"},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.CNAME, Name: "host.example.com"},
			wantErr:    "this plugin does not handle resource records of type 'CNAME' only '[AAAA]', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := (&BuiltinPluginAAAA{}).Render(tc.identifier, tc.rr)
			checkErr(t, err, tc.wantErr)
		})
	}
}
