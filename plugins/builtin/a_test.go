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
	"fmt"
	"testing"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
	"github.com/bcurnow/zonemgr/test"
	"github.com/bcurnow/zonemgr/version"
	"github.com/google/go-cmp/cmp"
)

var testAPlugin = &APlugin{}

func TestPluginVersion(t *testing.T) {
	ver, err := testAPlugin.PluginVersion()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if ver != version.Version() {
		t.Errorf("incorrect version %s, want %s", ver, version.Version())
	}
}

func TestPluginTypes(t *testing.T) {
	pluginTypes, err := testAPlugin.PluginTypes()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := []plugins.PluginType{plugins.RecordA}
	if !cmp.Equal(pluginTypes, want) {
		t.Errorf("unexpected plugin types %s, want %s", pluginTypes, want)
	}
}

func TestConfigure(t *testing.T) {
	if err := testAPlugin.Configure(nil); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestNormalize(t *testing.T) {
	testCases := []struct {
		rr         *schema.ResourceRecord
		name       string
		identifier string
		err        string
	}{
		{identifier: "ValidRecordWithoutAName", rr: &schema.ResourceRecord{Type: schema.A, Value: "1.2.3.4"}, name: "ValidRecordWithoutAName"},
		{identifier: "Valid record with a name", rr: &schema.ResourceRecord{Type: schema.A, Value: "1.2.3.4", Name: "name"}, name: "name"},
		{identifier: "Wrong resource record type for plugin", rr: &schema.ResourceRecord{Type: schema.CNAME, Value: "1.2.3.4"}, err: "this plugin does not handle resource records of type 'CNAME' only '[A]', identifier: 'Wrong resource record type for plugin'"},
		{identifier: "Wrong class", rr: &schema.ResourceRecord{Type: schema.A, Class: "bogus", Value: "1.2.3.4"}, err: "invalid A record, 'bogus' is not a valid class, identifier: 'Wrong class'"},
		{identifier: "Value set twice", rr: &schema.ResourceRecord{Type: schema.A, Value: "and again", Values: []*schema.ResourceRecordValue{{Value: "and again"}}}, err: "invalid A record, can not specify both value and values, identifier: 'Value set twice'"},
		{identifier: "Comment set twice", rr: &schema.ResourceRecord{Type: schema.A, Values: []*schema.ResourceRecordValue{{Comment: "once"}}, Comment: "once"}, err: "invalid A record, can not specify both comment and values, identifier: 'Comment set twice'"},
		{identifier: "Invalid name", rr: &schema.ResourceRecord{Type: schema.A, Name: "1invalidname"}, err: fmt.Sprintf("invalid A record, does not match regexp '%s': '1invalidname', identifier=Invalid name", `^([A-Za-z])([A-Za-z0-9-]{1,62})(\.[A-Za-z0-9-]{1,63})*\.{0,1}$`)},
		{identifier: "Value not IP", rr: &schema.ResourceRecord{Type: schema.A, Value: "notIP", Name: "validname"}, err: "invalid A record, 'notIP' must be a valid IP address, identifier: 'Value not IP'"},
	}

	for _, tc := range testCases {
		err := testAPlugin.Normalize(tc.identifier, tc.rr)
		if err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		} else {
			// make sure the name defaulting worked
			if tc.rr.Name != tc.name {
				t.Errorf("incorrect name: %s, want %s", tc.rr.Name, tc.name)
			}
		}
	}
}

func TestValidateZone(t *testing.T) {
	if err := testAPlugin.ValidateZone("noop", nil); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestRender(t *testing.T) {
	testCases := []struct {
		rr             *schema.ResourceRecord
		supportedTypes []plugins.PluginType
		want           string
		err            string
	}{
		{rr: &schema.ResourceRecord{Name: "test", Type: schema.A, Value: "1.2.3.4"}, want: fmt.Sprintf("%-40s %-6s %s", "test", schema.A, "1.2.3.4")},
		{rr: &schema.ResourceRecord{Name: "test", Type: schema.A, Value: "1.2.3.4", Class: schema.INTERNET, TTL: test.ToInt32Ptr(30), Comment: "test comment"}, want: fmt.Sprintf("%-40s %-6s %s %d %s ;%s", "test", schema.A, schema.INTERNET, 30, "1.2.3.4", "test comment")},
		{rr: &schema.ResourceRecord{Name: "test", Type: schema.A, Values: []*schema.ResourceRecordValue{{Value: "1.2.3.4", Comment: "test comment"}}, Class: schema.INTERNET, TTL: test.ToInt32Ptr(30)}, want: fmt.Sprintf("%-40s %-6s %s %d %s ;%s", "test", schema.A, schema.INTERNET, 30, "1.2.3.4", "test comment")},
	}

	for _, tc := range testCases {
		s, err := testAPlugin.Render("testing", tc.rr)
		if err != nil {
			if tc.err == "" {
				t.Errorf("unexpected error: %s", err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("incorrect error: %s, want %s", err, tc.err)
				}
			}
		} else {
			if s != tc.want {
				t.Errorf("incorrect value: %s, want %s", s, tc.want)
			}
		}
	}
}
