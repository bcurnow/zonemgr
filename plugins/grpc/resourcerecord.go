package grpc

import (
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/bcurnow/zonemgr/schema"
)

// Updates the passed in ResourceRecord with values from the protocol buff
func UpdateResourceRecordFromProtoBuf(p *proto.ResourceRecord, rr *schema.ResourceRecord) {
	//Update the resource record with the new values
	rr.Name = p.Name
	rr.Type = p.Type
	rr.Class = p.Class
	rr.TTL = p.Ttl
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
	ret := &proto.ResourceRecord{
		Name:    rr.Name,
		Type:    rr.Type,
		Class:   rr.Class,
		Ttl:     rr.TTL,
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
