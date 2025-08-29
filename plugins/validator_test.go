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

	"github.com/bcurnow/zonemgr/models"
)

var validations Validator = V()

func TestCommonValidations(t *testing.T) {
	testCases := []struct {
		rr         *models.ResourceRecord
		name       string
		identifier string
		err        string
	}{
		{identifier: "ValidRecord", rr: &models.ResourceRecord{Type: models.A, Value: "1.2.3.4"}, name: "ValidRecordWithoutAName"},
		{identifier: "Wrong resource record type for plugin", rr: &models.ResourceRecord{Type: models.CNAME, Value: "1.2.3.4"}, err: "this plugin does not handle resource records of type 'CNAME' only '[A]', identifier: 'Wrong resource record type for plugin'"},
		{identifier: "Wrong class", rr: &models.ResourceRecord{Type: models.A, Class: "bogus", Value: "1.2.3.4"}, err: "invalid A record, 'bogus' is not a valid class, identifier: 'Wrong class'"},
		{identifier: "Value set twice", rr: &models.ResourceRecord{Type: models.A, Values: []*models.ResourceRecordValue{{Value: "value once"}}, Value: "value again"}, err: "invalid A record, both value and values are set, identifier: 'Value set twice'"},
		{identifier: "Comment set twice", rr: &models.ResourceRecord{Type: models.A, Values: []*models.ResourceRecordValue{{Comment: "comment again"}}, Comment: "comment once"}, err: "invalid A record, both comment and values are set, identifier: 'Comment set twice'"},
	}

	for _, tc := range testCases {
		err := validations.CommonValidations(tc.identifier, tc.rr, A)
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
		rrType         models.ResourceRecordType
		supportedTypes []PluginType
		err            string
	}{
		{rrType: models.A, supportedTypes: []PluginType{A}},
		{rrType: models.A, supportedTypes: []PluginType{CNAME, NS}, err: "this plugin does not handle resource records of type 'A' only '[CNAME NS]', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		if err := validations.IsSupportedPluginType("testing", tc.rrType, tc.supportedTypes...); err != nil {
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
		{name: longRecordName, err: fmt.Sprintf("invalid A record, must be less than 255 characters: '%s', identifier:'testing'", longRecordName)},
		{name: "1isvalid"},
		{name: "withhyphenstart.-valid", err: "invalid A record, cannot start or end with a hyphen (-): 'withhyphenstart.-valid', identifier:'testing'"},
		{name: "withhyphenend.valid-", err: "invalid A record, cannot start or end with a hyphen (-): 'withhyphenend.valid-', identifier:'testing'"},
	}

	for _, tc := range testCases {
		if err := validations.IsValidRFC1035Name("testing", tc.name, models.A); err != nil {
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
	if err := validations.IsValidNameOrWildcard("testing", "@", models.A); err != nil {
		t.Errorf("unexpected error")
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
		{email: "bogus@example.com@example.com", err: "invalid A record, invalid email address: 'bogus@example.com@example.com', identifier:'testing'"},
	}

	for _, tc := range testCases {
		formattedEmail, err := validations.FormatEmail("testing", tc.email, models.A)
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
		{name: "name.domain.com", err: "invalid A record, must end with a trailing dot: 'name.domain.com', identifier: 'testing'"},
		{name: "name.", err: "invalid A record, must be fully qualified with at least two dots: 'name.', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		err := validations.IsFullyQualified("testing", tc.name, models.A)
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
