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
	"reflect"
	"testing"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

func TestGRPCServer(t *testing.T) {
	originalFunc := registerZonemgrPluginServer
	defer func() { registerZonemgrPluginServer = originalFunc }()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockPlugin := plugins.NewMockZoneMgrPlugin(mockController)
	version := "testing-ver"
	mockPlugin.EXPECT().PluginVersion().Return(version, nil)

	registerZonemgrPluginServer = func(s grpc.ServiceRegistrar, srv proto.ZonemgrPluginServer) {
		if !reflect.ValueOf(s).IsNil() {
			t.Errorf("incorrect call: s=%v, want: s=nil", s)
		}

		// Call a method to make sure this is the mockPlugin
		svrVersion, err := srv.PluginVersion(context.Background(), &proto.Empty{})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		} else {
			if svrVersion.Version != version {
				t.Errorf("incorrect result: '%s', want: '%s'", svrVersion, version)
			}
		}
	}

	if err := (&GRPCPlugin{Impl: mockPlugin}).GRPCServer(nil, nil); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestGRPCClient(t *testing.T) {
	originalFunc := newZonemgrPluginClient
	defer func() { newZonemgrPluginClient = originalFunc }()

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mockClient := proto.NewMockZonemgrPluginClient(mockController)

	newZonemgrPluginClient = func(cc grpc.ClientConnInterface) proto.ZonemgrPluginClient {
		if !reflect.ValueOf(cc).IsNil() {
			t.Errorf("incorrect call: cc: %v, want: cc:%v", cc, nil)
		}
		return mockClient
	}

	client, err := (&GRPCPlugin{}).GRPCClient(context.Background(), nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		grpcClient := client.(*GRPCClient)
		if grpcClient.client != mockClient {
			t.Errorf("incorrect result: client: %v, want: client: %v", grpcClient.client, mockClient)
		}
	}
}
