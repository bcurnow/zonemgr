/**
 * Copyright (C) 2025 Brian Curnow
 *
 * This file is part of zonemgr.
 *
 * zonemgr is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * zonemgr is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with zonemgr.  If not, see <https://www.gnu.org/licenses/>.
 */

package dns

import (
	"fmt"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
)

var validations = plugins.V()

type ZoneReverser interface {
	ReverseZone(sourceZoneName string, zone *models.Zone) (map[string]*models.Zone, error)
}

type zoneReverser struct {
	ZoneReverser
}

func Reverser() ZoneReverser {
	return &zoneReverser{}
}

func (zr *zoneReverser) ReverseZone(sourceZoneName string, zone *models.Zone) (map[string]*models.Zone, error) {
	reverseLookupZones := make(map[string]*models.Zone)

	for _, rr := range zone.ResourceRecords {
		// We only care about A and AAAA records as they're the ones we're trying to reverse
		if rr.Type == models.A || rr.Type == models.AAAA {
			ip, err := utils.ParseIP(rr.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid IP address in %s record %q: %w", rr.Type, rr.Name, err)
			}
			zoneName := ip.ReverseZoneName()
			reverseZone, ok := reverseLookupZones[zoneName]
			if !ok {
				reverseZone = &models.Zone{
					Config:          zone.Config,
					ResourceRecords: make(map[string]*models.ResourceRecord),
					TTL:             zone.TTL,
				}

				// Add the SOA record for the zone
				sourceSOA := zone.SOARecord()
				reverseZone.ResourceRecords[zoneName] = &models.ResourceRecord{
					// Copy the values from the SOA record in the source zone
					Name:    zoneName,
					Type:    models.SOA,
					Class:   sourceSOA.Class,
					TTL:     sourceSOA.TTL,
					Values:  copyResourceRecordValues(sourceSOA.Values),
					Value:   sourceSOA.Value,
					Comment: sourceSOA.Comment,
				}
				reverseLookupZones[zoneName] = reverseZone
			}

			ptr := zr.toPTR(sourceZoneName, ip, rr)
			reverseZone.ResourceRecords[ptr.Name] = ptr
		}
	}

	return reverseLookupZones, nil
}

// copyResourceRecordValues returns a deep copy of values so a reverse zone's SOA record doesn't share
// backing storage with the source zone's SOA record; normalizing one would otherwise mutate the other.
func copyResourceRecordValues(values []*models.ResourceRecordValue) []*models.ResourceRecordValue {
	if values == nil {
		return nil
	}
	copied := make([]*models.ResourceRecordValue, len(values))
	for i, v := range values {
		valueCopy := *v
		copied[i] = &valueCopy
	}
	return copied
}

func (zr *zoneReverser) toPTR(sourceZoneName string, ip utils.IP, rr *models.ResourceRecord) *models.ResourceRecord {
	// The PTR record must be fully qualified
	ptrName := rr.Name
	if err := validations.EnsureFullyQualified("generated record", ptrName, rr.Type); err != nil {
		ptrName = validations.EnsureTrailingDot(ptrName + "." + sourceZoneName)
	}

	return &models.ResourceRecord{
		Name:   ip.PTRRecordValue(),
		Type:   models.PTR,
		Class:  rr.Class,
		TTL:    rr.TTL,
		Values: []*models.ResourceRecordValue{},
		// Each value must be fully qualified
		Value:   ptrName,
		Comment: rr.Comment,
	}
}
