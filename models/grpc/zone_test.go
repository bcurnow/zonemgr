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

func TestZoneFromProtoBuf(t *testing.T) {
	testCases := []struct {
		zone  *models.Zone
		proto *proto.Zone
	}{
		{zone: nil, proto: nil},
		{zone: &models.Zone{}, proto: nil},
		{zone: nil, proto: &proto.Zone{}},
		{zone: &models.Zone{}, proto: &proto.Zone{}},
		{
			zone: &models.Zone{
				Config: &models.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				ResourceRecords: map[string]*models.ResourceRecord{
					"one": {Type: models.A, Name: "one"},
					"two": {Type: models.NS, Name: "ns.example.com."},
				},
				TTL: &models.TTL{Value: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
			proto: &proto.Zone{
				Config: &proto.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				ResourceRecords: map[string]*proto.ResourceRecord{
					"one": {Type: "A", Name: "one"},
					"two": {Type: "NS", Name: "ns.example.com."},
				},
				Ttl: &proto.TTL{Ttl: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
		},
		{
			zone: &models.Zone{
				Config: &models.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				TTL:    &models.TTL{Value: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
			proto: &proto.Zone{
				Config: &proto.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				Ttl:    &proto.TTL{Ttl: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
		},
	}

	for _, tc := range testCases {
		input := &models.Zone{}
		want := tc.zone
		if tc.zone == nil || tc.proto == nil || tc.proto.ResourceRecords == nil {
			input = nil
			want = nil
		}

		ZoneFromProtoBuf(tc.proto, input)

		if !cmp.Equal(input, want, cmpopts.IgnoreUnexported(models.Zone{})) {
			t.Errorf("incorrect result:\n%s", cmp.Diff(input, want, cmpopts.IgnoreUnexported(models.Zone{})))
		}
	}
}

func TestZoneToProtoBuf(t *testing.T) {
	testCases := []struct {
		zone  *models.Zone
		proto *proto.Zone
	}{
		{zone: nil, proto: nil},
		{
			zone: &models.Zone{
				Config: &models.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				TTL:    &models.TTL{Value: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
			proto: &proto.Zone{
				Config:          &proto.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				ResourceRecords: make(map[string]*proto.ResourceRecord),
				Ttl:             &proto.TTL{Ttl: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
		},
		{
			zone: &models.Zone{
				Config: &models.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				ResourceRecords: map[string]*models.ResourceRecord{
					"one": {Type: models.A, Name: "one"},
					"two": {Type: models.NS, Name: "ns.example.com."},
				},
				TTL: &models.TTL{Value: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
			proto: &proto.Zone{
				Config: &proto.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing-scid"},
				ResourceRecords: map[string]*proto.ResourceRecord{
					"one": {Type: "A", Name: "one", Values: make([]*proto.ResourceRecordValue, 0)},
					"two": {Type: "NS", Name: "ns.example.com.", Values: make([]*proto.ResourceRecordValue, 0)},
				},
				Ttl: &proto.TTL{Ttl: toInt32Ptr(99), Comment: "testing-ttl-comment"},
			},
		},
	}

	for _, tc := range testCases {
		result := ZoneToProtoBuf(tc.zone)

		if !cmp.Equal(result, tc.proto, cmpopts.IgnoreUnexported(
			proto.Zone{},
			proto.Config{},
			proto.ResourceRecord{},
			proto.ResourceRecordValue{},
			proto.TTL{},
		)) {
			t.Errorf("incorrect result:\n%s", cmp.Diff(result, tc.proto, cmpopts.IgnoreUnexported(
				proto.Zone{},
				proto.Config{},
				proto.ResourceRecord{},
				proto.ResourceRecordValue{},
				proto.TTL{},
			)))
		}
	}
}
