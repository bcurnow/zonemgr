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

func TestEnsureSupportedPluginType(t *testing.T) {
	testCases := []struct {
		rrType         models.ResourceRecordType
		supportedTypes []PluginType
		err            string
	}{
		{rrType: models.A, supportedTypes: []PluginType{A}},
		{rrType: models.A, supportedTypes: []PluginType{CNAME, NS}, err: "this plugin does not handle resource records of type 'A' only '[CNAME NS]', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		if err := validations.EnsureSupportedPluginType("testing", tc.rrType, tc.supportedTypes...); err != nil {
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

func TestEnsureValidRFC1035Name(t *testing.T) {
	longRecordName := "MoreThan255Character" + strings.Repeat("s", 255)

	testCases := []struct {
		name string
		err  string
	}{
		{name: "valid"},
		{name: longRecordName, err: fmt.Sprintf("invalid A record, must be less than 255 characters: '%s', identifier: 'testing'", longRecordName)},
		{name: "1isvalid"},
		{name: "$bogus", err: fmt.Sprintf("invalid A record, does not match regexp '%s': '$bogus', identifier: 'testing'", dnsNameRegexRFC1035String)},
		{name: "withhyphenstart.-valid", err: "invalid A record, cannot start or end with a hyphen (-): 'withhyphenstart.-valid', identifier: 'testing'"},
		{name: "withhyphenend.valid-", err: "invalid A record, cannot start or end with a hyphen (-): 'withhyphenend.valid-', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		if err := validations.EnsureValidRFC1035Name("testing", tc.name, models.A); err != nil {
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

func TestEnsureValidNameOrWildcard(t *testing.T) {
	// We just need to valid the @ case and anything that will cause IsValidRFC1035Name to fail
	// the others are tested elsewhere
	if err := validations.EnsureValidNameOrWildcard("testing", "@", models.A); err != nil {
		t.Errorf("unexpected error")
	}

	if err := validations.EnsureValidNameOrWildcard("testing", "$bogus", models.A); err != nil {
		want := fmt.Sprintf("invalid A record, does not match regexp '%s': '$bogus', identifier: 'testing'", dnsNameRegexRFC1035String)
		if err.Error() != want {
			t.Errorf("incorrrect error: '%s', want: '%s'", err, want)
		}
	} else {
		t.Error("expected an error, found none")
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
		{email: "name.example.com", want: "name.example.com.", err: "invalid A record, must end with a trailing dot: 'name.example.com', identifier: 'testing'"},
		{email: "bogus@example.com@example.com", err: "invalid A record, invalid email address: 'bogus@example.com@example.com', identifier: 'testing'"},
		{email: ".bogus@example.com", err: "invalid A record, invalid email address: '.bogus@example.com', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		formattedEmail, err := validations.FormatEmail("testing", tc.email, models.A)
		if err != nil {
			if tc.err == "" {
				t.Errorf("%s - unexpected error: %s", tc.email, err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("%s - incorrect error: %s, want %s", tc.email, err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("%s - expected an error, found none", tc.email)
			}
			if formattedEmail != tc.want {
				t.Errorf("%s - incorrect email: %s, want %s", tc.email, formattedEmail, tc.want)
			}
		}
	}
}

func TestEnsureFullyQualified(t *testing.T) {

	testCases := []struct {
		name string
		err  string
	}{
		{name: "name.domain.com."},
		{name: "$bogus", err: fmt.Sprintf("invalid A record, does not match regexp '%s': '$bogus', identifier: 'testing'", dnsNameRegexRFC1035String)},
		{name: "name.domain.com", err: "invalid A record, must end with a trailing dot: 'name.domain.com', identifier: 'testing'"},
		{name: "name.", err: "invalid A record, must be fully qualified with at least two dots: 'name.', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		err := validations.EnsureFullyQualified("testing", tc.name, models.A)
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

func TestEnsureIP(t *testing.T) {
	testCases := []struct {
		s   string
		err string
	}{
		{s: "1.2.3.4"},
		{s: "not an ip", err: "invalid A record, 'not an ip' must be a valid IP address, identifier: 'testing'"},
	}

	for _, tc := range testCases {
		if err := validations.EnsureIP("testing", tc.s, models.A); err != nil {
			if tc.err != "" {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: '%s', want: '%s'", err, tc.err)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if tc.err != "" {
				t.Error("expected an error, found none")
			}
		}
	}
}

func TestEnsureNotIP(t *testing.T) {
	testCases := []struct {
		s   string
		err string
	}{
		{s: "1.2.3.4", err: "invalid A record, '1.2.3.4' must not be an IP address, identifier: 'testing'"},
		{s: "not an ip"},
	}

	for _, tc := range testCases {
		if err := validations.EnsureNotIP("testing", tc.s, models.A); err != nil {
			if tc.err != "" {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: '%s', want: '%s'", err, tc.err)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if tc.err != "" {
				t.Error("expected an error, found none")
			}
		}
	}
}

func TestEnsurePositive(t *testing.T) {
	testCases := []struct {
		s   string
		err string
	}{
		{s: "bogus", err: "strconv.ParseInt: parsing \"bogus\": invalid syntax"},
		{s: "0"},
		{s: "1"},
		{s: "1234567890"},
		{s: "-1", err: "retry must not be less than 0 on a SOA record, was '-1', identifier: 'testing'"},
	}

	for _, tc := range testCases {
		if err := validations.EnsurePositive("testing", tc.s, "retry", models.SOA); err != nil {
			if tc.err != "" {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: '%s', want: '%s'", err, tc.err)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if tc.err != "" {
				t.Error("expected an error, found none")
			}
		}
	}
}
