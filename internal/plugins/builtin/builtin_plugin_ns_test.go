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

func TestNSNormalize(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.NS, Name: "ns.example.com", Value: "ns1.example.com."},
		},
		{
			name:       "name-defaults-to-wildcard",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.NS, Value: "ns1.example.com."},
		},
		{
			name:       "value-defaults-to-identifier",
			identifier: "ns1.example.com.",
			rr:         &models.ResourceRecord{Type: models.NS, Name: "ns.example.com"},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.A, Name: "ns.example.com", Value: "ns1.example.com."},
			wantErr:    "this plugin does not handle resource records of type 'A' only '[NS]', identifier: 'record1'",
		},
		{
			name:       "invalid-name",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.NS, Name: "-invalid", Value: "ns1.example.com."},
			wantErr:    "invalid NS record, cannot start or end with a hyphen (-): '-invalid', identifier: 'record1'",
		},
		{
			name:       "value-is-ip",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.NS, Name: "ns.example.com", Value: "1.2.3.4"},
			wantErr:    "invalid NS record, '1.2.3.4' must not be an IP address, identifier: 'record1'",
		},
		{
			name:       "value-not-fqdn",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.NS, Name: "ns.example.com", Value: "ns1.example.com"},
			wantErr:    "invalid NS record, must end with a trailing dot: 'ns1.example.com', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkErr(t, (&BuiltinPluginNS{}).Normalize(tc.identifier, tc.rr), tc.wantErr)
		})
	}
}

func TestNSRender(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.NS, Name: "ns.example.com", Value: "ns1.example.com."},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.A, Name: "ns.example.com"},
			wantErr:    "this plugin does not handle resource records of type 'A' only '[NS]', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := (&BuiltinPluginNS{}).Render(tc.identifier, tc.rr)
			checkErr(t, err, tc.wantErr)
		})
	}
}
