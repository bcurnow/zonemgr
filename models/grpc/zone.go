package grpc

import (
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

// Updates the passed in zone with values from the protocol buff
func ZoneFromProtoBuf(p *proto.Zone, z *models.Zone) {
	if p == nil || z == nil {
		return
	}
	// Make sure to default the Config and TTL if necessary otherwise we
	// can end up with a p that has a valid config but return a z that has it set to nil (because it came in nil)
	if z.Config == nil && p.Config != nil {
		z.Config = &models.Config{}
	}

	if z.TTL == nil && p.Ttl != nil {
		z.TTL = &models.TTL{}
	}

	ConfigFromProtoBuf(p.Config, z.Config)
	resourceRecordsFromProtoBuf(p, z)
	TTLFromProtoBuf(p.Ttl, z.TTL)
}

func ZoneToProtoBuf(z *models.Zone) *proto.Zone {
	if z == nil {
		return nil
	}
	return &proto.Zone{
		Config:          ConfigToProtoBuf(z.Config),
		ResourceRecords: resourceRecordsToProtoBuf(z.ResourceRecords),
		Ttl:             TTLToProtoBuf(z.TTL),
	}
}

func resourceRecordsFromProtoBuf(p *proto.Zone, z *models.Zone) {
	rrs := make(map[string]*models.ResourceRecord, len(p.ResourceRecords))
	z.ResourceRecords = rrs
	for identifier, prr := range p.ResourceRecords {
		rr := &models.ResourceRecord{}
		ResourceRecordFromProtoBuf(prr, rr)
		rrs[identifier] = rr
	}
}

func resourceRecordsToProtoBuf(rrs map[string]*models.ResourceRecord) map[string]*proto.ResourceRecord {
	protoRRs := make(map[string]*proto.ResourceRecord, len(rrs))

	for identifier, rr := range rrs {
		protoRRs[identifier] = ResourceRecordToProtoBuf(rr)
	}

	return protoRRs
}
