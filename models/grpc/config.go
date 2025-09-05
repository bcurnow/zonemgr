package grpc

import (
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

// Updates the passed in config with values from the protocol buff
func UpdateConfigFromProtoBuf(p *proto.Config, c *models.Config) {
	if c == nil {
		return
	}
	c.GenerateSerial = p.GenerateSerial
	c.GenerateReverseLookupZones = p.GenerateReverseLookupZones
	c.SerialChangeIndexDirectory = p.SerialChangeIndexDirectory
}

func ConfigFromProtoBuf(p *proto.Config) *models.Config {
	config := &models.Config{}
	UpdateConfigFromProtoBuf(p, config)
	return config
}

func ConfigToProtoBufTo(c *models.Config) *proto.Config {
	if nil == c {
		return &proto.Config{}
	}

	return &proto.Config{
		GenerateSerial:             c.GenerateSerial,
		GenerateReverseLookupZones: c.GenerateReverseLookupZones,
		SerialChangeIndexDirectory: c.SerialChangeIndexDirectory,
	}
}
