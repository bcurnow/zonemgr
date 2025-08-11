package grpc

import (
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/bcurnow/zonemgr/schema"
)

// Updates the passed in ResourceRecord with values from the protocol buff
func UpdateResourceRecordFromProtoBuf(p *proto.ResourceRecord, rr *schema.ResourceRecord) {
	var ttl *int32 = nil
	if p.Ttl != -1 {
		ttl = &p.Ttl
	}

	//Update the resource record with the new values
	rr.Name = p.Name
	rr.Type = p.Type
	rr.Class = p.Class
	rr.TTL = ttl
	UpdateResourceRecordValuesFromProtoBuf(p.Values, rr.Values)
	rr.Value = p.Value
	rr.Comment = p.Comment
}

func ResourceRecordFromProtoBuf(p *proto.ResourceRecord) *schema.ResourceRecord {
	rr := &schema.ResourceRecord{}
	UpdateResourceRecordFromProtoBuf(p, rr)
	return rr
}

func ResourceRecordToProtoBuf(rr *schema.ResourceRecord) *proto.ResourceRecord {
	var ttl int32 = -1
	if rr.TTL != nil {
		// We're using a negative number so we can check for it the other way as well and set appropriately
		ttl = *rr.TTL
	}
	ret := &proto.ResourceRecord{
		Name:    rr.Name,
		Type:    rr.Type,
		Class:   rr.Class,
		Ttl:     ttl,
		Value:   rr.Value,
		Values:  ResourceRecordValuesToProtoBuf(rr.Values),
		Comment: rr.Comment,
	}

	return ret
}

func UpdateResourceRecordValuesFromProtoBuf(p []*proto.ResourceRecordValue, rrs []*schema.ResourceRecordValue) {
	for i, value := range p {
		UpdateResourceRecordValueFromProtoBuf(value, rrs[i])
	}
}

func UpdateResourceRecordValueFromProtoBuf(p *proto.ResourceRecordValue, rr *schema.ResourceRecordValue) {
	rr.Value = p.Value
	rr.Comment = p.Comment
}

func ResourceRecordValueFromProtoBuf(p *proto.ResourceRecordValue) *schema.ResourceRecordValue {
	rr := &schema.ResourceRecordValue{}
	UpdateResourceRecordValueFromProtoBuf(p, rr)
	return rr
}

func ResourceRecordValuesToProtoBuf(rrvs []*schema.ResourceRecordValue) []*proto.ResourceRecordValue {
	protoValues := make([]*proto.ResourceRecordValue, len(rrvs))
	for i, rrv := range rrvs {
		protoValues[i] = ResourceRecordValueToProtoBuf(rrv)
	}
	return protoValues
}

func ResourceRecordValueToProtoBuf(rrv *schema.ResourceRecordValue) *proto.ResourceRecordValue {
	return &proto.ResourceRecordValue{Value: rrv.Value, Comment: rrv.Comment}
}
