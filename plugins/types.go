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
	"github.com/hashicorp/go-plugin"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// This is the goplugin.Plugin implementation
type Plugin struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl TypeHandler
}

func (p *Plugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	proto.RegisterZonemgrPluginServer(server, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *Plugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewZonemgrPluginClient(client)}, nil
}

// Validate that we're correctly implenting goplugin.GRPCPlugin
var _ goplugin.GRPCPlugin = &Plugin{}

// Represents an instance of a plugin along with a bit of metadata
type TypeHandlerPlugin struct {
	IsBuiltIn  bool
	PluginName string
	PluginCmd  string
	Plugin     TypeHandler
}
