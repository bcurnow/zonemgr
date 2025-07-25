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
package normalize

import (
	"fmt"

	"github.com/bcurnow/zonemgr/parse/schema"
)

// Performs all the normalization on a single zone
func NormalizeZone(zone *schema.Zone) error {
	// Check the SOA record fields
	if err := normalizeSOAData(zone); err != nil {
		return err
	}

	// Normalize the serial number
	if zone.GenerateSerial {
		index := zone.SerialChangeIndex
		if zone.SerialChangeIndex == 0 {
			index = uint32(1)
		}

		serial, err := generateSerial(index)
		if err != nil {
			return fmt.Errorf("Unable to generate serial: %w", err)
		}
		zone.Serial = serial
	} else {
		if zone.Serial == 0 {
			return fmt.Errorf("serial number must be provided when generate_serial is false")
		}
	}
	// Normalize the Resource Records
	for name, record := range zone.ResourceRecords {
		if err := normlizeResourceRecord(name, &record, zone.ResourceRecords); err != nil {
			return err
		}
	}

	return nil
}

func normalizeSOAData(zone *schema.Zone) error {
	if zone.Class == "" {
		zone.Class = "IN" // Default class is IN
	} else if !isValidClass(zone.Class) {
		return fmt.Errorf("Invalid SOA class: %s", zone.Class)
	}

	err := isFullyQualified(zone.Nameserver)
	if err != nil {
		return fmt.Errorf("Invalid nameserver in SOA: %s: %w", zone.Nameserver, err)
	}

	// Check if the administrator is valid
	// We require an email but the format requires some special formatting
	formattedEmail, err := formatEmail(zone.Administrator)
	if err != nil {
		return fmt.Errorf("Failed to format administrator email: %w", err)
	}
	zone.Administrator = formattedEmail

	// Check the values of the SOA fields
	if zone.Refresh.Value <= 0 || zone.Retry.Value <= 0 || zone.Expire.Value <= 0 || zone.Minimum.Value < 0 {
		return fmt.Errorf("Invalid SOA values for zone, Refresh: %d, Retry: %d, Expire: %d, Minimum: %d. Values must be greater than or equal to zero (0)",
			zone.Refresh.Value, zone.Retry.Value, zone.Expire.Value, zone.Minimum.Value)

	}

	return nil
}
