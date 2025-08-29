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
	"fmt"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/models/testingutils"
)

func TestNormalizeZones(t *testing.T) {
	dnsSetup(t)
	defer testingutils.Teardown(t)

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(testZone.Config).Times(2)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config).Times(2)

	// Each plugin should have normalize called for each zone
	mockAPlugin.EXPECT().Normalize("record1", testZone.ResourceRecords["record1"]).Times(2)
	mockCNAMEPlugin.EXPECT().Normalize("record2", testZone.ResourceRecords["record2"]).Times(2)

	// Each plugin should have validate called for each zone
	mockAPlugin.EXPECT().ValidateZone("zone 1", testZone)
	mockCNAMEPlugin.EXPECT().ValidateZone("zone 1", testZone)
	mockAPlugin.EXPECT().ValidateZone("zone 2", testZone)
	mockCNAMEPlugin.EXPECT().ValidateZone("zone 2", testZone)

	if err := PluginNormalizer(mockPlugins, mockMetadata).Normalize(testZones); err != nil {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}

func TestNormalizeZones_NoZones(t *testing.T) {
	dnsSetup(t)
	defer testingutils.Teardown(t)

	if err := PluginNormalizer(mockPlugins, mockMetadata).Normalize(map[string]*models.Zone{}); err != nil {
		if err.Error() != "no zones found" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}

func TestNormalizeZones_NoPluginForRecordType(t *testing.T) {
	dnsSetup(t)
	defer testingutils.Teardown(t)

	invalidZone := &models.Zone{
		Config: testZone.Config,
		ResourceRecords: map[string]*models.ResourceRecord{
			"bad type": {Type: "bogus"},
		},
		TTL: testZone.TTL,
	}

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(invalidZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(invalidZone.Config)

	if err := PluginNormalizer(mockPlugins, mockMetadata).Normalize(map[string]*models.Zone{"invalid zone": invalidZone}); err != nil {
		if err.Error() != "unable to normalize zone 'invalid zone', no plugin for resource record type 'bogus', identifier: 'bad type'" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}
func TestNormalizeZones_NormalizeError(t *testing.T) {
	dnsSetup(t)
	defer testingutils.Teardown(t)

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)

	mockAPlugin.EXPECT().Normalize("record1", testZone.ResourceRecords["record1"]).Return(fmt.Errorf("test normalize error"))

	if err := PluginNormalizer(mockPlugins, mockMetadata).Normalize(testZones); err != nil {
		if err.Error() != "test normalize error" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}

func TestNormalizeZones_ValidateError(t *testing.T) {
	dnsSetup(t)
	defer testingutils.Teardown(t)

	// Each plugin should be configured once for each zone
	mockAPlugin.EXPECT().Configure(testZone.Config)
	mockCNAMEPlugin.EXPECT().Configure(testZone.Config)

	// Each plugin should have normalize called for each zone
	mockAPlugin.EXPECT().Normalize("record1", testZone.ResourceRecords["record1"])
	mockCNAMEPlugin.EXPECT().Normalize("record2", testZone.ResourceRecords["record2"])

	// Each plugin should have validate called for each zone, we're not sure which order they will iterate in so EXPECT for both
	mockAPlugin.EXPECT().ValidateZone("zone 1", testZone).Return(fmt.Errorf("test validate error")).MaxTimes(1)
	mockCNAMEPlugin.EXPECT().ValidateZone("zone 1", testZone).Return(fmt.Errorf("test validate error")).MaxTimes(1)

	if err := PluginNormalizer(mockPlugins, mockMetadata).Normalize(testZones); err != nil {
		if err.Error() != "test validate error" {
			t.Errorf("Error NormalizingZones: %s", err)
		}
	} else {
		t.Errorf("Error NormalizingZones: %s", err)
	}
}
