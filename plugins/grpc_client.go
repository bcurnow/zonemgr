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

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/grpc"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

type GRPCClient struct {
	client proto.ZonemgrPluginClient
}

func (c *GRPCClient) PluginVersion() (string, error) {
	resp, err := c.client.PluginVersion(context.Background(), &proto.Empty{})
	if err != nil {
		return "", err
	}
	return resp.Version, nil
}

func (c *GRPCClient) PluginTypes() ([]PluginType, error) {
	resp, err := c.client.PluginTypes(context.Background(), &proto.Empty{})

	if err != nil {
		return nil, err
	}
	supportedPluginTypes := make([]PluginType, len(resp.SupportedTypes))
	for i, pluginTypeString := range resp.SupportedTypes {
		supportedPluginTypes[i] = PluginType(pluginTypeString)
	}

	return supportedPluginTypes, nil
}

func (c *GRPCClient) Configure(config *models.Config) error {
	_, err := c.client.Configure(context.Background(), &proto.ConfigureRequest{Config: grpc.ConfigToProtoBufTo(config)})
	if err != nil {
		return err
	}
	return nil
}

func (c *GRPCClient) Normalize(identifier string, rr *models.ResourceRecord) error {
	resp, err := c.client.Normalize(context.Background(), &proto.NormalizeRequest{
		Identifier:     identifier,
		ResourceRecord: grpc.ResourceRecordToProtoBuf(rr),
	})

	if err != nil {
		return err
	}
	// The expectation is that the Normalize method will make modifications to the ResourceRecord
	// Take the values from the respons and apply them to the original ResourceRecord pointer we received
	grpc.UpdateResourceRecordFromProtoBuf(resp.ResourceRecord, rr)
	return nil
}

func (c *GRPCClient) ValidateZone(name string, zone *models.Zone) error {
	_, err := c.client.ValidateZone(context.Background(), &proto.ValidateZoneRequest{
		Name: name,
		Zone: grpc.ZoneToProtoBuf(zone),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *GRPCClient) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	resp, err := c.client.Render(context.Background(), &proto.RenderRequest{
		Identifier:     identifier,
		ResourceRecord: grpc.ResourceRecordToProtoBuf(rr),
	})
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

var _ ZoneMgrPlugin = &GRPCClient{}
