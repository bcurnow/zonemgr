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

package grpc_plugins

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/models/grpc"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

func TestPluginVersion_Client(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err error
	}{
		{},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		call := mockClient.EXPECT().PluginVersion(context.Background(), emptyMessage)
		if tc.err != nil {
			call.Return(nil, tc.err)
		} else {
			call.Return(&proto.PluginVersionResponse{Version: "testing-ver"}, nil)
		}

		resp, err := grpcClient.PluginVersion()
		handleError(t, err, tc.err)
		if tc.err == nil {
			wanted := "testing-ver"
			if resp != wanted {
				t.Errorf("incorrect version: '%s', wanted: '%s'", resp, wanted)
			}
		}
	}
}

func TestPluginTypes_Client(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err error
	}{
		{},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		call := mockClient.EXPECT().PluginTypes(context.Background(), emptyMessage)
		if tc.err != nil {
			call.Return(nil, tc.err)
		} else {
			call.Return(&proto.PluginTypesResponse{SupportedTypes: []string{"A", "NS"}}, nil)
		}

		resp, err := grpcClient.PluginTypes()
		handleError(t, err, tc.err)
		if tc.err == nil {
			wanted := []plugins.PluginType{plugins.A, plugins.NS}
			if !slices.Equal(resp, wanted) {
				t.Errorf("incorrect version: '%s', wanted: '%s'", resp, wanted)
			}
		}
	}
}

func TestConfigure_Client(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err          error
		wantedConfig *models.Config
	}{
		{wantedConfig: &models.Config{GenerateSerial: true}},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		c := &models.Config{}
		call := mockClient.EXPECT().Configure(context.Background(), &proto.ConfigureRequest{Config: grpc.ConfigToProtoBufTo(c)})
		if tc.err != nil {
			call.Return(nil, tc.err)
		} else {
			call.Return(emptyMessage, nil)
		}

		err := grpcClient.Configure(c)
		handleError(t, err, tc.err)
		if reflect.DeepEqual(c, tc.wantedConfig) {
			t.Errorf("incorrect result: '%s', wanted: '%s'", c, tc.wantedConfig)
		}
	}
}

func TestNormalize_Client(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err      error
		wantedRR *models.ResourceRecord
	}{
		{wantedRR: &models.ResourceRecord{Type: models.A, Name: "testing"}},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		identifier := "testing"
		rr := &models.ResourceRecord{}
		call := mockClient.EXPECT().Normalize(context.Background(), &proto.NormalizeRequest{Identifier: identifier, ResourceRecord: grpc.ResourceRecordToProtoBuf(rr)})
		if tc.err != nil {
			call.Return(nil, tc.err)
		} else {
			call.Return(&proto.NormalizeResponse{ResourceRecord: grpc.ResourceRecordToProtoBuf(tc.wantedRR)}, nil)
		}

		err := grpcClient.Normalize(identifier, rr)
		handleError(t, err, tc.err)

		if tc.err == nil {
			if !reflect.DeepEqual(rr, tc.wantedRR) {
				t.Errorf("incorrect result: '%s', wanted: '%s'", rr, tc.wantedRR)
			}
		}
	}
}

func TestValidateZone_Client(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err        error
		wantedZone *models.Zone
	}{
		{wantedZone: &models.Zone{ResourceRecords: map[string]*models.ResourceRecord{"one": {Type: models.A}}}},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		name := "testing"
		z := &models.Zone{}
		call := mockClient.EXPECT().ValidateZone(context.Background(), &proto.ValidateZoneRequest{Name: name, Zone: grpc.ZoneToProtoBuf(z)})
		if tc.err != nil {
			call.Return(nil, tc.err)
		} else {
			call.Return(emptyMessage, nil)
		}

		err := grpcClient.ValidateZone(name, z)
		handleError(t, err, tc.err)

		if reflect.DeepEqual(z, tc.wantedZone) {
			t.Errorf("incorrect result: '%s', wanted: '%s'", z, tc.wantedZone)
		}
	}
}

func TestRender_Client(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err     error
		content string
	}{
		{content: "testing-content"},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		identifier := "testing"
		rr := &models.ResourceRecord{}
		call := mockClient.EXPECT().Render(context.Background(), &proto.RenderRequest{Identifier: identifier, ResourceRecord: grpc.ResourceRecordToProtoBuf(rr)})
		if tc.err != nil {
			call.Return(nil, tc.err)
		} else {
			call.Return(&proto.RenderResponse{Content: tc.content}, nil)
		}

		resp, err := grpcClient.Render(identifier, rr)
		handleError(t, err, tc.err)

		if resp != tc.content {
			t.Errorf("incorrect result: '%s', wanted: '%s'", resp, tc.content)
		}
	}
}
