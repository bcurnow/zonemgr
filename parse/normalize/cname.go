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
	"net"

	"github.com/bcurnow/zonemgr/parse/schema"
)

func normalizeCNAMERecord(name string, record *schema.ResourceRecord, resourceRecords map[string]schema.ResourceRecord) error {
	// Check if the name is valid
	if !isValidName(name) {
		return fmt.Errorf("Invalid CNAME record name: %s", name)
	}

	// Check if the value is a valid name (not an IP address)
	if net.ParseIP(record.Value) != nil {
		return fmt.Errorf("CNAME record value cannot be an IP address: %s", record.Value)
	}

	// CNAME values must also point to a valid name within the zone
	found, valid := resourceRecords[record.Value]
	if !valid || found.Type != "A" {
		return fmt.Errorf("CNAME record value: %s must point to a valid A record within the zone", record.Value)
	}

	return nil
}
