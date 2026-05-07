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
				mockPluginManager.EXPECT().Plugins().Return(testPlugins).Times(2)
				mockPluginManager.EXPECT().Metadata().Return(testMetadata).Times(2)
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

			call = mockZoneFileGenerator.EXPECT().GenerateZone("one", zoneOne, outputDir)

			if tc.zoneFileGeneratorErr {
				call.Return(errors.New("zoneFileGeneratorErr"))
			} else {
				call.Return(nil)
				mockZoneFileGenerator.EXPECT().GenerateZone("two", zoneTwo, outputDir).Return(nil)

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
