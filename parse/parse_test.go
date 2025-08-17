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
package parse

import (
	"fmt"
	"testing"

	"github.com/bcurnow/zonemgr/schema"
	"github.com/bcurnow/zonemgr/testing/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

var (
	mockController *gomock.Controller
	mockNormalizer *mocks.MockNormalizer
	defaultConfig  *schema.Config
)

func setup(t *testing.T) {
	mockController = gomock.NewController(t)
	mockNormalizer = mocks.NewMockNormalizer(mockController)
	normalizer = mockNormalizer
	defaultConfig = &schema.Config{}
	defaultConfig.ConfigDefaults()
}

func teardown(_ *testing.T) {
	mockController.Finish()
}

func TestParse(t *testing.T) {
	setup(t)
	defer teardown(t)

	testCases := []struct {
		count     int
		inputFile string
		err       string
	}{
		{0, "empty.zones.yaml", "failed to unmarshal from empty.zones.yaml: no zones found in input file"},
		{0, "only-zone-name.zones.yaml", "invalid input file only-zone-name.zones.yaml, no zone information for zone and_now_for_something_completely_unexpected"},
		{0, "invalid.zones.yaml", "failed to unmarshal from invalid.zones.yaml: failed to parse input YAML: yaml: unmarshal errors:\n  line 21: cannot unmarshal !!str `hello!` into bool"},
		{1, "minimal.zones.yaml", ""},
		{5, "multiple.zones.yaml", ""},
		{5, "missing.zones.yaml", "failed to open input missing.zones.yaml: open missing.zones.yaml: no such file or directory"},
	}

	mockNormalizer.EXPECT().Normalize(gomock.Any()).MaxTimes(len(testCases))

	for _, tc := range testCases {
		zones, err := Parser().Parse(tc.inputFile)

		if err != nil {
			if tc.err == "" {
				t.Errorf("%s, did not expect error: %s", tc.inputFile, err)
			} else {
				if err.Error() != tc.err {
					t.Errorf("%s, unexpected error:%s, want %s", tc.inputFile, err, tc.err)
				}
			}
		} else {
			if tc.err != "" {
				t.Errorf("%s, expected error %s", tc.inputFile, tc.err)
			}
		}

		if tc.err == "" {
			if len(zones) != tc.count {
				t.Errorf("%s, zone count=%d, want %d", tc.inputFile, len(zones), 1)
			}

			for _, zone := range zones {
				// Make sure the config defaulting works
				if diff := cmp.Diff(zone.Config, defaultConfig); diff != "" {
					t.Errorf("%s, incorrect config:\n%s", tc.inputFile, diff)
				}
			}
		}

	}
}

func TestParse_NormalizerError(t *testing.T) {
	setup(t)
	defer teardown(t)

	mockNormalizer.EXPECT().Normalize(gomock.Any()).Return(fmt.Errorf("testing normalizer error"))

	_, err := Parser().Parse("minimal.zones.yaml")
	if err == nil {
		t.Errorf("expected error")
	} else {
		want := "failed to normalize zones: testing normalizer error"
		if err.Error() != want {
			t.Errorf("incorrect error:%s, want %s", err, want)
		}
	}

}
