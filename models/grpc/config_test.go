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

func TestConfigFromProtoBuf(t *testing.T) {
	testCases := []struct {
		config *models.Config
		proto  *proto.Config
	}{
		{config: nil, proto: nil},
		{config: &models.Config{}, proto: nil},
		{config: &models.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing"}, proto: &proto.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing"}},
	}

	for _, tc := range testCases {
		inputConfig := &models.Config{}
		if tc.config == nil {
			inputConfig = nil
		}
		ConfigFromProtoBuf(tc.proto, inputConfig)

		if tc.config == nil {
			if inputConfig != nil {
				t.Errorf("expected nil config input to return nil")
			}
		} else {
			if !reflect.DeepEqual(inputConfig, tc.config) {
				t.Errorf("incorrect result: %s, want: %s", inputConfig, tc.config)
			}
		}
	}
}

func TestConfigToProtoBuf(t *testing.T) {
	testCases := []struct {
		proto  *proto.Config
		config *models.Config
	}{
		{config: nil, proto: &proto.Config{}},
		{config: &models.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing"}, proto: &proto.Config{GenerateSerial: true, GenerateReverseLookupZones: true, SerialChangeIndexDirectory: "testing"}},
	}

	for _, tc := range testCases {
		result := ConfigToProtoBuf(tc.config)

		if !reflect.DeepEqual(result, tc.proto) {
			t.Errorf("incorrect result: %s, want: %s", result, tc.proto)
		}
	}
}
