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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bcurnow/zonemgr/sourceyaml"
)

type ReverseZoneTemplateData struct {
	Name            string
	ForwardZone     string // This is the name of the original zone (e.g. example.com.) that we're creating the reverse zone for
	Zone            *sourceyaml.Zone
	ResourceRecords map[string]sourceyaml.ResourceRecord
}

// Generates a set of reverse zone files for the specific forward zone
func GenerateReverseLookupZones(forwardZone string, zone *sourceyaml.Zone, outputDir string) error {
	// Convert the full set of resource records for this zone into a set of reverse lookup zones
	reverseLookupZones := toReverseZones(zone.ResourceRecords)

	// Write the various reverse lookup zone file
	err := toReverseZoneFiles(forwardZone, zone, reverseLookupZones, outputDir)
	if err != nil {
		return fmt.Errorf("Error generating reverse lookup zone files: %w", err)
	}
	return nil
}

// Takes a map of resource records and returns a set of reverse lookup zone names with only the necessary ResourceRecords ("A" records) included
func toReverseZones(resourceRecords map[string]sourceyaml.ResourceRecord) map[string]map[string]sourceyaml.ResourceRecord {
	reverseLookupZones := make(map[string]map[string]sourceyaml.ResourceRecord)
	for name, record := range resourceRecords {
		if record.Type == "A" {
			records := initReverseLookupZone(reverseZoneName(record.Value), reverseLookupZones)
			records[name] = record
		}
	}

	return reverseLookupZones
}

func initReverseLookupZone(name string, reverseLookupZones map[string]map[string]sourceyaml.ResourceRecord) map[string]sourceyaml.ResourceRecord {
	_, valid := reverseLookupZones[name]
	if !valid {
		reverseLookupRecords := make(map[string]sourceyaml.ResourceRecord)
		reverseLookupZones[name] = reverseLookupRecords
		return reverseLookupRecords
	}
	return reverseLookupZones[name]
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

func toReverseZoneFiles(forwardZone string, zone *sourceyaml.Zone, reverseZones map[string]map[string]sourceyaml.ResourceRecord, outputDir string) error {
	funcMap := template.FuncMap{
		"lastOctet": lastOctet,
	}

	template, err := template.New("reversezonefile.tmpl").Funcs(funcMap).Parse(reverseZoneFileTemplate)
	if err != nil {
		return fmt.Errorf("Failed to parse template: %w", err)
	}

	for name, resourceRecords := range reverseZones {
		outputFile, err := os.Create(filepath.Join(outputDir, name))
		if err != nil {
			return fmt.Errorf("Failed to create output file for reverse lookkup zone %s: %w", name, err)
		}
		defer outputFile.Close()

		fmt.Printf("Generating %s for reverse lookup zone %s\n", outputFile.Name(), name)
		err = template.Execute(outputFile, ReverseZoneTemplateData{Name: name, Zone: zone, ResourceRecords: resourceRecords, ForwardZone: forwardZone})
	}

	return nil
}

// Retrieves the last octet of an IPv4 address
// Given 10.2.2.76, this function would return 76
func lastOctet(ip string) string {
	octets := strings.Split(ip, ".")
	return octets[len(octets)-1]
}
