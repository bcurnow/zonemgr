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

package utils

import (
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		input string
		err   string
	}{
		{input: "1.2.3.4", err: ""},
		{input: "fdda:5cc1:23:4::1f", err: ""},
		{input: "fdda:5cc1:0023:0004:0000:0000:0000:001f", err: ""},
		{input: "invalid", err: "ParseAddr(\"invalid\"): unable to parse IP"},
	}

	for _, tc := range testCases {
		_, err := ParseIP(tc.input)
		if err != nil && tc.err == "" {
			t.Errorf("unexpected error: %v", err)
		}
		if err != nil && tc.err != "" && err.Error() != tc.err {
			t.Errorf("incorrect error: '%v', want '%v'", err, tc.err)
		}
	}
}

func TestReverseZoneName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "1.2.3.4", expected: "3.2.1.in-addr.arpa."},
		{input: "fdda:5cc1:23:4::1f", expected: "f.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa."},
		{input: "fdda:5cc1:0023:0004:0000:0000:0000:001f", expected: "f.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa."},
	}

	for _, tc := range testCases {
		ip, err := ParseIP(tc.input)
		if err != nil {
			t.Fatalf("Failed to parse IP: %v", err)
		}
		actual := ip.ReverseZoneName()
		if actual != tc.expected {
			t.Errorf("incorrect name: '%s', want '%s'", actual, tc.expected)
		}
	}
}

func TestPTRRecordValue(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "1.2.3.4", expected: "4"},
		{input: "fdda:5cc1:23:4::1f", expected: "4.0.0.0.3.2.0.0.1.c.c.5.a.d.d.f"},
		{input: "fdda:5cc1:0023:0004:0000:0000:0000:001f", expected: "4.0.0.0.3.2.0.0.1.c.c.5.a.d.d.f"},
	}

	for _, tc := range testCases {
		ip, err := ParseIP(tc.input)
		if err != nil {
			t.Fatalf("Failed to parse IP: %v", err)
		}
		actual := ip.PTRRecordValue()
		if actual != tc.expected {
			t.Errorf("incorrect name: '%s', want '%s'", actual, tc.expected)
		}
	}
}
