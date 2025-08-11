package grpc

import (
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/bcurnow/zonemgr/schema"
)

// Updates the passed in zone with values from the protocol buff
func UpdateZoneFromProtoBuf(p *proto.Zone, zone *schema.Zone) {
	UpdateConfigFromProtoBuf(p.Config, zone.Config)
	UpdateResourceRecordsFromProtoBuf(p.ResourceRecords, zone.ResourceRecords)
	UpdateTTLFromProtoBuf(p.Ttl, zone.TTL)
}

func ZoneFromProtoBuf(p *proto.Zone) *schema.Zone {
	zone := &schema.Zone{}
	UpdateZoneFromProtoBuf(p, zone)
	return zone
}

func ZoneToProtoBuf(z *schema.Zone) *proto.Zone {
	return &proto.Zone{
		Config:          ConfigToProtoBufTo(z.Config),
		ResourceRecords: ResourceRecordsToProtoBuf(z.ResourceRecords),
		Ttl:             TTLToProtoBuf(z.TTL),
	}
}

func UpdateResourceRecordsFromProtoBuf(p map[string]*proto.ResourceRecord, rrs map[string]*schema.ResourceRecord) {
	for identifier, rr := range p {
		UpdateResourceRecordFromProtoBuf(rr, rrs[identifier])
	}
}

func ResoureRecordsFromProtoBuf(p map[string]*proto.ResourceRecord) map[string]*schema.ResourceRecord {
	rrs := make(map[string]*schema.ResourceRecord, len(p))
	for identifier, prr := range p {
		rr := &schema.ResourceRecord{}
		UpdateResourceRecordFromProtoBuf(prr, rr)
		rrs[identifier] = rr
	}

	return rrs
}

func ResourceRecordsToProtoBuf(rrs map[string]*schema.ResourceRecord) map[string]*proto.ResourceRecord {
	protoRRS := make(map[string]*proto.ResourceRecord, len(rrs))
	for identifier, rr := range rrs {
		protoRRS[identifier] = ResourceRecordToProtoBuf(rr)
	}
	return protoRRS
}
