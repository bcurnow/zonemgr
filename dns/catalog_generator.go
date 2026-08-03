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

package dns

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
)

// catalogSchemaVersion is the RFC 9432 "version" TXT record value, currently fixed at 2.
const catalogSchemaVersion = "2"

type CatalogGenerator interface {
	// Additively injects the RFC 9432 catalog records (a "version" TXT record and one PTR record
	// per member zone) into zone.ResourceRecords. Fails if a member record's identifier collides
	// with a record the user already authored in the zone.
	AddCatalogRecords(catalogZoneName string, zone *models.Zone, memberZoneNames []string) error
}

type pluginCatalogGenerator struct {
	CatalogGenerator
	plugins  map[plugins.Type]plugins.ZoneMgrPlugin
	metadata map[plugins.Type]*plugins.Metadata
}

func PluginCatalogGenerator(pluginMap map[plugins.Type]plugins.ZoneMgrPlugin, metadata map[plugins.Type]*plugins.Metadata) CatalogGenerator {
	return &pluginCatalogGenerator{plugins: pluginMap, metadata: metadata}
}

func (cg *pluginCatalogGenerator) AddCatalogRecords(catalogZoneName string, zone *models.Zone, memberZoneNames []string) error {
	if err := plugins.WithSortedPlugins(cg.plugins, cg.metadata, func(pluginType plugins.Type, p plugins.ZoneMgrPlugin, metadata *plugins.Metadata) error {
		return p.Configure(zone.Config)
	}); err != nil {
		return err
	}

	// Snapshot the identifiers the user authored before we start injecting so we can tell a genuine
	// collision apart from re-encountering a record we ourselves already injected for a duplicate member.
	preExisting := make(map[string]struct{}, len(zone.ResourceRecords))
	for identifier := range zone.ResourceRecords {
		preExisting[identifier] = struct{}{}
	}

	// The catalog zone was already validated during parsing, so it is guaranteed to have exactly one SOA.
	class := zone.SOARecord().Class

	versionIdentifier := "version." + catalogZoneName
	versionRecord := &models.ResourceRecord{
		Name:  versionIdentifier,
		Type:  models.TXT,
		Class: class,
		Value: catalogSchemaVersion,
	}
	if err := cg.addRecord(zone, preExisting, versionIdentifier, versionRecord); err != nil {
		return err
	}

	for _, memberName := range sortedCatalogMembers(memberZoneNames) {
		identifier := catalogMemberLabel(memberName) + ".zones." + catalogZoneName
		ptrRecord := &models.ResourceRecord{
			Name:   identifier,
			Type:   models.PTR,
			Class:  class,
			Values: []*models.ResourceRecordValue{},
			Value:  memberName,
		}
		if err := cg.addRecord(zone, preExisting, identifier, ptrRecord); err != nil {
			return err
		}
	}

	return nil
}

// addRecord normalizes and inserts rr under identifier, unless identifier was already present in the
// zone before injection started (a genuine collision with a user-authored record) or has already been
// injected earlier in this call (a duplicate member name, which is a no-op).
func (cg *pluginCatalogGenerator) addRecord(zone *models.Zone, preExisting map[string]struct{}, identifier string, rr *models.ResourceRecord) error {
	if _, ok := preExisting[identifier]; ok {
		return fmt.Errorf("unable to generate catalog zone, identifier '%s' is already defined in the catalog zone", identifier)
	}
	if _, ok := zone.ResourceRecords[identifier]; ok {
		return nil
	}

	plugin := cg.plugins[plugins.Type(rr.Type)]
	if nil == plugin {
		return fmt.Errorf("unable to generate catalog zone, no plugin for resource record type '%s', identifier: '%s'", rr.Type, identifier)
	}
	if err := plugin.Normalize(identifier, rr); err != nil {
		return err
	}

	zone.ResourceRecords[identifier] = rr
	return nil
}

// sortedCatalogMembers lowercases and trailing-dots each member name, removes duplicates and returns
// the result in sorted order so catalog generation is deterministic regardless of input ordering.
func sortedCatalogMembers(memberZoneNames []string) []string {
	memberSet := make(map[string]struct{}, len(memberZoneNames))
	for _, name := range memberZoneNames {
		memberSet[validations.EnsureTrailingDot(strings.ToLower(name))] = struct{}{}
	}

	members := make([]string, 0, len(memberSet))
	for name := range memberSet {
		members = append(members, name)
	}
	sort.Strings(members)
	return members
}

// catalogMemberLabel returns the RFC 9432 unique label for a member zone: the SHA-1 hash, hex encoded,
// of the zone name in DNS wire format. This matches BIND's catz convention for stable labels.
func catalogMemberLabel(memberName string) string {
	hash := sha1.Sum(dnsWireFormat(memberName))
	return hex.EncodeToString(hash[:])
}

// dnsWireFormat converts a dotted, trailing-dot DNS name into RFC1035 wire format: length-prefixed
// labels terminated by the zero-length root label.
func dnsWireFormat(name string) []byte {
	labels := strings.Split(strings.TrimSuffix(name, "."), ".")
	wire := make([]byte, 0, len(name)+1)
	for _, label := range labels {
		wire = append(wire, byte(len(label)))
		wire = append(wire, []byte(label)...)
	}
	return append(wire, 0)
}
