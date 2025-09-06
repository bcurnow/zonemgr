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

func TestResourceRecordFromProtoBuf(t *testing.T) {
	testCases := []struct {
		rr    *models.ResourceRecord
		proto *proto.ResourceRecord
	}{
		{rr: nil, proto: nil},
		{rr: nil, proto: &proto.ResourceRecord{}},
		{rr: &models.ResourceRecord{}, proto: nil},
		{rr: &models.ResourceRecord{}, proto: &proto.ResourceRecord{}},
		{
			rr: &models.ResourceRecord{
				Type:    models.A,
				Name:    "testing",
				Class:   models.INTERNET,
				TTL:     toInt32Ptr(99),
				Values:  []*models.ResourceRecordValue{{Value: "test-values", Comment: "test-values-comment"}},
				Value:   "test-value",
				Comment: "test-comment",
			},
			proto: &proto.ResourceRecord{
				Type:    "A",
				Name:    "testing",
				Class:   "IN",
				Ttl:     toInt32Ptr(99),
				Values:  []*proto.ResourceRecordValue{{Value: "test-values", Comment: "test-values-comment"}},
				Value:   "test-value",
				Comment: "test-comment",
			},
		},
	}

	for _, tc := range testCases {
		input := &models.ResourceRecord{}
		if tc.rr == nil {
			input = nil
		}

		ResourceRecordFromProtoBuf(tc.proto, input)

		if !cmp.Equal(input, tc.rr) {
			t.Errorf("incorrect result: %s, want: %s", input, tc.rr)
		}
	}
}

func TestUpdateResourceRecordToProtoBuf(t *testing.T) {
	testCases := []struct {
		rr    *models.ResourceRecord
		proto *proto.ResourceRecord
	}{
		{rr: nil, proto: nil},
		{
			rr: &models.ResourceRecord{
				Type:    models.A,
				Name:    "testing",
				Class:   models.INTERNET,
				TTL:     toInt32Ptr(99),
				Values:  []*models.ResourceRecordValue{{Value: "test-values", Comment: "test-values-comment"}},
				Value:   "test-value",
				Comment: "test-comment",
			},
			proto: &proto.ResourceRecord{
				Type:    "A",
				Name:    "testing",
				Class:   "IN",
				Ttl:     toInt32Ptr(99),
				Values:  []*proto.ResourceRecordValue{{Value: "test-values", Comment: "test-values-comment"}},
				Value:   "test-value",
				Comment: "test-comment",
			},
		},
		{
			rr: &models.ResourceRecord{
				Type:    models.A,
				Name:    "testing",
				Class:   models.INTERNET,
				TTL:     toInt32Ptr(99),
				Values:  make([]*models.ResourceRecordValue, 0),
				Value:   "test-value",
				Comment: "test-comment",
			},
			proto: &proto.ResourceRecord{
				Type:    "A",
				Name:    "testing",
				Class:   "IN",
				Ttl:     toInt32Ptr(99),
				Value:   "test-value",
				Comment: "test-comment",
			},
		},
	}

	for _, tc := range testCases {
		result := ResourceRecordToProtoBuf(tc.rr)

		if !cmp.Equal(result, tc.proto, cmpopts.IgnoreTypes(proto.ResourceRecord{}, proto.ResourceRecordValue{})) {
			t.Errorf("incorrect result: %s, want: %s", result, tc.proto)
		}
	}
}
