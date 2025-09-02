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

package models

import "sort"

func WithSortedZones(zones map[string]*Zone, fn func(name string, zone *Zone) error) error {
	for _, zoneName := range sortedZoneKeys(zones) {
		if err := fn(zoneName, zones[zoneName]); err != nil {
			return err
		}
	}
	return nil
}

func sortedZoneKeys(zones map[string]*Zone) []string {
	keys := make([]string, 0, len(zones))
	for k := range zones {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
