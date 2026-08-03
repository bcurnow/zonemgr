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
	"slices"
	"strings"
	"testing"

	"github.com/bcurnow/zonemgr/models"
)

func txtRenderPrefix(name string) string {
	return fmt.Sprintf(models.ResourceRecordNameFormatString, name) + " " + fmt.Sprintf(models.ResourceRecordTypeFormatString, models.TXT) + " "
}

func TestTXTNormalize(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: "v=spf1 -all"},
		},
		{
			name:       "name-defaulting",
			identifier: "example.com.",
			rr:         &models.ResourceRecord{Type: models.TXT, Value: "v=spf1 -all"},
		},
		{
			name:       "wildcard-name",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "@", Value: "v=spf1 -all"},
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.A, Name: "example.com.", Value: "v=spf1 -all"},
			wantErr:    "this plugin does not handle resource records of type 'A' only '[TXT]', identifier: 'record1'",
		},
		{
			name:       "invalid-name",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "-invalid", Value: "v=spf1 -all"},
			wantErr:    "invalid TXT record, cannot start or end with a hyphen (-): '-invalid', identifier: 'record1'",
		},
		{
			// <character-string> is binary data per RFC1035 3.3, not restricted to printable ASCII; non-printable
			// or non-ASCII bytes are legal. Escaping (if any is needed) is the caller's responsibility, not
			// something Normalize validates.
			name:       "non-printable-and-non-ascii-values-are-accepted",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: "bad\tcafé"},
		},
		{
			name:       "multiple-values",
			identifier: "record1",
			rr: &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Values: []*models.ResourceRecordValue{
				{Value: "part one"},
				{Value: "part two"},
			}},
		},
		{
			name:       "value-entry-exceeds-255-characters",
			identifier: "record1",
			rr: &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Values: []*models.ResourceRecordValue{
				{Value: strings.Repeat("a", 256)},
			}},
			wantErr: "invalid TXT record, character-string exceeds 255 characters: '" + strings.Repeat("a", 256) + "', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkErr(t, (&BuiltinPluginTXT{}).Normalize(tc.identifier, tc.rr), tc.wantErr)
		})
	}
}

func TestTXTRender(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		rr         *models.ResourceRecord
		want       string
		wantErr    string
	}{
		{
			name:       "valid",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: "v=spf1 -all"},
			want:       txtRenderPrefix("example.com.") + `"v=spf1 -all"`,
		},
		{
			// A bare, unescaped quote is the one thing Render must still escape itself, since otherwise it
			// would prematurely terminate our own generated quoted string.
			name:       "escapes-a-bare-quote",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: `say "hi"`},
			want:       txtRenderPrefix("example.com.") + `"say \"hi\""`,
		},
		{
			// The caller is responsible for escaping \ and " themselves; an already-escaped quote or backslash
			// is left untouched rather than being escaped again (which would corrupt it, e.g. \" becoming \\").
			name:       "leaves-already-escaped-sequences-untouched",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: `say \"hi\" and a\\b and \068`},
			want:       txtRenderPrefix("example.com.") + `"say \"hi\" and a\\b and \068"`,
		},
		{
			// Non-printable/non-ASCII bytes the caller left un-escaped are passed through as-is; escaping is
			// entirely the caller's responsibility, not something Render does on their behalf.
			name:       "passes-through-non-printable-and-non-ascii-bytes-unescaped",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: "bad\tcafé"},
			want:       txtRenderPrefix("example.com.") + "\"bad\tcafé\"",
		},
		{
			name:       "with-comment",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: "v=spf1 -all", Comment: "SPF record"},
			want:       txtRenderPrefix("example.com.") + `"v=spf1 -all" ;SPF record`,
		},
		{
			name:       "wrong-type",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.A, Name: "example.com."},
			wantErr:    "this plugin does not handle resource records of type 'A' only '[TXT]', identifier: 'record1'",
		},
		{
			name:       "splits-values-over-255-characters",
			identifier: "record1",
			rr:         &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Value: strings.Repeat("a", 256)},
			want:       txtRenderPrefix("example.com.") + `"` + strings.Repeat("a", 255) + `" "a"`,
		},
		{
			name:       "renders-each-values-entry-as-its-own-character-string",
			identifier: "record1",
			rr: &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Values: []*models.ResourceRecordValue{
				{Value: "part one"},
				{Value: "part two"},
			}},
			want: txtRenderPrefix("example.com.") + `"part one" "part two"`,
		},
		{
			// Render must reject an oversized values entry itself rather than relying on Normalize having
			// already run and silently splitting it.
			name:       "errors-on-oversized-values-entry-without-normalize",
			identifier: "record1",
			rr: &models.ResourceRecord{Type: models.TXT, Name: "example.com.", Values: []*models.ResourceRecordValue{
				{Value: strings.Repeat("a", 256)},
			}},
			wantErr: "invalid TXT record, character-string exceeds 255 characters: '" + strings.Repeat("a", 256) + "', identifier: 'record1'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := (&BuiltinPluginTXT{}).Render(tc.identifier, tc.rr)
			checkErr(t, err, tc.wantErr)
			if tc.wantErr == "" && got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestEscapeTXTValue(t *testing.T) {
	testCases := []struct {
		name  string
		value string
		want  string
	}{
		{name: "empty", value: "", want: ""},
		{name: "plain-ascii", value: "hello world", want: "hello world"},
		{name: "bare-quote-is-escaped", value: `say "hi"`, want: `say \"hi\"`},
		{name: "already-escaped-quote-is-untouched", value: `say \"hi\"`, want: `say \"hi\"`},
		{name: "already-escaped-backslash-is-untouched", value: `a\\b`, want: `a\\b`},
		{name: "already-escaped-decimal-value-is-untouched", value: `\068`, want: `\068`},
		{name: "trailing-lone-backslash-is-untouched", value: `a\`, want: `a\`},
		{name: "non-printable-and-non-ascii-bytes-pass-through", value: "a\tcafé", want: "a\tcafé"},
		{
			// An escaped backslash (\\) immediately followed by a quote is a complete, self-contained escape
			// pair (one literal backslash) followed by a genuinely bare quote, which still needs escaping,
			// unlike a single \ immediately before a quote (\"), which is already a complete escaped quote.
			name:  "bare-quote-after-escaped-backslash-is-still-escaped",
			value: `a\\"b`,
			want:  `a\\\"b`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := escapeTXTValue(tc.value)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSplitTXTValue(t *testing.T) {
	testCases := []struct {
		name  string
		value string
		want  []string
	}{
		{name: "empty", value: "", want: []string{""}},
		{name: "under-limit", value: "hello", want: []string{"hello"}},
		{name: "exactly-at-limit", value: strings.Repeat("a", 255), want: []string{strings.Repeat("a", 255)}},
		{
			name:  "one-over-limit",
			value: strings.Repeat("a", 256),
			want:  []string{strings.Repeat("a", 255), "a"},
		},
		{
			name:  "multiple-chunks",
			value: strings.Repeat("a", 255) + strings.Repeat("b", 255) + "c",
			want:  []string{strings.Repeat("a", 255), strings.Repeat("b", 255), "c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := splitTXTValue(tc.value)
			if !slices.Equal(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
