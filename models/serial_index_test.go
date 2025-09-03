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

import "testing"

func TestSerial(t *testing.T) {
	testCases := []struct {
		base        *uint32
		changeIndex *uint32
		want        string
	}{
		{want: ""},
		{base: toUint32Ptr(1), want: ""},
		{changeIndex: toUint32Ptr(1), want: ""},
		{base: toUint32Ptr(0), changeIndex: toUint32Ptr(0), want: "00"},
		{base: toUint32Ptr(1), changeIndex: toUint32Ptr(1), want: "11"},
		{base: toUint32Ptr(20250903), changeIndex: toUint32Ptr(1), want: "202509031"},
	}

	for _, tc := range testCases {
		si := &SerialIndex{Base: tc.base, ChangeIndex: tc.changeIndex}

		if si.Serial() != tc.want {
			t.Errorf("incorrect serial: '%s', want: '%s'", si.Serial(), tc.want)
		}
	}
}
