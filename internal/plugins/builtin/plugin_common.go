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
package builtin

import (
	"sort"
	"strings"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
)

var validations plugins.Validator = plugins.V()

func sortedIdentifiers(records map[string]*models.ResourceRecord) []string {
	identifiers := make([]string, 0, len(records))
	for identifier := range records {
		identifiers = append(identifiers, identifier)
	}
	sort.Strings(identifiers)
	return identifiers
}

// canonicalName produces a comparable form of a DNS name: lowercase (RFC 4343), with '@' and relative names
// expanded against the zone name, and always ending in a trailing dot.
func canonicalName(zoneName string, recordName string) string {
	zone := strings.ToLower(strings.TrimSuffix(zoneName, ".")) + "."
	recordName = strings.ToLower(recordName)

	switch {
	case recordName == "@" || recordName == "":
		return zone
	case strings.HasSuffix(recordName, "."):
		return recordName
	default:
		return recordName + "." + zone
	}
}
