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
	ConfigFromProtoBuf(p.Config, z.Config)
	resourceRecordsFromProtoBuf(p, z)
	TTLFromProtoBuf(p.Ttl, z.TTL)
}

func ZoneToProtoBuf(z *models.Zone) *proto.Zone {
	if z == nil {
		return nil
	}
	proto :=
		&proto.Zone{
			Config:          ConfigToProtoBuf(z.Config),
			ResourceRecords: resourceRecordsToProtoBuf(z),
			Ttl:             TTLToProtoBuf(z.TTL),
		}
	return proto
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

func resourceRecordsToProtoBuf(z *models.Zone) map[string]*proto.ResourceRecord {
	if z.ResourceRecords == nil {
		rrs := make(map[string]*proto.ResourceRecord, 0)
		return rrs
	}
	rrs := make(map[string]*proto.ResourceRecord, len(z.ResourceRecords))

	for identifier, rr := range z.ResourceRecords {
		rrs[identifier] = ResourceRecordToProtoBuf(rr)
	}

	return rrs
}
