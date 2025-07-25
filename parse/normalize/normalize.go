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

// This function reviews various types of records to ensure that there are no errors that the YAML parser can't catch
// For example, not ending a value with a '.', using an IP vs a name in a CNAME record, etc.
func Normalize(zones map[string]*schema.Zone) error {
	for name, zone := range zones {
		// Check the zone name
		err := isFullyQualified(name)
		if err != nil {
			return fmt.Errorf("Invalid zone name %s: %w", name, err)
		}

		// Check the remaining fields in the zone
		err = NormalizeZone(zone)
		if err != nil {
			return fmt.Errorf("Failed to normalize zone %s: %w", name, err)
		}

		zones[name] = zone
	}

	return nil
}
