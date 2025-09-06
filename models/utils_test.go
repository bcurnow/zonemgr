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
	"errors"
	"testing"
)

func TestInt32ToString(t *testing.T) {
	testCases := []struct {
		input *int32
		want  string
	}{
		{input: nil, want: "<nil>"},
		{input: toInt32Ptr(99), want: "99"},
	}

	for _, tc := range testCases {
		actual := int32ToString(tc.input)
		if actual != tc.want {
			t.Errorf("incorrect result: '%s', want: '%s'", actual, tc.want)
		}
	}
}

func TestUInt32ToString(t *testing.T) {
	testCases := []struct {
		input *uint32
		want  string
	}{
		{input: nil, want: "<nil>"},
		{input: toUint32Ptr(99), want: "99"},
	}

	for _, tc := range testCases {
		actual := uint32ToString(tc.input)
		if actual != tc.want {
			t.Errorf("incorrect result: '%s', want: '%s'", actual, tc.want)
		}
	}
}

func TestWithSortedZones(t *testing.T) {
	testCases := []struct {
		zones map[string]*Zone
		order []string
		fnErr bool
	}{
		{zones: map[string]*Zone{"1": {}, "2": {}, "3": {}}, order: []string{"1", "2", "3"}},
		{zones: map[string]*Zone{"3": {}, "1": {}, "2": {}}, order: []string{"1", "2", "3"}},
		{zones: map[string]*Zone{"3": {}, "1": {}, "2": {}}, fnErr: true},
	}

	for _, tc := range testCases {
		index := 0
		testFn := func(name string, zone *Zone) error {
			if tc.fnErr {
				return errors.New("fnErr")
			} else {
				if name != tc.order[index] {
					t.Errorf("incorrect zone: '%s', want: '%s'", name, tc.order[index])
				}
				index++
				return nil
			}
		}

		err := WithSortedZones(tc.zones, testFn)
		if err != nil {
			if tc.fnErr {
				want := "fnErr"
				if err.Error() != want {
					t.Errorf("incorrect error: '%s', want: '%s'", err, want)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if tc.fnErr {
				t.Error("expected an error, found none")
			}
		}
	}
}
