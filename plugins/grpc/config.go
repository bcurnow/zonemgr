package grpc

import (
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/bcurnow/zonemgr/schema"
)

// Updates the passed in config with values from the protocol buff
func UpdateConfigFromProtoBuf(p *proto.Config, c *schema.Config) {
	c.GenerateSerial = p.GenerateSerial
	c.SerialChangeIndex = p.SerialChangeIndex
	c.GenerateReverseLookupZones = p.GenerateReverseLookupZones
	c.PluginsDirectory = p.PluginsDirectory
	c.SerialChangeIndexDirectory = p.PluginsDirectory
}

func ConfigFromProtoBuf(p *proto.Config) *schema.Config {
	config := &schema.Config{}
	UpdateConfigFromProtoBuf(p, config)
	return config
}

func ConfigToProtoBufTo(c *schema.Config) *proto.Config {
	return &proto.Config{
		GenerateSerial:             c.GenerateSerial,
		SerialChangeIndex:          c.SerialChangeIndex,
		GenerateReverseLookupZones: c.GenerateReverseLookupZones,
		PluginsDirectory:           c.PluginsDirectory,
		SerialChangeIndexDirectory: c.SerialChangeIndexDirectory,
	}
}
