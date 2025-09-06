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

package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestString_Zone(t *testing.T) {
	zone := &Zone{
		Config: &Config{},
		ResourceRecords: map[string]*ResourceRecord{
			"example.com.":     {Type: SOA},
			"ns1.example.com.": {Type: NS},
		},
		TTL: &TTL{Value: toInt32Ptr(33), Comment: "ttl comment"},
	}
	want := "Zone{\n" +
		"   Config: Config{ GenerateSerial: false, GenerateReverseLookupZones: false, SerialChangeIndexDirectory:  }\n" +
		"   ResourceRecords:\n" +
		"     example.com. -> ResourceRecord{\n" +
		"       Name: \n" +
		"       Type: SOA\n" +
		"       Class: \n" +
		"       TTL: <nil>\n" +
		"       Values: []\n" +
		"       Value: \n" +
		"       Comment: \n" +
		"     }\n" +
		"     ns1.example.com. -> ResourceRecord{\n" +
		"       Name: \n" +
		"       Type: NS\n" +
		"       Class: \n" +
		"       TTL: <nil>\n" +
		"       Values: []\n" +
		"       Value: \n" +
		"       Comment: \n" +
		"     }\n" +
		"   TTL: TTL{ Value: 33, Comment: ttl comment }\n" +
		"}"

	if cmp.Diff(zone.String(), want) != "" {
		t.Errorf("incorrect string:\n%s", cmp.Diff(zone.String(), want))
	}

	zone = &Zone{}
	want = "Zone{\n" +
		"   Config: <nil>\n" +
		"   ResourceRecords:\n" +
		"   TTL: <nil>\n" +
		"}"

	if cmp.Diff(zone.String(), want) != "" {
		t.Errorf("incorrect string:\n%s", cmp.Diff(zone.String(), want))
	}
}

func TestSOARecord(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want *ResourceRecord
	}{
		{rr: &ResourceRecord{Type: SOA}, want: &ResourceRecord{Type: SOA}},
		{rr: &ResourceRecord{Type: NS}, want: nil},
		{rr: nil, want: nil},
	}

	for _, tc := range testCases {
		z := &Zone{ResourceRecords: map[string]*ResourceRecord{"one": tc.rr}}
		if tc.rr == nil {
			delete(z.ResourceRecords, "one")
		}
		actual := z.SOARecord()
		if !cmp.Equal(actual, tc.want) {
			t.Errorf("incorrect record: '%s', want: '%s'", actual, tc.want)
		}
	}
}

func TestResourceRecordsByType(t *testing.T) {
	rrSOA := &ResourceRecord{Type: SOA}
	rrNS := &ResourceRecord{Type: NS}
	testCases := []struct {
		zone *Zone
		want map[ResourceRecordType]map[string]*ResourceRecord
	}{
		{
			zone: &Zone{
				ResourceRecords: map[string]*ResourceRecord{
					"one": rrSOA,
					"two": rrNS,
				},
			},
			want: map[ResourceRecordType]map[string]*ResourceRecord{
				SOA: {"one": rrSOA},
				NS:  {"two": rrNS},
			},
		},
		{
			zone: &Zone{ResourceRecords: make(map[string]*ResourceRecord)},
			want: make(map[ResourceRecordType]map[string]*ResourceRecord),
		},
	}

	for _, tc := range testCases {
		actual := tc.zone.ResourceRecordsByType()

		if cmp.Diff(actual, tc.want) != "" {
			t.Errorf("incorrect results:\n%s", cmp.Diff(actual, tc.want))
		}
	}
}

func TestSortedResourceRecordKeys(t *testing.T) {
	testCases := []struct {
		zone *Zone
		want []string
	}{
		{
			zone: &Zone{
				ResourceRecords: map[string]*ResourceRecord{
					"one":   {},
					"two":   {},
					"three": {},
					"four":  {},
				},
			},
			want: []string{"four", "one", "three", "two"},
		},
		{
			zone: &Zone{ResourceRecords: make(map[string]*ResourceRecord)},
			want: []string{},
		},
	}

	for _, tc := range testCases {
		actual := tc.zone.sortedResourceRecordKeys()

		if cmp.Diff(actual, tc.want) != "" {
			t.Errorf("incorrect results:\n%s", cmp.Diff(actual, tc.want))
		}
	}
}
