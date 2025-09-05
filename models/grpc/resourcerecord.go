package grpc

import (
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

// Updates the passed in ResourceRecord with values from the protocol buff
func UpdateResourceRecordFromProtoBuf(p *proto.ResourceRecord, rr *models.ResourceRecord) {
	//Update the resource record with the new values
	rr.Name = p.Name
	rr.Type = models.ResourceRecordType(p.Type)
	rr.Class = models.ResourceRecordClass(p.Class)
	rr.TTL = p.Ttl
	UpdateResourceRecordValuesFromProtoBuf(p.Values, rr.Values)
	rr.Value = p.Value
	rr.Comment = p.Comment
}

func ResourceRecordFromProtoBuf(p *proto.ResourceRecord) *models.ResourceRecord {
	rr := &models.ResourceRecord{}
	UpdateResourceRecordFromProtoBuf(p, rr)
	return rr
}

func ResourceRecordToProtoBuf(rr *models.ResourceRecord) *proto.ResourceRecord {
	ret := &proto.ResourceRecord{
		Name:    rr.Name,
		Type:    string(rr.Type),
		Class:   string(rr.Class),
		Ttl:     rr.TTL,
		Value:   rr.Value,
		Values:  ResourceRecordValuesToProtoBuf(rr.Values),
		Comment: rr.Comment,
	}

	return ret
}

func UpdateResourceRecordValuesFromProtoBuf(p []*proto.ResourceRecordValue, rrs []*models.ResourceRecordValue) {
	for i, value := range p {
		UpdateResourceRecordValueFromProtoBuf(value, rrs[i])
	}
}

func UpdateResourceRecordValueFromProtoBuf(p *proto.ResourceRecordValue, rr *models.ResourceRecordValue) {
	rr.Value = p.Value
	rr.Comment = p.Comment
}

func ResourceRecordValueFromProtoBuf(p *proto.ResourceRecordValue) *models.ResourceRecordValue {
	rr := &models.ResourceRecordValue{}
	UpdateResourceRecordValueFromProtoBuf(p, rr)
	return rr
}

func ResourceRecordValuesToProtoBuf(rrvs []*models.ResourceRecordValue) []*proto.ResourceRecordValue {
	protoValues := make([]*proto.ResourceRecordValue, len(rrvs))
	for i, rrv := range rrvs {
		protoValues[i] = ResourceRecordValueToProtoBuf(rrv)
	}
	return protoValues
}

func ResourceRecordValueToProtoBuf(rrv *models.ResourceRecordValue) *proto.ResourceRecordValue {
	return &proto.ResourceRecordValue{Value: rrv.Value, Comment: rrv.Comment}
}
