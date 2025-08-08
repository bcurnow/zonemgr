/*
Copyright © 2025 Brian Curnow

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package schema

import (
	"strings"

	"github.com/bcurnow/zonemgr/plugins/proto"
)

// Represents the overall Zone file structure, the YAML file is an array of these
type Zone struct {
	TTL                   TTL                       `yaml:"ttl"`
	ResourceRecords       map[string]ResourceRecord `yaml:"resource_records"`
	Config                Config                    `yaml:"config"`
	resourceRecordsByType map[string]map[string]ResourceRecord
}

func (z *Zone) ResourceRecordsByType() map[string]map[string]ResourceRecord {
	if z.resourceRecordsByType == nil {
		z.resourceRecordsByType = make(map[string]map[string]ResourceRecord)

		for identifier, rr := range z.ResourceRecords {
			rrType := strings.ToUpper(rr.Type)
			_, ok := z.resourceRecordsByType[rrType]
			if !ok {
				z.resourceRecordsByType[rrType] = make(map[string]ResourceRecord)
			}
			z.resourceRecordsByType[rrType][identifier] = rr
		}
	}
	return z.resourceRecordsByType
}

func (z Zone) ToProtoBuf() *proto.Zone {
	return &proto.Zone{
		Config:          z.Config.ToProtoBuf(),
		ResourceRecords: z.toProtoBufResourceRecords(),
		Ttl:             z.TTL.ToProtoBuf(),
	}
}

func (z Zone) FromProtoBuf(p *proto.Zone) Zone {
	return Zone{
		Config:          Config.FromProtoBuf(Config{}, p.Config),
		ResourceRecords: z.fromProtoBufResourceRecords(p.ResourceRecords),
		TTL:             TTL.FromProtoBuf(TTL{}, p.Ttl),
	}
}

func (z Zone) toProtoBufResourceRecords() map[string]*proto.ResourceRecord {
	rrs := make(map[string]*proto.ResourceRecord, len(z.ResourceRecords))
	for identifier, rr := range z.ResourceRecords {
		rrs[identifier] = rr.ToProtoBuf()
	}
	return rrs
}

func (z Zone) fromProtoBufResourceRecords(p map[string]*proto.ResourceRecord) map[string]ResourceRecord {
	rrs := make(map[string]ResourceRecord, len(p))
	for identifier, rr := range p {
		rrs[identifier] = ResourceRecord.FromProtoBuf(ResourceRecord{}, rr)
	}
	return rrs
}

type TTL struct {
	Value   *int32 `yaml:"value"`
	Comment string `yaml:"comment"`
}

func (ttl TTL) ToProtoBuf() *proto.TTL {
	return &proto.TTL{Ttl: *ttl.Value, Comment: ttl.Comment}
}

func (ttl TTL) FromProtoBuf(p *proto.TTL) TTL {
	return TTL{Value: &p.Ttl, Comment: p.Comment}
}
