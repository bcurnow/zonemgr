package grpc

import (
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

func ConfigFromProtoBuf(p *proto.Config, c *models.Config) {
	if nil == p || nil == c {
		return
	}

	c.GenerateSerial = p.GenerateSerial
	c.GenerateReverseLookupZones = p.GenerateReverseLookupZones
	c.SerialChangeIndexDirectory = p.SerialChangeIndexDirectory
}

func ConfigToProtoBuf(c *models.Config) *proto.Config {
	if nil == c {
		return &proto.Config{}
	}

	return &proto.Config{
		GenerateSerial:             c.GenerateSerial,
		GenerateReverseLookupZones: c.GenerateReverseLookupZones,
		SerialChangeIndexDirectory: c.SerialChangeIndexDirectory,
	}
}
