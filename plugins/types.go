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

package plugins

import (
	"context"

	"github.com/bcurnow/zonemgr/plugins/proto"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// This is the goplugin.Plugin implementation
type GRPCPlugin struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl ZoneMgrPlugin
}

func (p *GRPCPlugin) GRPCServer(broker *goplugin.GRPCBroker, server *grpc.Server) error {
	proto.RegisterZonemgrPluginServer(server, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewZonemgrPluginClient(client)}, nil
}

// Validate that we're correctly implenting goplugin.GRPCPlugin
var _ goplugin.GRPCPlugin = &GRPCPlugin{}

// Represents an instance of a plugin along with a bit of metadata
type Plugin struct {
	IsBuiltIn  bool
	PluginName string
	PluginCmd  string
	Plugin     ZoneMgrPlugin
}
