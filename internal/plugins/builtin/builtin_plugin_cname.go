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
	"fmt"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
)

var _ plugins.ZoneMgrPlugin = &BuiltinPluginCNAME{}

type BuiltinPluginCNAME struct {
	plugins.ZoneMgrPlugin
}

func (p *BuiltinPluginCNAME) PluginVersion() (string, error) {
	return utils.Version(), nil
}

func (p *BuiltinPluginCNAME) PluginTypes() ([]plugins.Type, error) {
	return plugins.PluginTypes(plugins.CNAME), nil
}

func (p *BuiltinPluginCNAME) Configure(config *models.Config) error {
	// no config
	return nil
}

func (p *BuiltinPluginCNAME) Normalize(identifier string, rr *models.ResourceRecord) error {
	if err := validations.CommonValidations(identifier, rr, plugins.CNAME); err != nil {
		return err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := validations.EnsureValidNameOrWildcard(identifier, rr.Name, rr.Type); err != nil {
		return err
	}

	// Make sure the value isn't an IP
	if err := validations.EnsureNotIP(identifier, rr.RetrieveSingleValue(), rr.Type); err != nil {
		return err
	}

	return nil
}

// ValidateZone enforces the CNAME rules from RFC 1034 section 3.6.2, RFC 1912 section 2.4 and RFC 2181 section 10.1:
//   - a name that owns a CNAME must not own any other record type (this also forbids a CNAME at the zone apex, which always has SOA and NS)
//   - a name must own at most one CNAME
//   - CNAME chains are legal but must not loop
//
// The target is deliberately not required to exist in this zone or to have any particular record type; it may be an A, AAAA,
// another CNAME, or a name in another zone entirely.
func (p *BuiltinPluginCNAME) ValidateZone(name string, zone *models.Zone) error {
	resourceRecordsByType := zone.ResourceRecordsByType()
	cnameRecords := resourceRecordsByType[models.CNAME]
	if len(cnameRecords) == 0 {
		return nil
	}

	otherOwnerTypes := make(map[string]models.ResourceRecordType)
	for recordType, records := range resourceRecordsByType {
		if recordType == models.CNAME {
			continue
		}
		for _, resourceRecord := range records {
			otherOwnerTypes[canonicalName(name, resourceRecord.Name)] = recordType
		}
	}

	cnameIdentifiers := sortedIdentifiers(cnameRecords)
	cnameOwnerIdentifiers := make(map[string]string, len(cnameRecords))
	cnameTargets := make(map[string]string, len(cnameRecords))
	for _, identifier := range cnameIdentifiers {
		cnameRecord := cnameRecords[identifier]
		ownerName := canonicalName(name, cnameRecord.Name)

		if existingIdentifier, duplicate := cnameOwnerIdentifiers[ownerName]; duplicate {
			return fmt.Errorf("invalid CNAME record, '%s' and '%s' both define a CNAME for the name '%s', a name can only have one CNAME, zone: '%s'", existingIdentifier, identifier, cnameRecord.Name, name)
		}
		cnameOwnerIdentifiers[ownerName] = identifier
		cnameTargets[ownerName] = canonicalName(name, cnameRecord.RetrieveSingleValue())
	}

	for _, identifier := range cnameIdentifiers {
		cnameRecord := cnameRecords[identifier]
		ownerName := canonicalName(name, cnameRecord.Name)

		if otherType, conflict := otherOwnerTypes[ownerName]; conflict {
			return fmt.Errorf("invalid CNAME record, '%s' has the name '%s' which also has a record of type '%s', a name with a CNAME cannot have any other record types, zone: '%s'", identifier, cnameRecord.Name, otherType, name)
		}

		if err := checkForCNAMELoop(name, identifier, ownerName, cnameTargets); err != nil {
			return err
		}
	}

	return nil
}

func checkForCNAMELoop(zoneName string, identifier string, startName string, cnameTargets map[string]string) error {
	visited := map[string]bool{startName: true}
	current := cnameTargets[startName]
	for {
		if visited[current] {
			return fmt.Errorf("invalid CNAME record, '%s' is part of a CNAME loop, zone: '%s'", identifier, zoneName)
		}
		visited[current] = true

		next, isCNAME := cnameTargets[current]
		if !isCNAME {
			return nil
		}
		current = next
	}
}

func (p *BuiltinPluginCNAME) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	if err := validations.EnsureSupportedPluginType(identifier, rr.Type, plugins.CNAME); err != nil {
		return "", err
	}
	return rr.RenderSingleValueResource(), nil
}

func init() {
	registerBuiltIn(plugins.CNAME, &BuiltinPluginCNAME{})
}
