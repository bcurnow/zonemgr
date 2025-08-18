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

package plugins

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bcurnow/zonemgr/schema"
)

func TestStandardValidations(t *testing.T) {
	testCases := []struct {
		rr         *schema.ResourceRecord
		name       string
		identifier string
		err        string
	}{
		{identifier: "ValidRecord", rr: &schema.ResourceRecord{Type: schema.A, Value: "1.2.3.4"}, name: "ValidRecordWithoutAName"},
		{identifier: "Wrong resource record type for plugin", rr: &schema.ResourceRecord{Type: schema.CNAME, Value: "1.2.3.4"}, err: "this plugin does not handle resource records of type 'CNAME' only '[A]', identifier: 'Wrong resource record type for plugin'"},
		{identifier: "Wrong class", rr: &schema.ResourceRecord{Type: schema.A, Class: "bogus", Value: "1.2.3.4"}, err: "invalid A record, 'bogus' is not a valid class, identifier: 'Wrong class'"},
		{identifier: "Value set twice", rr: &schema.ResourceRecord{Type: schema.A, Value: "and again", Values: []*schema.ResourceRecordValue{{Value: "and again"}}}, err: "invalid A record, can not specify both value and values, identifier: 'Value set twice'"},
		{identifier: "Comment set twice", rr: &schema.ResourceRecord{Type: schema.A, Values: []*schema.ResourceRecordValue{{Comment: "once"}}, Comment: "once"}, err: "invalid A record, can not specify both comment and values, identifier: 'Comment set twice'"},
	}

	for _, tc := range testCases {
		err := StandardValidations(tc.identifier, tc.rr, []PluginType{RecordA})
		if err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		}
	}
}

func TestIsSupportedPluginType(t *testing.T) {
	testCases := []struct {
		rrType         schema.ResourceRecordType
		supportedTypes []PluginType
		err            string
	}{
		{rrType: schema.A, supportedTypes: []PluginType{RecordA}},
		{rrType: schema.A, supportedTypes: []PluginType{RecordCNAME, RecordNS}, err: "this plugin does not handle resource records of type 'A' only '[CNAME NS]', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		if err := IsSupportedPluginType("testing", tc.rrType, tc.supportedTypes); err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("expected error")
			}
		}
	}
}

func TestIsValidRFC1035Name(t *testing.T) {
	longRecordName := "MoreThan255Character" + strings.Repeat("s", 255)

	testCases := []struct {
		name string
		err  string
	}{
		{name: "valid"},
		{name: longRecordName, err: fmt.Sprintf("invalid A record, must be less than 255 characters: '%s', identifier=testing", longRecordName)},
		{name: "1invalid", err: fmt.Sprintf("invalid A record, does not match regexp '%s': '1invalid', identifier=testing", `^([A-Za-z])([A-Za-z0-9-]{1,62})(\.[A-Za-z0-9-]{1,63})*\.{0,1}$`)},
		{name: "withhyphenstart.-valid", err: "invalid A record, can not start or end with a hyphen (-): 'withhyphenstart.-valid', identifier=testing"},
		{name: "withhyphenend.valid-", err: "invalid A record, can not start or end with a hyphen (-): 'withhyphenend.valid-', identifier=testing"},
	}

	for _, tc := range testCases {
		if err := IsValidRFC1035Name("testing", tc.name, schema.A); err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("expected error")
			}
		}
	}
}

func TestIValidNameOrWildcard(t *testing.T) {
	// We just need to valid the @ case, the others are already tested
	if err := IsValidNameOrWildcard("testing", "@", schema.A); err != nil {
		t.Errorf("unexpected exception")
	}
}

func TestFormatEmail(t *testing.T) {

	testCases := []struct {
		email string
		err   string
		want  string
	}{
		{email: "name@example.com", want: "name.example.com."},
		{email: "name.example.com.", want: "name.example.com."},
		{email: "name.example.com", want: "name.example.com."},
		{email: "bogus@example.com@example.com", err: "invalid A record, invalid email address: bogus@example.com@example.com, identifier=testing"},
	}

	for _, tc := range testCases {
		formattedEmail, err := FormatEmail("testing", tc.email, schema.A)
		if err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		} else {
			if formattedEmail != tc.want {
				t.Errorf("incorrect email: %s, want %s", formattedEmail, tc.want)
			}
		}
	}
}

func TestIsFullyQualified(t *testing.T) {

	testCases := []struct {
		name string
		err  string
	}{
		{name: "name.domain.com."},
		{name: "name.domain.com", err: "invalid A record, must end with a trailing dot: 'name.domain.com', identifier=testing"},
		{name: "name.", err: "invalid A record, must be fully qualified with at least two dots: 'name.', identifier=testing"},
	}

	for _, tc := range testCases {
		err := IsFullyQualified("testing", tc.name, schema.A)
		if err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("expected error")
			}
		}
	}
}
