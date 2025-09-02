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
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/google/go-cmp/cmp"
)

func TestNormalize(t *testing.T) {
	dnsSetup(t)
	defer dnsTeardown(t)

	testCases := []struct {
		name                     string
		expectedConfig           *models.Config
		zones                    map[string]*models.Zone
		absPathErr               bool
		missingPluginErr         bool
		validateZoneErr          bool
		normalizeErr             bool
		missingPluginMetadataErr bool
	}{
		{name: "no-zones", zones: make(map[string]*models.Zone)},
		{name: "no-config-defaulting", expectedConfig: testZone.Config, zones: testZones},
		{name: "config-defaulting", expectedConfig: globalConfig, zones: map[string]*models.Zone{"nil-config-zone": {Config: nil}}},
		{name: "abs-path-error", expectedConfig: testZone.Config, zones: testZones, absPathErr: true},
		{name: "no-plugin-for-resource-record-type", expectedConfig: testZone.Config, zones: testZones, missingPluginErr: true},
		{name: "missing-plugin-metadata", expectedConfig: testZone.Config, zones: testZones, missingPluginMetadataErr: true},
		{name: "validation-error", expectedConfig: testZone.Config, zones: testZones, validateZoneErr: true},
		{name: "normalize-error", expectedConfig: testZone.Config, zones: testZones, normalizeErr: true},
	}

	dnsSetup(t)
	defer dnsTeardown(t)

	for _, tc := range testCases {
		// If we don't have any zones then we won't make any calls
		if len(tc.zones) != 0 {
			iterationZones := tc.zones
			if tc.absPathErr || tc.normalizeErr || tc.validateZoneErr || tc.missingPluginErr || tc.missingPluginMetadataErr {
				// We need to adjust the number of zones we'll iterate over because we won't get past the first one
				var firstZoneName string
				var firstZone *models.Zone
				for k, v := range tc.zones {
					firstZoneName = k
					firstZone = v
					break
				}
				iterationZones = map[string]*models.Zone{firstZoneName: firstZone}
			}

			models.WithSortedZones(iterationZones, func(zoneName string, zone *models.Zone) error {
				if tc.absPathErr {
					// This is the first thing we'll do so we only need this mock call
					mockFs.EXPECT().ToAbsoluteFilePath(tc.expectedConfig.SerialChangeIndexDirectory).Return("", errors.New("abs-path-testing"))
					return nil
				} else {
					mockFs.EXPECT().ToAbsoluteFilePath(tc.expectedConfig.SerialChangeIndexDirectory).Return(tc.expectedConfig.SerialChangeIndexDirectory, nil)
				}

				// If we don't have any plugins, Configure won't get called
				if !tc.missingPluginErr && !tc.missingPluginMetadataErr {
					// Each plugin should be configured once for each zone
					mockAPlugin.EXPECT().Configure(tc.expectedConfig)
					mockCNAMEPlugin.EXPECT().Configure(tc.expectedConfig)

					for identifier, rr := range zone.ResourceRecords {
						if rr.Type == models.A {
							if tc.normalizeErr {
								mockAPlugin.EXPECT().Normalize(identifier, rr).Return(errors.New("normalize-error"))
								// We won't call any other plugins
								break
							} else {
								mockAPlugin.EXPECT().Normalize(identifier, rr)
							}
						} else {
							if !tc.normalizeErr {
								mockCNAMEPlugin.EXPECT().Normalize(identifier, rr)
							}
						}

					}

					if tc.validateZoneErr {
						mockAPlugin.EXPECT().ValidateZone(zoneName, zone).Return(errors.New("validate-zone-error"))
					} else {
						if !tc.normalizeErr {
							// Each plugin should have validate called for each zone
							mockAPlugin.EXPECT().ValidateZone(zoneName, zone)
							mockCNAMEPlugin.EXPECT().ValidateZone(zoneName, zone)
						}
					}
				}

				return nil
			})
		}

		thePlugins := mockPlugins
		if tc.missingPluginErr {
			thePlugins = make(map[plugins.PluginType]plugins.ZoneMgrPlugin)
		}

		theMetadata := mockMetadata
		if tc.missingPluginMetadataErr {
			theMetadata = make(map[plugins.PluginType]*plugins.PluginMetadata)
		}

		if err := PluginNormalizer(thePlugins, theMetadata).Normalize(tc.zones, tc.expectedConfig); err != nil {
			// Determine which error we want
			want := ""
			if tc.absPathErr {
				want = "abs-path-testing"
			} else if len(tc.zones) == 0 {
				want = "no zones found"
			} else if tc.missingPluginErr {
				want = "unable to normalize zone 'zone1', no plugin for resource record type 'A', identifier: 'record1'"
			} else if tc.validateZoneErr {
				want = "validate-zone-error"
			} else if tc.normalizeErr {
				want = "normalize-error"
			} else if tc.missingPluginMetadataErr {
				want = "could not find plugin metadata for plugin type: A"
			}

			if want != "" {
				if err.Error() != want {
					t.Errorf("incorrect error: '%s', want '%s'", err, want)
				}
			} else {
				t.Errorf("%s - unexpected error:'%s'", tc.name, err)
			}
		} else {
			if tc.absPathErr || len(tc.zones) == 0 || tc.missingPluginErr {
				t.Errorf("%s - expected an error and found none", tc.name)
			} else {
				for zoneName, zone := range tc.zones {
					if diff := cmp.Diff(zone.Config, tc.expectedConfig); diff != "" {
						t.Errorf("%s - incorrect config for zone '%s':\n%s", tc.name, zoneName, diff)
					}
				}
			}
		}
	}
}
