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

	"github.com/bcurnow/zonemgr/logging"
	"github.com/bcurnow/zonemgr/schema"
)

var logger = logging.Logger().Named("zonefile")

func ToZoneFiles(zones map[string]*schema.Zone, outputDir string) error {
	for name, zone := range zones {
		if err := generateZone(name, zone, outputDir); err != nil {
			return err
		}

		if zone.Config.GenerateReverseLookupZones {
			logger.Debug("Zone has generate reverse lookup zones turned on", "zone", name)
			err := generateReverseLookupZones(zone, outputDir)
			if err != nil {
				return fmt.Errorf("Unable to generate reverse lookup zones for zone %s: %w\n", name, err)
			}
		}
	}
	return nil
}
