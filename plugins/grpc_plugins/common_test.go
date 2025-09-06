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
	"testing"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/golang/mock/gomock"
)

var (
	emptyMessage   = &proto.Empty{}
	grpcServer     *GRPCServer
	grpcClient     *GRPCClient
	mockController *gomock.Controller
	mockImpl       *plugins.MockZoneMgrPlugin
	mockClient     *proto.MockZonemgrPluginClient
)

func setup_grpc(t *testing.T) {
	mockController = gomock.NewController(t)
	mockImpl = plugins.NewMockZoneMgrPlugin(mockController)
	grpcServer = &GRPCServer{Impl: mockImpl}
	mockClient = proto.NewMockZonemgrPluginClient(mockController)
	grpcClient = &GRPCClient{client: mockClient}
}

func teardown_grpc(_ *testing.T) {
	mockController.Finish()
}

func handleError(t *testing.T, err error, wanted error) {
	if wanted != nil && err == nil {
		t.Errorf("expected error")
	} else if wanted != nil && err != nil {
		if err.Error() != wanted.Error() {
			t.Errorf("unexpected error:\n'%s'\nwant\n'%s'", err, wanted)
		}
	} else if wanted == nil && err != nil {
		t.Errorf("unexpected error:\n'%s'", err)
	}
}
