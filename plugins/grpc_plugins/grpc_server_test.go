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

func TestPluginVersion_Server(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)

	testCases := []struct {
		wantVersion string
		err         error
	}{
		{wantVersion: "testing-ver"},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		call := mockImpl.EXPECT().PluginVersion()
		if tc.err != nil {
			call.Return("", tc.err)
		} else {
			call.Return(tc.wantVersion, nil)
		}

		resp, err := grpcServer.PluginVersion(context.Background(), emptyMessage)
		handleError(t, err, tc.err)
		if tc.err == nil {
			if resp.Version != tc.wantVersion {
				t.Errorf("incorrect version: '%s', want: '%s'", err, tc.err)
			}
		}
	}
}

func TestPluginTypes_Server(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)

	testCases := []struct {
		wantTypes []plugins.PluginType
		err       error
	}{
		{wantTypes: plugins.PluginTypes(plugins.NS, plugins.A)},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		stringTypes := make([]string, len(tc.wantTypes))
		for i, pluginType := range tc.wantTypes {
			stringTypes[i] = string(pluginType)
		}

		call := mockImpl.EXPECT().PluginTypes()
		if tc.err != nil {
			call.Return(nil, tc.err)
		} else {
			call.Return(tc.wantTypes, nil)
		}

		resp, err := grpcServer.PluginTypes(context.Background(), emptyMessage)
		handleError(t, err, tc.err)
		if tc.err == nil {
			if !slices.Equal(resp.SupportedTypes, stringTypes) {
				t.Errorf("incorrect type: '%s', want: '%s'", err, tc.err)
			}
		}
	}
}

func TestConfigure_Server(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)

	testCases := []struct {
		err error
	}{
		{},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		c := &models.Config{}
		call := mockImpl.EXPECT().Configure(c)
		if tc.err != nil {
			call.Return(tc.err)
		} else {
			call.Return(nil)
		}

		req := &proto.ConfigureRequest{Config: grpc.ConfigToProtoBufTo(c)}
		resp, err := grpcServer.Configure(context.Background(), req)
		handleError(t, err, tc.err)
		if tc.err == nil {
			if !reflect.DeepEqual(resp, emptyMessage) {
				t.Errorf("incorrect response: '%v', want: '%v'", resp, emptyMessage)
			}
		}
	}
}

func TestNormalize_Server(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err error
	}{
		{},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		identifier := "testing"
		rr := &models.ResourceRecord{Type: models.A, Name: "Testing"}
		call := mockImpl.EXPECT().Normalize(identifier, rr)
		if tc.err != nil {
			call.Return(tc.err)
		} else {
			call.Return(nil)
		}

		req := &proto.NormalizeRequest{Identifier: identifier, ResourceRecord: grpc.ResourceRecordToProtoBuf(rr)}
		resp, err := grpcServer.Normalize(context.Background(), req)
		handleError(t, err, tc.err)
		if tc.err == nil {
			wantResp := &proto.NormalizeResponse{ResourceRecord: grpc.ResourceRecordToProtoBuf(rr)}
			if !reflect.DeepEqual(resp, wantResp) {
				t.Errorf("incorrect response: '%v', want: '%v'", resp, wantResp)
			}
		}
	}
}

func TestValidateZone_Server(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err error
	}{
		{},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		name := "testing"
		zone := &models.Zone{}
		call := mockImpl.EXPECT().ValidateZone(name, zone)
		if tc.err != nil {
			call.Return(tc.err)
		} else {
			call.Return(nil)
		}

		req := &proto.ValidateZoneRequest{Name: name, Zone: grpc.ZoneToProtoBuf(zone)}
		resp, err := grpcServer.ValidateZone(context.Background(), req)
		handleError(t, err, tc.err)
		if tc.err == nil {
			if !reflect.DeepEqual(resp, emptyMessage) {
				t.Errorf("incorrect response: '%v', want: '%v'", resp, emptyMessage)
			}
		}
	}
}
func TestRender_Server(t *testing.T) {
	setup_grpc(t)
	defer teardown_grpc(t)
	testCases := []struct {
		err error
	}{
		{},
		{err: errors.New("testing-err")},
	}

	for _, tc := range testCases {
		identifier := "testing"
		rr := &models.ResourceRecord{Type: models.A, Name: "Testing"}
		call := mockImpl.EXPECT().Render(identifier, rr)
		if tc.err != nil {
			call.Return("", tc.err)
		} else {
			call.Return("rendered", nil)
		}

		req := &proto.RenderRequest{Identifier: identifier, ResourceRecord: grpc.ResourceRecordToProtoBuf(rr)}
		resp, err := grpcServer.Render(context.Background(), req)
		handleError(t, err, tc.err)
		if tc.err == nil {
			wantResp := &proto.RenderResponse{Content: "rendered"}
			if !reflect.DeepEqual(resp, wantResp) {
				t.Errorf("incorrect response: '%v', want: '%v'", resp, wantResp)
			}
		}
	}
}
