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
package grpc

import (
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTTLFromProtoBuf(t *testing.T) {
	testCases := []struct {
		ttl   *models.TTL
		proto *proto.TTL
	}{
		{ttl: nil, proto: nil},
		{ttl: &models.TTL{}, proto: nil},
		{ttl: nil, proto: &proto.TTL{}},
		{
			ttl:   &models.TTL{Value: toInt32Ptr(99), Comment: "testing-comment"},
			proto: &proto.TTL{Ttl: toInt32Ptr(99), Comment: "testing-comment"},
		},
	}

	for _, tc := range testCases {
		input := &models.TTL{}
		want := tc.ttl
		if tc.ttl == nil || tc.proto == nil {
			input = nil
			want = nil
		}
		TTLFromProtoBuf(tc.proto, input)

		if !cmp.Equal(input, want) {
			t.Errorf("incorrect result: %s, want: %s", input, want)
		}
	}
}

func TestTTLToProtoBuf(t *testing.T) {
	testCases := []struct {
		ttl   *models.TTL
		proto *proto.TTL
	}{
		{ttl: nil, proto: nil},
		{
			ttl:   &models.TTL{Value: toInt32Ptr(99), Comment: "testing-comment"},
			proto: &proto.TTL{Ttl: toInt32Ptr(99), Comment: "testing-comment"},
		},
	}

	for _, tc := range testCases {
		result := TTLToProtoBuf(tc.ttl)

		if !cmp.Equal(result, tc.proto, cmpopts.IgnoreUnexported(proto.TTL{})) {
			t.Errorf("incorrect result: %s, want: %s", result, tc.proto)
		}
	}
}
