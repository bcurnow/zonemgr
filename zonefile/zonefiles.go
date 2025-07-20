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

	"github.com/bcurnow/zonemgr/sourceyaml"
)

func ToZoneFiles(zones map[string]*sourceyaml.Zone, outputDir string) error {
	for name, zone := range zones {
		GenerateZone(name, zone, outputDir)

		if zone.GenerateReverseLookupZones {
			fmt.Printf("Zone %s has generate reverse lookup zones turned on...\n", name)
			err := GenerateReverseLookupZones(zone, outputDir)
			if err != nil {
				return fmt.Errorf("Unable to generate reverse lookup zones for zone %s: %s\n", name, err)
			}
		}
	}
	return nil
}
