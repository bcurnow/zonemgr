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
package dns

import (
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLastOctet(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	testCases := []struct {
		value string
		want  string
	}{
		{value: "4.3.2.1", want: "1"},
		{value: "255.255.255.0", want: "0"},
		{value: "1.2.3", want: "3"},
		{value: "bogus", want: "bogus"},
	}

	for _, tc := range testCases {
		res := (&zoneReverser{}).lastOctet(tc.value)

		if res != tc.want {
			t.Errorf("Unexpected result: '%s', want '%s'", res, tc.want)
		}
	}
}

func TestReverseZoneName(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	testCases := []struct {
		value string
		want  string
	}{
		{value: "4.3.2.1", want: "2.3.4.in-addr.arpa."},
		{value: "255.255.255.0", want: "255.255.255.in-addr.arpa."},
		{value: "1.2.3", want: "2.1.in-addr.arpa."},
		{value: "bogus", want: ".in-addr.arpa."},
	}

	for _, tc := range testCases {
		res := (&zoneReverser{}).reverseZoneName(tc.value)

		if res != tc.want {
			t.Errorf("Unexpected result: '%s', want '%s'", res, tc.want)
		}
	}
}

func TestToPTR(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	testCases := []struct {
		name string
		rr   *models.ResourceRecord
		want *models.ResourceRecord
	}{
		{
			name: "short name",
			rr:   &models.ResourceRecord{Type: models.A, Name: "testing", Value: "1.2.3.4"},
			want: &models.ResourceRecord{Type: models.PTR, Name: "4", Value: "testing.example.com.", Values: []*models.ResourceRecordValue{}},
		},
		{
			name: "fully qualified",
			rr:   &models.ResourceRecord{Type: models.A, Name: "testing.somedomain.com.", Value: "1.2.3.4"},
			want: &models.ResourceRecord{Type: models.PTR, Name: "4", Value: "testing.somedomain.com.", Values: []*models.ResourceRecordValue{}},
		},
		{
			name: "fully qualified with testing domain",
			rr:   &models.ResourceRecord{Type: models.A, Name: "testing.example.com.", Value: "1.2.3.4"},
			want: &models.ResourceRecord{Type: models.PTR, Name: "4", Value: "testing.example.com.", Values: []*models.ResourceRecordValue{}},
		},
		{
			name: "in testing domain but missing trailing dot",
			rr:   &models.ResourceRecord{Type: models.A, Name: "testing.example.com", Value: "1.2.3.4"},
			want: &models.ResourceRecord{Type: models.PTR, Name: "4", Value: "testing.example.com.example.com.", Values: []*models.ResourceRecordValue{}},
		},
	}

	for _, tc := range testCases {
		ptr := (&zoneReverser{}).toPTR("example.com", tc.rr)

		if !cmp.Equal(ptr, tc.want) {
			t.Errorf("Unexpected result for %s:\n%s", tc.name, cmp.Diff(ptr, tc.want))
		}
	}
}

func TestReverseZone(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	zone := &models.Zone{
		Config: &models.Config{},
		TTL:    &models.TTL{},
		ResourceRecords: map[string]*models.ResourceRecord{
			"record1": {Type: models.A, Name: "one", Value: "1.2.3.4"},
			"record2": {Type: models.A, Name: "two", Value: "10.2.2.5"},
			"record3": {Type: models.NS, Name: "doesn't matter", Value: "also doesn't matter"},
			"record4": {Type: models.SOA, Name: "SOA", Value: "SOA"},
		},
	}

	wantedReverseZones := map[string]*models.Zone{
		"3.2.1.in-addr.arpa.": {
			Config: zone.Config,
			TTL:    zone.TTL,
			ResourceRecords: map[string]*models.ResourceRecord{
				"3.2.1.in-addr.arpa.": {Type: models.SOA, Name: "3.2.1.in-addr.arpa.", Value: "SOA"},
				"4":                   {Type: models.PTR, Name: "4", Value: "one.testing.example.com.", Values: []*models.ResourceRecordValue{}},
			},
		},
		"2.2.10.in-addr.arpa.": {
			Config: zone.Config,
			TTL:    zone.TTL,
			ResourceRecords: map[string]*models.ResourceRecord{
				"2.2.10.in-addr.arpa.": {Type: models.SOA, Name: "2.2.10.in-addr.arpa.", Value: "SOA"},
				"5":                    {Type: models.PTR, Name: "5", Value: "two.testing.example.com.", Values: []*models.ResourceRecordValue{}},
			},
		},
	}

	reverseZones := (&zoneReverser{}).ReverseZone("testing.example.com.", zone)
	if len(reverseZones) != 2 {
		t.Errorf("Expected 2 reverse zones, got %d", len(reverseZones))
	}

	for zoneName, wantedReverseZone := range wantedReverseZones {
		reverseZone, ok := reverseZones[zoneName]
		if !ok {
			t.Errorf("Expected to find zone '%s' but was missing", zoneName)
		}

		if !cmp.Equal(reverseZone, wantedReverseZone, cmpopts.IgnoreUnexported(models.Zone{})) {
			t.Errorf("Incorrect reverse zone:\n%s", cmp.Diff(reverseZone, wantedReverseZone, cmpopts.IgnoreUnexported(models.Zone{})))
		}
	}
}

func TestReverseZone_NoResourceRecords(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	zone := &models.Zone{
		Config:          &models.Config{},
		TTL:             &models.TTL{},
		ResourceRecords: map[string]*models.ResourceRecord{},
	}

	reverseZones := (&zoneReverser{}).ReverseZone("testing.example.com.", zone)
	if len(reverseZones) != 0 {
		t.Errorf("Expected no zones but found %d", len(reverseZones))
	}
}

func TestReverser(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	one := Reverser()
	two := Reverser()

	if one == two {
		t.Error("Expected two different instances")
	}
}
