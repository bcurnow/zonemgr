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
	"github.com/bcurnow/zonemgr/schema"
)

type GRPCServer struct {
	Impl TypeHandler
}

func (s *GRPCServer) PluginVersion(ctx context.Context, req *proto.Empty) (*proto.PluginVersionResponse, error) {
	version, err := s.Impl.PluginVersion()
	if err != nil {
		return nil, err
	}
	return &proto.PluginVersionResponse{Version: version}, nil
}

func (s *GRPCServer) PluginTypes(ctx context.Context, req *proto.Empty) (*proto.PluginTypesResponse, error) {
	supportedPluginTypes, err := s.Impl.PluginTypes()
	if err != nil {
		return nil, err
	}
	supportedPluginTypesStrings := make([]string, len(supportedPluginTypes))
	for i, pluginType := range supportedPluginTypes {
		supportedPluginTypesStrings[i] = string(pluginType)
	}
	return &proto.PluginTypesResponse{SupportedTypes: supportedPluginTypesStrings}, nil
}

func (s *GRPCServer) Configure(ctx context.Context, req *proto.ConfigureRequest) (*proto.Empty, error) {
	err := s.Impl.Configure(schema.Config.FromProtoBuf(schema.Config{}, req.Config))
	return &proto.Empty{}, err
}

func (s *GRPCServer) Normalize(ctx context.Context, req *proto.NormalizeRequest) (*proto.NormalizeResponse, error) {
	rr, err := s.Impl.Normalize(req.Identifier, schema.ResourceRecord.FromProtoBuf(schema.ResourceRecord{}, req.ResourceRecord))
	return &proto.NormalizeResponse{ResourceRecord: rr.ToProtoBuf()}, err
}

func (s *GRPCServer) ValidateZone(ctx context.Context, req *proto.ValidateZoneRequest) (*proto.Empty, error) {
	err := s.Impl.ValidateZone(req.Name, schema.Zone.FromProtoBuf(schema.Zone{}, req.Zone))
	return &proto.Empty{}, err
}

func (s *GRPCServer) Render(ctx context.Context, req *proto.RenderRequest) (*proto.RenderResonse, error) {
	record, err := s.Impl.Render(req.Identifier, schema.ResourceRecord.FromProtoBuf(schema.ResourceRecord{}, req.ResourceRecord))
	return &proto.RenderResonse{Content: record}, err
}

var _ proto.ZonemgrPluginServer = &GRPCServer{}
