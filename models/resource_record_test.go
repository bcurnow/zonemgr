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
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestString_ResourceRecord(t *testing.T) {
	rr := ResourceRecord{
		Name:    "name",
		Type:    A,
		Class:   INTERNET,
		TTL:     toInt32Ptr(30),
		Values:  []*ResourceRecordValue{},
		Value:   "value",
		Comment: "comment",
	}
	want := "ResourceRecord{\n" +
		"       Name: name\n" +
		"       Type: A\n" +
		"       Class: IN\n" +
		"       TTL: 30\n" +
		"       Values: []\n" +
		"       Value: value\n" +
		"       Comment: comment\n" +
		"     }"

	if cmp.Diff(rr.String(), want) != "" {
		t.Errorf("incorrect string:\n%s", cmp.Diff(rr.String(), want))
	}

	rr = ResourceRecord{}
	want = "ResourceRecord{\n" +
		"       Name: \n" +
		"       Type: \n" +
		"       Class: \n" +
		"       TTL: <nil>\n" +
		"       Values: []\n" +
		"       Value: \n" +
		"       Comment: \n" +
		"     }"

	if cmp.Diff(rr.String(), want) != "" {
		t.Errorf("incorrect string:\n%s", cmp.Diff(rr.String(), want))
	}
}

func TestRetrieveSingleValue(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want string
	}{
		{rr: &ResourceRecord{}, want: ""},
		{rr: &ResourceRecord{Value: "value"}, want: "value"},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Value: "values"}}}, want: "values"},
		{rr: &ResourceRecord{Value: "value", Values: []*ResourceRecordValue{{Value: "values"}}}, want: "values"},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Value: "value1"}, {Value: "value2"}}}, want: "value1"},
	}

	for _, tc := range testCases {
		if tc.rr.RetrieveSingleValue() != tc.want {
			t.Errorf("incorrect value: '%s', want '%s'", tc.rr.RetrieveSingleValue(), tc.want)
		}
	}
}

func TestRetrieveSingleComment(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want string
	}{
		{rr: &ResourceRecord{}, want: ""},
		{rr: &ResourceRecord{Comment: "comment"}, want: "comment"},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Comment: "comments"}}}, want: "comments"},
		{rr: &ResourceRecord{Comment: "comment", Values: []*ResourceRecordValue{{Comment: "comments"}}}, want: "comments"},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Comment: "comment1"}, {Comment: "comment2"}}}, want: "comment1"},
	}

	for _, tc := range testCases {
		if tc.rr.RetrieveSingleComment() != tc.want {
			t.Errorf("incorrect comment: '%s', want '%s'", tc.rr.RetrieveSingleComment(), tc.want)
		}
	}
}

func TestIsValueSetInOnePlace(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want bool
	}{
		{rr: &ResourceRecord{}, want: true},
		{rr: &ResourceRecord{Value: "value"}, want: true},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Value: "values"}}}, want: true},
		{rr: &ResourceRecord{Value: "value", Values: []*ResourceRecordValue{{Value: "values"}}}, want: false},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Value: "value1"}, {Value: "value2"}}}, want: true},
	}

	for _, tc := range testCases {
		if tc.rr.IsValueSetInOnePlace() != tc.want {
			t.Errorf("incorrect value: '%v', want '%v'", tc.rr.IsValueSetInOnePlace(), tc.want)
		}
	}
}

func TestIsCommentSetInOnePlace(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want bool
	}{
		{rr: &ResourceRecord{}, want: true},
		{rr: &ResourceRecord{Comment: "comment"}, want: true},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Comment: "comments"}}}, want: true},
		{rr: &ResourceRecord{Comment: "comment", Values: []*ResourceRecordValue{{Comment: "comments"}}}, want: false},
		{rr: &ResourceRecord{Values: []*ResourceRecordValue{{Comment: "comment1"}, {Comment: "comment2"}}}, want: true},
	}

	for _, tc := range testCases {
		if tc.rr.IsCommentSetInOnePlace() != tc.want {
			t.Errorf("incorrect comment: '%v', want '%v'", tc.rr.IsCommentSetInOnePlace(), tc.want)
		}
	}
}

func TestRenderResourceWithoutValue(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want string
	}{
		// I don't like the way the wants are put together but haven't come up with a better idea
		{rr: &ResourceRecord{}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" ", "", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Value: "1.2.3.4"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" ", "name", "A")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "1.2.3.4"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s ", "name", "A", "IN", "30")},
	}

	for _, tc := range testCases {
		if tc.rr.RenderResourceWithoutValue() != tc.want {
			t.Errorf("incorrect render: '%s', want '%s'", tc.rr.RenderResourceWithoutValue(), tc.want)
		}
	}
}

func TestRenderSingleValueResource(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want string
	}{
		// I don't like the way the wants are put together but haven't come up with a better idea
		{rr: &ResourceRecord{}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" ", "", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Value: "1.2.3.4"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s", "name", "A", "1.2.3.4")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "1.2.3.4"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s %s", "name", "A", "IN", "30", "1.2.3.4")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "1.2.3.4", Comment: "testing"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s %s ;%s", "name", "A", "IN", "30", "1.2.3.4", "testing")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Values: []*ResourceRecordValue{{Value: "1.2.3.4"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s %s", "name", "A", "IN", "30", "1.2.3.4")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Values: []*ResourceRecordValue{{Value: "1.2.3.4", Comment: "testing"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s %s ;%s", "name", "A", "IN", "30", "1.2.3.4", "testing")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "main value", Values: []*ResourceRecordValue{{Value: "1.2.3.4"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s %s", "name", "A", "IN", "30", "1.2.3.4")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "main value", Comment: "main comment", Values: []*ResourceRecordValue{{Value: "1.2.3.4", Comment: "testing"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s %s ;%s", "name", "A", "IN", "30", "1.2.3.4", "testing")},
	}

	for _, tc := range testCases {
		if tc.rr.RenderSingleValueResource() != tc.want {
			t.Errorf("incorrect render: '%s', want '%s'", tc.rr.RenderSingleValueResource(), tc.want)
		}
	}
}

func TestRenderMultiValueResource(t *testing.T) {
	testCases := []struct {
		rr   *ResourceRecord
		want string
	}{
		// I don't like the way the wants are put together but haven't come up with a better idea
		{rr: &ResourceRecord{}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" (\n%48s)", "", "", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Value: "1.2.3.4"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" (\n%48s)", "name", "A", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "1.2.3.4"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s (\n%54s)", "name", "A", "IN", "30", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "1.2.3.4", Comment: "testing"}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s (\n%54s)", "name", "A", "IN", "30", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Values: []*ResourceRecordValue{{Value: "1.2.3.4"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s (\n%54s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+"\n%54s)", "name", "A", "IN", "30", "", "", "1.2.3.4", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Values: []*ResourceRecordValue{{Value: "1.2.3.4", Comment: "testing"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s (\n%54s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%54s)", "name", "A", "IN", "30", "", "", "1.2.3.4", "testing", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "main value", Values: []*ResourceRecordValue{{Value: "1.2.3.4"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s (\n%54s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+"\n%54s)", "name", "A", "IN", "30", "", "", "1.2.3.4", "")},
		{rr: &ResourceRecord{Name: "name", Type: A, Class: INTERNET, TTL: toInt32Ptr(30), Value: "main value", Comment: "main comment", Values: []*ResourceRecordValue{{Value: "1.2.3.4", Comment: "testing"}}}, want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" %s %s (\n%54s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%54s)", "name", "A", "IN", "30", "", "", "1.2.3.4", "testing", "")},
		{rr: &ResourceRecord{
			Name: "example.com.",
			Type: SOA,
			Values: []*ResourceRecordValue{
				{Value: "ns1.example.com.", Comment: "name server"},
				{Value: "admin@example.com", Comment: "admin"},
				{Value: "12345678", Comment: "serial"},
				{Value: "3", Comment: "refresh"},
				{Value: "4", Comment: "retry"},
				{Value: "5", Comment: "expire"},
				{Value: "6", Comment: "ncache"},
			},
		},
			want: fmt.Sprintf(ResourceRecordNameFormatString+" "+ResourceRecordTypeFormatString+" (\n%48s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%48s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%48s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%48s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%48s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%48s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%48s"+ResourceRecordMultivalueIndentFormatString+ResourceRecordNameFormatString+" ;%s\n%48s)", "example.com.", "SOA", "", "", "ns1.example.com.", "name server", "", "", "admin@example.com", "admin", "", "", "12345678", "serial", "", "", "3", "refresh", "", "", "4", "retry", "", "", "5", "expire", "", "", "6", "ncache", ""),
		},
	}

	for _, tc := range testCases {
		if cmp.Diff(tc.rr.RenderMultivalueResource(), tc.want) != "" {
			t.Errorf(`incorrect render:
			"""
			%s
			"""`, cmp.Diff(tc.rr.RenderMultivalueResource(), tc.want))
		}
	}
}
