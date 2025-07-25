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

func normalizeNSRecord(nameserver string, record *schema.ResourceRecord) error {
	// Check the nameserver
	err := isFullyQualified(nameserver)
	if err != nil {
		return fmt.Errorf("Invalid nameserver %s: %w", nameserver, err)
	}

	// Check if the value is a valid name (not an IP address)
	if net.ParseIP(nameserver) != nil {
		return fmt.Errorf("NS record value cannot be an IP address: %s", nameserver)
	}

	// Check if the value is valid, in this case, the value" is the domain because this is optional in an NS record
	if record.Value != "" && !isValidNameOrWildcard(record.Value) {
		return fmt.Errorf("Invalid NS record value: %s", record.Value)
	}

	// Check if the class is valid
	if !isValidClass(record.Class) {
		return fmt.Errorf("Invalid NS record class: %s", record.Class)
	}

	return nil
}
