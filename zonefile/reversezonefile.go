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
package zonefile

import (
	"strings"

	"github.com/bcurnow/zonemgr/normalize"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
	"github.com/hashicorp/go-hclog"
)

// Generates a set of reverse zone files for the specific forward zone
func generateReverseLookupZones(name string, zone *schema.Zone, outputDir string) error {
	// Convert the full set of resource records for this zone into a set of reverse lookup zones
	hclog.L().Trace("Generating reverse lookup zones", "zoneName", name)
	reverseLookupZones := toReverseZones(name, zone)
	if err := normalize.NormalizeZones(reverseLookupZones); err != nil {
		return err
	}

	for name, zone := range reverseLookupZones {
		if err := writeZoneFile(name, zone, outputDir); err != nil {
			return err
		}
	}
	return nil
}

// Converts a Zone to a a set of reverse lookup zones
func toReverseZones(sourceZoneName string, zone *schema.Zone) map[string]*schema.Zone {
	reverseLookupZones := make(map[string]*schema.Zone)

	for _, rr := range zone.ResourceRecords {
		// We only care about A records as they're the ones we're trying to reverse
		// TODO should we also reverse CNAMEs?
		if rr.Type == schema.A {
			zoneName := reverseZoneName(rr.Value)
			reverseZone, ok := reverseLookupZones[zoneName]
			if !ok {
				reverseZone = &schema.Zone{
					Config:          zone.Config,
					ResourceRecords: make(map[string]*schema.ResourceRecord),
					TTL:             zone.TTL,
				}

				// Add the SOA record for the zone
				sourceSOA := zone.SOARecord()
				reverseZone.ResourceRecords[zoneName] = &schema.ResourceRecord{
					// Copy the values from the SOZ record in the source zone
					Name:    zoneName,
					Type:    schema.SOA,
					Class:   sourceSOA.Class,
					TTL:     sourceSOA.TTL,
					Values:  sourceSOA.Values,
					Value:   sourceSOA.Value,
					Comment: sourceSOA.Comment,
				}
				reverseLookupZones[zoneName] = reverseZone
			}

			ptr := toPTR(sourceZoneName, rr)
			reverseZone.ResourceRecords[ptr.Name] = ptr
		}
	}

	return reverseLookupZones
}

func toPTR(sourceZoneName string, rr *schema.ResourceRecord) *schema.ResourceRecord {
	ptrName := rr.Name
	if err := plugins.IsFullyQualified(ptrName, rr.Value, rr); err != nil {
		ptrName = plugins.EnsureTrailingDot(ptrName + "." + sourceZoneName)

	}

	return &schema.ResourceRecord{
		Name:   lastOctet(rr.Value), //An A records Name/identifier should be an IP, the name of the PTR record is just the last octet
		Type:   schema.PTR,
		Class:  rr.Class,
		TTL:    rr.TTL,
		Values: []*schema.ResourceRecordValue{},
		// Each value must be fully qualified
		Value:   ptrName,
		Comment: rr.Comment,
	}
}

func reverseZoneName(ip string) string {
	// Reverse zones are named based on the reverse of the first three octets of an IP
	// For example, if the IP is 10.2.2.10 the reverse zone name would be 2.2.10-in-addr.arpa
	// Get the last three octets
	octets := strings.Split(ip, ".")
	octets = octets[:len(octets)-1]

	// Reverse the octets
	for i, j := 0, len(octets)-1; i < j; i, j = i+1, j-1 {
		octets[i], octets[j] = octets[j], octets[i] // Swapping elements
	}

	// NOTE the zone must end with a dot (.) or it won't actually work, ORIGINs must be fully qualified!
	return strings.Join(octets, ".") + ".in-addr.arpa."
}

// Retrieves the last octet of an IPv4 address
// Given 10.2.2.76, this function would return 76
func lastOctet(ip string) string {
	octets := strings.Split(ip, ".")
	return octets[len(octets)-1]
}
