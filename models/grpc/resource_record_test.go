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
	"reflect"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

func TestUpdateResourceRecordFromProtoBuf(t *testing.T) {
	testCases := []struct {
		rr    *models.ResourceRecord
		proto *proto.ResourceRecord
	}{
		{rr: nil, proto: nil},
		{rr: nil, proto: &proto.ResourceRecord{}},
		{rr: &models.ResourceRecord{}, proto: nil},
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

		if !reflect.DeepEqual(input, tc.rr) {
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
	}

	for _, tc := range testCases {
		result := ResourceRecordToProtoBuf(tc.rr)

		if !reflect.DeepEqual(result, tc.proto) {
			t.Errorf("incorrect result: %s, want: %s", result, tc.proto)
		}
	}
}
