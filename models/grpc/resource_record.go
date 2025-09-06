package grpc

import (
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

func ResourceRecordFromProtoBuf(p *proto.ResourceRecord, rr *models.ResourceRecord) {
	if p == nil || rr == nil {
		return
	}
	rr.Name = p.Name
	rr.Type = models.ResourceRecordType(p.Type)
	rr.Class = models.ResourceRecordClass(p.Class)
	rr.TTL = p.Ttl
	resourceRecordValuesFromProtoBuf(p, rr)
	rr.Value = p.Value
	rr.Comment = p.Comment
}

func ResourceRecordToProtoBuf(rr *models.ResourceRecord) *proto.ResourceRecord {
	if rr == nil {
		return nil
	}
	ret := &proto.ResourceRecord{
		Name:    rr.Name,
		Type:    string(rr.Type),
		Class:   string(rr.Class),
		Ttl:     rr.TTL,
		Value:   rr.Value,
		Values:  resourceRecordValuesToProtoBuf(rr.Values),
		Comment: rr.Comment,
	}

	return ret
}

func resourceRecordValuesFromProtoBuf(p *proto.ResourceRecord, rr *models.ResourceRecord) {
	if p == nil || rr == nil || p.Values == nil {
		return
	}

	// Reset the Values so we end up with only the onces from the proto
	rr.Values = make([]*models.ResourceRecordValue, len(p.Values))
	for i, value := range p.Values {
		rr.Values[i] = &models.ResourceRecordValue{}
		rr.Values[i].Value = value.Value
		rr.Values[i].Comment = value.Comment
	}
}

func resourceRecordValuesToProtoBuf(rrvs []*models.ResourceRecordValue) []*proto.ResourceRecordValue {
	protoValues := make([]*proto.ResourceRecordValue, len(rrvs))
	for i, rrv := range rrvs {
		protoValues[i] = &proto.ResourceRecordValue{Value: rrv.Value, Comment: rrv.Comment}
	}
	return protoValues
}
