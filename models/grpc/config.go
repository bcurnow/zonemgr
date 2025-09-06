package grpc

import (
	"fmt"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

func ConfigFromProtoBuf(p *proto.Config, c *models.Config) {
	fmt.Printf("p=%s, c=%s\n", p, c)
	if nil == p || nil == c {
		return
	}

	c.GenerateSerial = p.GenerateSerial
	c.GenerateReverseLookupZones = p.GenerateReverseLookupZones
	c.SerialChangeIndexDirectory = p.SerialChangeIndexDirectory
}

func ConfigToProtoBuf(c *models.Config) *proto.Config {
	if nil == c {
		return nil
	}

	return &proto.Config{
		GenerateSerial:             c.GenerateSerial,
		GenerateReverseLookupZones: c.GenerateReverseLookupZones,
		SerialChangeIndexDirectory: c.SerialChangeIndexDirectory,
	}
}
