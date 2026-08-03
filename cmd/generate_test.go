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

package cmd

import (
	"errors"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/spf13/viper"
)

func TestPreRunE_Generate(t *testing.T) {
	setup(t)
	defer teardown(t)
	testCases := []struct {
		absErrOutput               bool
		absErrInput                bool
		generateReverseLookupZones bool
		generateSerial             bool
	}{
		{},
		{generateReverseLookupZones: true},
		{generateSerial: true},
		{generateReverseLookupZones: true, generateSerial: true},
		{absErrOutput: true},
		{absErrInput: true},
	}

	for _, tc := range testCases {
		v = viper.New()
		v.BindPFlags(generateCmd.Flags())
		call := mockFs.EXPECT().ToAbsoluteFilePath("testing-dir")
		if tc.absErrOutput {
			call.Return("", errors.New("absErrOutput"))
		} else {
			call.Return("testing-dir", nil)
			call = mockFs.EXPECT().ToAbsoluteFilePath("testing")
			if tc.absErrInput {
				call.Return("", errors.New("absErrInput"))
			} else {
				call.Return("testing", nil)
				mockPluginManager.EXPECT().Plugins().Return(testPlugins).Times(3)
				mockPluginManager.EXPECT().Metadata().Return(testMetadata).Times(3)
			}
		}

		args := []string{"--input-file", "testing", "--output-dir", "testing-dir", "--serial-change-index-directory", "testing-scid"}
		if tc.generateReverseLookupZones {
			args = append(args, "--generate-reverse-lookup-zones")
		}
		if tc.generateSerial {
			args = append(args, "--generate-serial")
		}

		generateCmd.ParseFlags(args)
		if err := generateCmd.PreRunE(generateCmd, args); err != nil {
			want := ""
			if tc.absErrOutput {
				want = "absErrOutput"
			} else if tc.absErrInput {
				want = "absErrInput"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s", err, want)
			}
		} else {
			if tc.absErrOutput || tc.absErrInput {
				t.Error("expected an error, found none")
			}

			want := "testing-dir"
			if outputDir != want {
				t.Errorf("incorrect output dir: '%s', want: '%s'", outputDir, want)
			}

			want = "testing"
			if inputFile != want {
				t.Errorf("incorrect input file dir: '%s', want: '%s'", inputFile, want)
			}

			if zoneFileGenerator == mockZoneFileGenerator {
				t.Errorf("expected zoneFileGenerator to not be a mock")
			}

			if normalizer == mockNormalizer {
				t.Errorf("expected normalizer to not be a mock")
			}

			if parser == mockParser {
				t.Errorf("expected parser to not be a mock")
			}

			if catalogGenerator == mockCatalogGenerator {
				t.Errorf("expected catalogGenerator to not be a mock")
			}
		}
	}
}

func TestRunE_Generate(t *testing.T) {
	setup(t)
	defer teardown(t)

	testCases := []struct {
		parseErr                    bool
		zoneFileGeneratorErr        bool
		reverseZoneErr              bool
		normalizerErr               bool
		reverseZoneFileGeneratorErr bool
	}{
		{},
		{parseErr: true},
		{zoneFileGeneratorErr: true},
		{reverseZoneErr: true},
		{normalizerErr: true},
		{reverseZoneFileGeneratorErr: true},
	}

	inputFile = "testing"
	outputDir = "testing-dir"
	for _, tc := range testCases {
		call := mockParser.EXPECT().Parse(inputFile)
		if tc.parseErr {
			call.Return(nil, errors.New("parseErr"))
		} else {
			zones := make(map[string]*models.Zone)
			zoneOne := &models.Zone{Config: &models.Config{}}
			zoneTwo := &models.Zone{Config: &models.Config{GenerateReverseLookupZones: true}}
			zones["one"] = zoneOne
			zones["two"] = zoneTwo
			call.Return(zones, nil)

			// Pass 1 (compute) runs to completion, for every zone, before any file is written in pass 2 -
			// so ReverseZone/Normalize are always attempted regardless of how pass 2 (GenerateZone) will fare.
			reverseZones := make(map[string]*models.Zone)
			reverseZoneOne := &models.Zone{}
			reverseZoneTwo := &models.Zone{}
			reverseZones["reverse-one"] = reverseZoneOne
			reverseZones["reverse-two"] = reverseZoneTwo
			call = mockZoneReverser.EXPECT().ReverseZone("two", zoneTwo)
			if tc.reverseZoneErr {
				call.Return(nil, errors.New("reverseZoneErr"))
			} else {
				call.Return(reverseZones, nil)

				call = mockNormalizer.EXPECT().Normalize(reverseZones)
				if tc.normalizerErr {
					call.Return(errors.New("normalizerErr"))
				} else {
					call.Return(nil)

					call = mockZoneFileGenerator.EXPECT().GenerateZone("one", zoneOne, outputDir)
					if tc.zoneFileGeneratorErr {
						call.Return(errors.New("zoneFileGeneratorErr"))
					} else {
						call.Return(nil)
						mockZoneFileGenerator.EXPECT().GenerateZone("two", zoneTwo, outputDir).Return(nil)

						call = mockZoneFileGenerator.EXPECT().GenerateZone("reverse-one", reverseZoneOne, outputDir)
						if tc.reverseZoneFileGeneratorErr {
							call.Return(errors.New("reverseZoneFileGeneratorErr"))
						} else {
							call.Return(nil)
							mockZoneFileGenerator.EXPECT().GenerateZone("reverse-two", reverseZoneTwo, outputDir).Return(nil)
						}
					}
				}
			}
		}

		if err := generateCmd.RunE(generateCmd, []string{}); err != nil {
			want := ""
			if tc.parseErr {
				want = "failed to parse input file testing: parseErr"
			} else if tc.zoneFileGeneratorErr {
				want = "zoneFileGeneratorErr"
			} else if tc.reverseZoneErr {
				want = "reverseZoneErr"
			} else if tc.normalizerErr {
				want = "normalizerErr"
			} else if tc.reverseZoneFileGeneratorErr {
				want = "reverseZoneFileGeneratorErr"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.parseErr || tc.zoneFileGeneratorErr || tc.reverseZoneErr || tc.normalizerErr || tc.reverseZoneFileGeneratorErr {
				t.Error("expected an error, found none")
			}
		}
	}
}

func TestRunE_Generate_CatalogZone(t *testing.T) {
	setup(t)
	defer teardown(t)

	inputFile = "testing"
	outputDir = "testing-dir"

	zoneOne := &models.Zone{Config: &models.Config{}}
	catalogZone := &models.Zone{Config: &models.Config{IsCatalog: true}}
	zones := map[string]*models.Zone{
		"one":                  zoneOne,
		"catalog.example.com.": catalogZone,
	}

	mockParser.EXPECT().Parse(inputFile).Return(zones, nil)
	mockCatalogGenerator.EXPECT().AddCatalogRecords("catalog.example.com.", catalogZone, []string{"one"}).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("one", zoneOne, outputDir).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("catalog.example.com.", catalogZone, outputDir).Return(nil)

	if err := generateCmd.RunE(generateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunE_Generate_CatalogZone_IncludeReverseZones(t *testing.T) {
	setup(t)
	defer teardown(t)

	inputFile = "testing"
	outputDir = "testing-dir"

	zoneOne := &models.Zone{Config: &models.Config{GenerateReverseLookupZones: true}}
	catalogZone := &models.Zone{Config: &models.Config{IsCatalog: true, CatalogIncludeReverseZones: true}}
	zones := map[string]*models.Zone{
		"one":                  zoneOne,
		"catalog.example.com.": catalogZone,
	}
	reverseZones := map[string]*models.Zone{
		"reverse.arpa.": {},
	}

	mockParser.EXPECT().Parse(inputFile).Return(zones, nil)
	mockZoneReverser.EXPECT().ReverseZone("one", zoneOne).Return(reverseZones, nil)
	mockNormalizer.EXPECT().Normalize(reverseZones).Return(nil)
	mockCatalogGenerator.EXPECT().AddCatalogRecords("catalog.example.com.", catalogZone, []string{"one", "reverse.arpa."}).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("one", zoneOne, outputDir).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("reverse.arpa.", reverseZones["reverse.arpa."], outputDir).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("catalog.example.com.", catalogZone, outputDir).Return(nil)

	if err := generateCmd.RunE(generateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunE_Generate_CatalogGeneratorErr(t *testing.T) {
	setup(t)
	defer teardown(t)

	inputFile = "testing"
	outputDir = "testing-dir"

	zoneOne := &models.Zone{Config: &models.Config{}}
	catalogZone := &models.Zone{Config: &models.Config{IsCatalog: true}}
	zones := map[string]*models.Zone{
		"one":                  zoneOne,
		"catalog.example.com.": catalogZone,
	}

	mockParser.EXPECT().Parse(inputFile).Return(zones, nil)
	mockCatalogGenerator.EXPECT().AddCatalogRecords("catalog.example.com.", catalogZone, []string{"one"}).Return(errors.New("catalogGeneratorErr"))
	// No GenerateZone calls are expected: nothing should be written until every zone, including
	// catalog zones, has been fully computed.

	err := generateCmd.RunE(generateCmd, []string{})
	if err == nil {
		t.Fatal("expected an error, found none")
	}
	if err.Error() != "catalogGeneratorErr" {
		t.Errorf("incorrect error: '%s', want: 'catalogGeneratorErr'", err)
	}
}

func TestRunE_Generate_MergesReverseZonesFromMultipleSourceZones(t *testing.T) {
	setup(t)
	defer teardown(t)

	inputFile = "testing"
	outputDir = "testing-dir"

	zoneOne := &models.Zone{Config: &models.Config{GenerateReverseLookupZones: true}}
	zoneTwo := &models.Zone{Config: &models.Config{GenerateReverseLookupZones: true}}
	zones := map[string]*models.Zone{
		"one": zoneOne,
		"two": zoneTwo,
	}

	sharedZoneFromOne := &models.Zone{
		ResourceRecords: map[string]*models.ResourceRecord{
			"shared.arpa.": {Type: models.SOA, Name: "shared.arpa.", Value: "SOA-from-one"},
			"1":            {Type: models.PTR, Name: "1", Value: "host-one.example.com."},
		},
	}
	sharedZoneFromTwo := &models.Zone{
		ResourceRecords: map[string]*models.ResourceRecord{
			"shared.arpa.": {Type: models.SOA, Name: "shared.arpa.", Value: "SOA-from-two"},
			"2":            {Type: models.PTR, Name: "2", Value: "host-two.example.com."},
		},
	}

	mockParser.EXPECT().Parse(inputFile).Return(zones, nil)
	mockZoneReverser.EXPECT().ReverseZone("one", zoneOne).Return(map[string]*models.Zone{"shared.arpa.": sharedZoneFromOne}, nil)
	mockZoneReverser.EXPECT().ReverseZone("two", zoneTwo).Return(map[string]*models.Zone{"shared.arpa.": sharedZoneFromTwo}, nil)
	mockNormalizer.EXPECT().Normalize(map[string]*models.Zone{"shared.arpa.": sharedZoneFromOne}).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("one", zoneOne, outputDir).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("two", zoneTwo, outputDir).Return(nil)
	mockZoneFileGenerator.EXPECT().GenerateZone("shared.arpa.", sharedZoneFromOne, outputDir).Return(nil)

	if err := generateCmd.RunE(generateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "one" sorts first, so it establishes the SOA; "two"'s PTR is merged in alongside "one"'s.
	if sharedZoneFromOne.ResourceRecords["shared.arpa."].Value != "SOA-from-one" {
		t.Errorf("expected the first source zone's SOA to be retained, got: %s", sharedZoneFromOne.ResourceRecords["shared.arpa."].Value)
	}
	if sharedZoneFromOne.ResourceRecords["2"] == nil || sharedZoneFromOne.ResourceRecords["2"].Value != "host-two.example.com." {
		t.Error("expected the second source zone's PTR record to be merged in")
	}
}

func TestRunE_Generate_MergesReverseZones_Conflict(t *testing.T) {
	setup(t)
	defer teardown(t)

	inputFile = "testing"
	outputDir = "testing-dir"

	zoneOne := &models.Zone{Config: &models.Config{GenerateReverseLookupZones: true}}
	zoneTwo := &models.Zone{Config: &models.Config{GenerateReverseLookupZones: true}}
	zones := map[string]*models.Zone{
		"one": zoneOne,
		"two": zoneTwo,
	}

	sharedZoneFromOne := &models.Zone{
		ResourceRecords: map[string]*models.ResourceRecord{
			"1": {Type: models.PTR, Name: "1", Value: "host-one.example.com."},
		},
	}
	sharedZoneFromTwo := &models.Zone{
		ResourceRecords: map[string]*models.ResourceRecord{
			"1": {Type: models.PTR, Name: "1", Value: "host-two.example.com."},
		},
	}

	mockParser.EXPECT().Parse(inputFile).Return(zones, nil)
	mockZoneReverser.EXPECT().ReverseZone("one", zoneOne).Return(map[string]*models.Zone{"shared.arpa.": sharedZoneFromOne}, nil)
	mockZoneReverser.EXPECT().ReverseZone("two", zoneTwo).Return(map[string]*models.Zone{"shared.arpa.": sharedZoneFromTwo}, nil)
	// No further calls are expected: the conflict is detected during pass 1, before Normalize or any writes.

	err := generateCmd.RunE(generateCmd, []string{})
	if err == nil {
		t.Fatal("expected an error, found none")
	}
}
