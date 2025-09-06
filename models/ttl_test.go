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

func TestString_TTL(t *testing.T) {
	ttl := &TTL{Value: toInt32Ptr(99), Comment: "comment"}
	want := "TTL{ Value: 99, Comment: comment }"

	if ttl.String() != want {
		t.Errorf("incorrect string: '%s', want: '%s'", ttl.String(), want)
	}

	ttl = &TTL{}
	want = "TTL{ Value: <nil>, Comment:  }"

	if ttl.String() != want {
		t.Errorf("incorrect string: '%s', want: '%s'", ttl.String(), want)
	}
}

func TestRender_TTL(t *testing.T) {
	testCases := []struct {
		ttl  *TTL
		want string
	}{
		{ttl: &TTL{Value: toInt32Ptr(150), Comment: "render comment"}, want: "$TTL 150 ;render comment"},
		{ttl: &TTL{Value: toInt32Ptr(999999)}, want: "$TTL 999999"},
		{ttl: &TTL{Value: nil}, want: ""},
		{ttl: &TTL{Value: nil, Comment: "doesn't matter if this is here"}, want: ""},
	}

	for _, tc := range testCases {
		if tc.ttl.Render() != tc.want {
			t.Errorf("incorrect string: '%s', want: '%s'", tc.ttl.Render(), tc.want)
		}
	}
}
