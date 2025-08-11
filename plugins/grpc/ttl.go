package grpc

import (
	"github.com/bcurnow/zonemgr/plugins/proto"
	"github.com/bcurnow/zonemgr/schema"
)

func UpdateTTLFromProtoBuf(p *proto.TTL, ttl *schema.TTL) {
	ttl.Value = &p.Ttl
	ttl.Comment = p.Comment
}

func TTLFromProtoBuf(p *proto.TTL) *schema.TTL {
	ttl := &schema.TTL{}
	UpdateTTLFromProtoBuf(p, ttl)
	return ttl
}

func TTLToProtoBuf(ttl *schema.TTL) *proto.TTL {
	return &proto.TTL{Ttl: *ttl.Value, Comment: ttl.Comment}
}
