package grpc

import (
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins/proto"
)

func TTLFromProtoBuf(p *proto.TTL, ttl *models.TTL) {
	if p == nil || ttl == nil {
		return
	}
	ttl.Value = p.Ttl
	ttl.Comment = p.Comment
}

func TTLToProtoBuf(ttl *models.TTL) *proto.TTL {
	if ttl == nil {
		return nil
	}

	return &proto.TTL{Ttl: ttl.Value, Comment: ttl.Comment}
}
