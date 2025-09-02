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

	"github.com/bcurnow/zonemgr/internal/mocks"
	"github.com/bcurnow/zonemgr/models"
	"github.com/golang/mock/gomock"
)

var (
	mockNormalizer *mocks.MockNormalizer
)

func testZoneYamlParserSetup(t *testing.T) {
	dnsSetup(t)
	mockNormalizer = mocks.NewMockNormalizer(mockController)
}

func TestParse(t *testing.T) {
	testZoneYamlParserSetup(t)
	defer dnsTeardown(t)

	testCases := []struct {
		count     int
		inputFile string
		err       string
	}{
		{0, "empty.zones.yaml", "no zones found in input file"},
		{0, "only-zone-name.zones.yaml", "invalid input file only-zone-name.zones.yaml, no zone information for zone and_now_for_something_completely_unexpected"},
		{0, "invalid.zones.yaml", "failed to parse input YAML: yaml: unmarshal errors:\n  line 21: cannot unmarshal !!str `hello!` into bool"},
		{1, "minimal.zones.yaml", ""},
		{5, "multiple.zones.yaml", ""},
		{5, "missing.zones.yaml", "failed to open 'missing.zones.yaml': open missing.zones.yaml: no such file or directory"},
	}

	mockNormalizer.EXPECT().Normalize(gomock.Any(), globalConfig).MaxTimes(len(testCases))

	for _, tc := range testCases {
		zones, err := YamlZoneParser(mockNormalizer).Parse(tc.inputFile, &models.Config{})

		if err != nil {
			if tc.err == "" {
				t.Errorf("%s, unexpected error: %s", tc.inputFile, err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("%s, incorrect error: %s, want %s", tc.inputFile, err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("%s, expected error '%s', found none", tc.inputFile, tc.err)
			}
		}

		if tc.err == "" {
			if len(zones) != tc.count {
				t.Errorf("%s, zone count=%d, want %d", tc.inputFile, len(zones), 1)
			}
		}
	}
}

func TestParse_NormalizerError(t *testing.T) {
	testZoneYamlParserSetup(t)
	defer dnsTeardown(t)

	mockNormalizer.EXPECT().Normalize(gomock.Any(), globalConfig).Return(fmt.Errorf("testing normalizer error"))

	_, err := YamlZoneParser(mockNormalizer).Parse("minimal.zones.yaml", &models.Config{})
	if err == nil {
		t.Errorf("expected error")
	} else {
		want := "failed to normalize zones: testing normalizer error"
		if err.Error() != want {
			t.Errorf("incorrect error:%s, want %s", err, want)
		}
	}

}
