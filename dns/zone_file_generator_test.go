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
	"errors"
	"os"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/utils"
)

func TestPluginZoneFileGenerator(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	res1 := PluginZoneFileGenerator(mockPlugins, mockMetadata)
	res2 := PluginZoneFileGenerator(mockPlugins, mockMetadata)

	if res1 == res2 {
		t.Errorf("Expected a new instance on each call, got same instance")
	}
}

func TestGenerateZone(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)
	// We want to use the actual implementation for this test
	fs = utils.FS()
	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)
	mockAPlugin.EXPECT().Render("record1", &models.ResourceRecord{Type: models.A}).Return("record1", nil)
	mockCNAMEPlugin.EXPECT().Render("record2", &models.ResourceRecord{Type: models.CNAME}).Return("record2", nil)

	if err := PluginZoneFileGenerator(mockPlugins, mockMetadata).GenerateZone("testing", testZone, "."); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	defer os.Remove("./testing")
}

func TestGenerate_NoPlugin(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	g := &pluginZoneFileGenerator{plugins: mockPlugins, metadata: mockMetadata}

	// Insert an unknown resource record type
	// Replace the resourceRecords
	testZone.ResourceRecords = map[string]*models.ResourceRecord{"bogus": {Type: models.AAAA}}

	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)

	_, err := g.generate("testing", testZone)
	if err == nil {
		t.Errorf("expected error")
	}

	want := "unable to write zone 'testing', no plugin for resource record type 'AAAA', identifier: 'bogus'"

	if err.Error() != want {
		t.Errorf("Unexpected error: %s, want %s", err, want)
	}
}

func TestGenerate_RenderError(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	g := &pluginZoneFileGenerator{plugins: mockPlugins, metadata: mockMetadata}

	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)
	want := "testing error"
	mockAPlugin.EXPECT().Render("record1", &models.ResourceRecord{Type: models.A}).Return("", errors.New(want))

	_, err := g.generate("testing", testZone)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != want {
		t.Errorf("Unexpected error: %s, want %s", err, want)
	}
}

func TestGenerate_NoTTL(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	g := &pluginZoneFileGenerator{plugins: mockPlugins, metadata: mockMetadata}

	testZone.TTL = nil
	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)
	mockAPlugin.EXPECT().Render("record1", &models.ResourceRecord{Type: models.A}).Return("record1", nil)
	mockCNAMEPlugin.EXPECT().Render("record2", &models.ResourceRecord{Type: models.CNAME}).Return("record2", nil)

	content, err := g.generate("testing", testZone)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := "$ORIGIN testing\nrecord1\nrecord2\n"

	if string(content) != want {
		t.Errorf("Unexpected content:\n'%s'\nwant\n'%s'\n", string(content), want)
	}
}

func TestGenerate_WithTTL(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)
	g := &pluginZoneFileGenerator{plugins: mockPlugins, metadata: mockMetadata}

	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)
	mockAPlugin.EXPECT().Render("record1", &models.ResourceRecord{Type: models.A}).Return("record1", nil)
	mockCNAMEPlugin.EXPECT().Render("record2", &models.ResourceRecord{Type: models.CNAME}).Return("record2", nil)

	content, err := g.generate("testing", testZone)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := "$ORIGIN testing\n$TTL 30 ;testZone-TTL\nrecord1\nrecord2\n"

	if string(content) != want {
		t.Errorf("Unexpected content:\n'%s'\nwant\n'%s'\n", string(content), want)
	}
}
