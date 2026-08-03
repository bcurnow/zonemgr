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
	"testing"

	"github.com/bcurnow/zonemgr/internal/plugins/builtin"
	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
)

func TestCatalogMemberLabel(t *testing.T) {
	testCases := []struct {
		name       string
		memberName string
		want       string
	}{
		{name: "two labels", memberName: "example.com.", want: "c5e4b4da1e5a620ddaa3635e55c3732a5b49c7f4"},
		{name: "single label", memberName: "home.", want: "e7a4b5c9b47166b6a40cedce0df4432f2b41d57d"},
	}

	for _, tc := range testCases {
		got := catalogMemberLabel(tc.memberName)
		if got != tc.want {
			t.Errorf("unexpected label for %s: got '%s', want '%s'", tc.name, got, tc.want)
		}
	}
}

func TestAddCatalogRecords(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	zone := &models.Zone{
		Config: &models.Config{},
		ResourceRecords: map[string]*models.ResourceRecord{
			"catalog.example.com.": {Type: models.SOA, Name: "catalog.example.com.", Class: models.INTERNET, Value: "SOA"},
		},
	}

	catalogGenerator := PluginCatalogGenerator(catalogTestPlugins(), catalogTestMetadata())
	// Duplicate, case-varied member names must hash to the same label and only produce one PTR record.
	if err := catalogGenerator.AddCatalogRecords("catalog.example.com.", zone, []string{"example.com.", "Example.com."}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versionRecord, ok := zone.ResourceRecords["version.catalog.example.com."]
	if !ok {
		t.Fatal("expected a version TXT record to be added")
	}
	if versionRecord.Type != models.TXT || versionRecord.Value != catalogSchemaVersion {
		t.Errorf("unexpected version record: %+v", versionRecord)
	}

	ptrIdentifier := "c5e4b4da1e5a620ddaa3635e55c3732a5b49c7f4.zones.catalog.example.com."
	ptrRecord, ok := zone.ResourceRecords[ptrIdentifier]
	if !ok {
		t.Fatalf("expected a PTR record at '%s'", ptrIdentifier)
	}
	if ptrRecord.Type != models.PTR || ptrRecord.Value != "example.com." {
		t.Errorf("unexpected PTR record: %+v", ptrRecord)
	}

	// SOA + version TXT + a single PTR (the duplicate, case-varied member name must not produce two records)
	if len(zone.ResourceRecords) != 3 {
		t.Errorf("expected 3 resource records, got %d: %v", len(zone.ResourceRecords), zone.ResourceRecords)
	}
}

func TestAddCatalogRecords_EmptyMemberList(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	zone := &models.Zone{
		Config: &models.Config{},
		ResourceRecords: map[string]*models.ResourceRecord{
			"catalog.example.com.": {Type: models.SOA, Name: "catalog.example.com.", Value: "SOA"},
		},
	}

	catalogGenerator := PluginCatalogGenerator(catalogTestPlugins(), catalogTestMetadata())
	if err := catalogGenerator.AddCatalogRecords("catalog.example.com.", zone, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// SOA + version TXT only, no PTRs
	if len(zone.ResourceRecords) != 2 {
		t.Errorf("expected 2 resource records, got %d: %v", len(zone.ResourceRecords), zone.ResourceRecords)
	}
}

func TestAddCatalogRecords_CollisionWithUserAuthoredRecord(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	zone := &models.Zone{
		Config: &models.Config{},
		ResourceRecords: map[string]*models.ResourceRecord{
			"catalog.example.com.":         {Type: models.SOA, Name: "catalog.example.com.", Value: "SOA"},
			"version.catalog.example.com.": {Type: models.TXT, Name: "version.catalog.example.com.", Value: "user-authored"},
		},
	}

	catalogGenerator := PluginCatalogGenerator(catalogTestPlugins(), catalogTestMetadata())
	err := catalogGenerator.AddCatalogRecords("catalog.example.com.", zone, []string{})
	if err == nil {
		t.Fatal("expected an error for a colliding identifier, found none")
	}
}

func TestAddCatalogRecords_MemberMissingTrailingDot(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	zone := &models.Zone{
		Config: &models.Config{},
		ResourceRecords: map[string]*models.ResourceRecord{
			"catalog.example.com.": {Type: models.SOA, Name: "catalog.example.com.", Value: "SOA"},
		},
	}

	catalogGenerator := PluginCatalogGenerator(catalogTestPlugins(), catalogTestMetadata())
	if err := catalogGenerator.AddCatalogRecords("catalog.example.com.", zone, []string{"example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ptrIdentifier := "c5e4b4da1e5a620ddaa3635e55c3732a5b49c7f4.zones.catalog.example.com."
	ptrRecord, ok := zone.ResourceRecords[ptrIdentifier]
	if !ok {
		t.Fatal("expected the member name to be trailing-dotted before hashing")
	}
	if ptrRecord.Value != "example.com." {
		t.Errorf("expected the PTR value to carry a trailing dot, got '%s'", ptrRecord.Value)
	}
}

func TestAddCatalogRecords_SingleLabelMemberRejected(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	zone := &models.Zone{
		Config: &models.Config{},
		ResourceRecords: map[string]*models.ResourceRecord{
			"catalog.": {Type: models.SOA, Name: "catalog.", Value: "SOA"},
		},
	}

	catalogGenerator := PluginCatalogGenerator(catalogTestPlugins(), catalogTestMetadata())
	err := catalogGenerator.AddCatalogRecords("catalog.", zone, []string{"home."})
	if err == nil {
		t.Fatal("expected an error for a single-label member zone name, found none")
	}
}

func TestPluginCatalogGenerator(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	one := PluginCatalogGenerator(catalogTestPlugins(), catalogTestMetadata())
	two := PluginCatalogGenerator(catalogTestPlugins(), catalogTestMetadata())

	if one == two {
		t.Error("expected two different instances")
	}
}

// catalogTestPlugins wires the real builtin TXT/PTR plugins so injected catalog records are validated
// through the same code path user-authored records get.
func catalogTestPlugins() map[plugins.Type]plugins.ZoneMgrPlugin {
	return map[plugins.Type]plugins.ZoneMgrPlugin{
		plugins.TXT: &builtin.BuiltinPluginTXT{},
		plugins.PTR: &builtin.BuiltinPluginPTR{},
	}
}

func catalogTestMetadata() map[plugins.Type]*plugins.Metadata {
	return map[plugins.Type]*plugins.Metadata{
		plugins.TXT: {Name: string(plugins.TXT), Command: "Built In", BuiltIn: true},
		plugins.PTR: {Name: string(plugins.PTR), Command: "Built In", BuiltIn: true},
	}
}
