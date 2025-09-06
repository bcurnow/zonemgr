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

package serial

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	testSerialString = "12345678"
	testSerial       = toUint32Ptr(12345678)
)

func setup_SerialNumberGenerator(_ *testing.T) {
	sprintf = func(format string, a ...any) string {
		if len(a) == 2 {
			return fmt.Sprintf("%s%02d", testSerialString, a[1].(uint32))
		}
		return testSerialString
	}
}

func teardown_SerialNumberGenerator(_ *testing.T) {
	sprintf = fmt.Sprintf
	parseUint = strconv.ParseUint
}

func TestGenerateBase(t *testing.T) {
	setup_SerialNumberGenerator(t)
	defer teardown_SerialNumberGenerator(t)
	serial, err := (&TimeBasedGenerator{}).GenerateBase()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !cmp.Equal(serial, testSerial) {
		t.Errorf("incorrect result: %d, want: %s", *serial, testSerialString)
	}
}

func TestGenerate(t *testing.T) {
	setup_SerialNumberGenerator(t)
	defer teardown_SerialNumberGenerator(t)
	testCases := []struct {
		index uint32
		want  *uint32
	}{
		{index: uint32(32), want: toUint32Ptr(1234567832)},
		{index: uint32(1), want: toUint32Ptr(1234567801)},
	}

	for _, tc := range testCases {

		serial, err := (&TimeBasedGenerator{}).Generate(tc.index)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !cmp.Equal(serial, tc.want) {
			t.Errorf("incorrect result: %d, want: %d", *serial, *tc.want)
		}
	}
}

func TestGenerate_Error(t *testing.T) {
	setup_SerialNumberGenerator(t)
	defer teardown_SerialNumberGenerator(t)
	parseUint = func(_ string, _, _ int) (uint64, error) {
		return uint64(0), errors.New("testing")
	}
	_, err := (&TimeBasedGenerator{}).Generate(uint32(1))
	if err == nil {
		t.Error("expected an error, found none")
	} else {
		want := "unable to generate a serial number from string '12345678': testing"
		if err.Error() != want {
			t.Errorf("incorrect error: '%s', want: '%s'", err, want)
		}
	}
}

func TestFromString(t *testing.T) {
	setup_SerialNumberGenerator(t)
	defer teardown_SerialNumberGenerator(t)
	testCases := []struct {
		strconvErr bool
	}{
		{},
		{strconvErr: true},
	}

	for _, tc := range testCases {
		input := testSerialString
		if tc.strconvErr {
			input = "not an integer string"
		}

		serial, err := (&TimeBasedGenerator{}).FromString(input)
		if err != nil {
			if tc.strconvErr {
				want := "unable to generate a serial number from string 'not an integer string': strconv.ParseUint: parsing \"not an integer string\": invalid syntax"
				if err.Error() != want {
					t.Errorf("incorrect error: '%s', want: '%s'", err, want)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if tc.strconvErr {
				t.Error("expected an error, found none")
			} else {
				if *serial != *testSerial {
					t.Errorf("incorrect result: %d, want: %d", *serial, *testSerial)
				}
			}
		}
	}
}

func toUint32Ptr(i uint32) *uint32 {
	return &i
}
