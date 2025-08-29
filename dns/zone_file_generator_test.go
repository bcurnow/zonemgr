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
	"path/filepath"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/models/testingutils"
)

var outputDir = "zone_file_generator_test"

func TestPluginZoneFileGenerator(t *testing.T) {
	dnsSetup(t)
	defer testingutils.Teardown(t)

	res1 := PluginZoneFileGenerator(mockPlugins)
	res2 := PluginZoneFileGenerator(mockPlugins)

	if res1 == res2 {
		t.Errorf("Expected a new instance on each call, got same instance")
	}
}

func TestGenerateZone_InvalidOutputDir(t *testing.T) {
	err := PluginZoneFileGenerator(mockPlugins).GenerateZone("testing", testZone, "bogus")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestGenerateZone_NoPlugin(t *testing.T) {
	dnsSetup(t)
	fileSetup(t)
	defer fileTeardown(t)
	g := PluginZoneFileGenerator(mockPlugins)

	// Insert an unknown resource record type
	// Replace the resourceRecords
	testZone.ResourceRecords = map[string]*models.ResourceRecord{"bogus": &models.ResourceRecord{Type: models.AAAA}}

	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)

	err := g.GenerateZone("testing", testZone, outputDir)
	if err == nil {
		t.Errorf("expected error")
	}

	want := "unable to write zone 'testing', no plugin for resource record type 'AAAA', identifier: 'bogus'"

	if err.Error() != want {
		t.Errorf("Unexpected error: %s, want %s", err, want)
	}
}

func TestGenerateZone_RenderError(t *testing.T) {
	dnsSetup(t)
	fileSetup(t)
	defer fileTeardown(t)
	g := PluginZoneFileGenerator(mockPlugins)

	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)
	want := "testing error"
	mockAPlugin.EXPECT().Render("record1", &models.ResourceRecord{Type: models.A}).Return("", errors.New(want))

	err := g.GenerateZone("testing", testZone, outputDir)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != want {
		t.Errorf("Unexpected error: %s, want %s", err, want)
	}
}

func TestGenerateZone_NoTTL(t *testing.T) {
	dnsSetup(t)
	fileSetup(t)
	defer fileTeardown(t)
	g := PluginZoneFileGenerator(mockPlugins)

	testZone.TTL = nil
	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)
	mockAPlugin.EXPECT().Render("record1", &models.ResourceRecord{Type: models.A}).Return("record1", nil)
	mockCNAMEPlugin.EXPECT().Render("record2", &models.ResourceRecord{Type: models.CNAME}).Return("record2", nil)

	err := g.GenerateZone("testing", testZone, outputDir)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := "$ORIGIN testing\nrecord1\nrecord2\n"

	outputBytes, err := os.ReadFile(filepath.Join(outputDir, "testing"))
	if err != nil {
		t.Errorf("Unable to read the output file: %s", err)
	}
	output := string(outputBytes)
	if output != want {
		t.Errorf("Unexpected content:\n'%s'\nwant\n'%s'\n", output, want)
	}
}

func TestGenerateZone_WithTTL(t *testing.T) {
	dnsSetup(t)
	fileSetup(t)
	defer fileTeardown(t)
	g := PluginZoneFileGenerator(mockPlugins)

	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)
	mockAPlugin.EXPECT().Render("record1", &models.ResourceRecord{Type: models.A}).Return("record1", nil)
	mockCNAMEPlugin.EXPECT().Render("record2", &models.ResourceRecord{Type: models.CNAME}).Return("record2", nil)

	err := g.GenerateZone("testing", testZone, outputDir)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := "$ORIGIN testing\n$TTL 30 ; testZone-TTL\nrecord1\nrecord2\n"

	outputBytes, err := os.ReadFile(filepath.Join(outputDir, "testing"))
	if err != nil {
		t.Errorf("Unable to read the output file: %s", err)
	}
	output := string(outputBytes)
	if output != want {
		t.Errorf("Unexpected content:\n'%s'\nwant\n'%s'\n", output, want)
	}
}
func fileSetup(t *testing.T) {
	// Remove any previous testing directory
	err := os.RemoveAll(outputDir)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Recreate the empty directory
	err = os.Mkdir(outputDir, 0755)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}
}

func fileTeardown(t *testing.T) {
	err := os.RemoveAll(outputDir)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
