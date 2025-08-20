package grpc

import (
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

func UpdateTTLFromProtoBuf(p *proto.TTL, ttl *models.TTL) {
	ttl.Value = p.Ttl
	ttl.Comment = p.Comment
}

func TTLFromProtoBuf(p *proto.TTL) *models.TTL {
	ttl := &models.TTL{}
	UpdateTTLFromProtoBuf(p, ttl)
	return ttl
}

func TTLToProtoBuf(ttl *models.TTL) *proto.TTL {
	return &proto.TTL{Ttl: ttl.Value, Comment: ttl.Comment}
}
