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

package ctx

import (
	"fmt"
	"testing"

	"github.com/bcurnow/zonemgr/models/testingutils"
	"github.com/spf13/pflag"
)

var (
	mockEnvironment                    *MockEnvironment
	mockFlags                          *pflag.FlagSet
	mockGenerateReverseLookupZoneFlag  *pflag.Flag
	mockGenerateSerialFlag             *pflag.Flag
	mockPluginsDirectoryFlag           *pflag.Flag
	mockSerialChangeIndexDirectoryFlag *pflag.Flag
)

func testContextSetup(t *testing.T) {
	testingutils.Setup(t)
	mockEnvironment = NewMockEnvironment(testingutils.MockController)
	env = mockEnvironment
	pc = &pluginContext{}

	mockGenerateReverseLookupZoneFlag = &pflag.Flag{Name: FlagGenerateReverseLookupZone, Changed: false}
	mockGenerateSerialFlag = &pflag.Flag{Name: FlagGenerateSerial, Changed: false}
	mockPluginsDirectoryFlag = &pflag.Flag{Name: FlagPluginsDirectory, Changed: false}
	mockSerialChangeIndexDirectoryFlag = &pflag.Flag{Name: FlagSerialChangeIndexDirectory, Changed: false}

	mockFlags = pflag.NewFlagSet("testing", pflag.ContinueOnError)
	mockFlags.AddFlag(mockGenerateReverseLookupZoneFlag)
	mockFlags.AddFlag(mockGenerateSerialFlag)
	mockFlags.AddFlag(mockPluginsDirectoryFlag)
	mockFlags.AddFlag(mockSerialChangeIndexDirectoryFlag)
}

type tcFlag struct {
	generateReverseLookupZones        string
	generateReverseLookupZonesChanged bool
	generateSerial                    string
	generateSerialChanged             bool
	pluginsDirectory                  string
	pluginsDirectoryChanged           bool
	serialChangeIndexDirectory        string
	serialChangeIndexDirectoryChanged bool
}

type tcWant struct {
	generateReverseLookupZones bool
	generateSerial             bool
	pluginsDirectory           string
	serialChangeIndexDirectory string
}

func TestInitPluginContext(t *testing.T) {
	testContextSetup(t)
	defer testingutils.Teardown(t)

	testCases := []struct {
		want                          tcWant
		flag                          tcFlag
		generateReverseLookupZonesErr string
		generateSerialErr             string
		pluginsDirectoryErr           string
		serialChangeIndexDirectoryErr string
	}{
		{want: tcWant{generateReverseLookupZones: true}, flag: tcFlag{generateReverseLookupZones: "true", generateReverseLookupZonesChanged: true}},
		{want: tcWant{generateReverseLookupZones: true}, flag: tcFlag{generateReverseLookupZones: "true", generateReverseLookupZonesChanged: true}, generateReverseLookupZonesErr: "trying to get bool value of flag of type string"},
		{want: tcWant{generateReverseLookupZones: true}, flag: tcFlag{}},
		{want: tcWant{generateSerial: true}, flag: tcFlag{generateSerial: "true", generateSerialChanged: true}},
		{want: tcWant{generateSerial: true}, flag: tcFlag{generateSerial: "true", generateSerialChanged: true}, generateSerialErr: "trying to get bool value of flag of type string"},
		{want: tcWant{generateSerial: true}, flag: tcFlag{}},
		{want: tcWant{pluginsDirectory: "pluginsDir"}, flag: tcFlag{pluginsDirectory: "pluginsDir", pluginsDirectoryChanged: true}},
		{want: tcWant{pluginsDirectory: "pluginsDir"}, flag: tcFlag{pluginsDirectory: "pluginsDir", pluginsDirectoryChanged: true}, pluginsDirectoryErr: "trying to get string value of flag of type bool"},
		{want: tcWant{pluginsDirectory: "pluginsDir"}, flag: tcFlag{}},
		{want: tcWant{serialChangeIndexDirectory: "sciDir"}, flag: tcFlag{serialChangeIndexDirectory: "sciDir", serialChangeIndexDirectoryChanged: true}},
		{want: tcWant{serialChangeIndexDirectory: "sciDir"}, flag: tcFlag{serialChangeIndexDirectory: "sciDir", serialChangeIndexDirectoryChanged: true}, serialChangeIndexDirectoryErr: "trying to get string value of flag of type bool"},
		{want: tcWant{serialChangeIndexDirectory: "sciDir"}, flag: tcFlag{}},
	}

	for _, tc := range testCases {
		mockEnvironment.EXPECT().LoadValues()
		// There's a lot of duplication here but it basically all comes down to
		// If the flag value has changed, we don't need to mock out the environment calls
		// Always make sure the .Changed value for the flag handled in the block is set to the correct value
		// If there's an error, change the flag value to generate one or mock out an error call
		// If there's an error before this, don't mock out the calls
		//
		if tc.flag.generateReverseLookupZonesChanged {
			mockGenerateReverseLookupZoneFlag.Changed = true
			if tc.generateReverseLookupZonesErr == "" {
				mockGenerateReverseLookupZoneFlag.Value = &testFlagValue{value: tc.flag.generateReverseLookupZones, flagType: "bool"}
			} else {
				mockGenerateReverseLookupZoneFlag.Value = &testFlagValue{value: "not a boolean", flagType: "string"}
			}
		} else {
			mockGenerateReverseLookupZoneFlag.Changed = false
			mockEnvironment.EXPECT().GenerateReverseLookupZones().Return(tc.want.generateReverseLookupZones)
		}
		if tc.flag.generateSerialChanged {
			mockGenerateSerialFlag.Changed = true
			if tc.generateSerialErr == "" {
				mockGenerateSerialFlag.Value = &testFlagValue{value: tc.flag.generateSerial, flagType: "bool"}
			} else {
				mockGenerateSerialFlag.Value = &testFlagValue{value: "not a boolean", flagType: "string"}
			}
		} else {
			mockGenerateSerialFlag.Changed = false
			if tc.generateReverseLookupZonesErr == "" {
				mockEnvironment.EXPECT().GenerateSerial().Return(tc.want.generateSerial)
			}
		}
		if tc.flag.pluginsDirectoryChanged {
			mockPluginsDirectoryFlag.Changed = true
			if tc.pluginsDirectoryErr == "" {
				mockPluginsDirectoryFlag.Value = &testFlagValue{value: "pluginsDir", flagType: "string"}
			} else {
				mockPluginsDirectoryFlag.Value = &testFlagValue{value: "true", flagType: "bool"}
			}
		} else {
			mockPluginsDirectoryFlag.Changed = false
			if tc.generateReverseLookupZonesErr == "" && tc.generateSerialErr == "" {
				mockEnvironment.EXPECT().PluginsDirectory().Return(tc.want.pluginsDirectory)
			}
		}
		if tc.flag.serialChangeIndexDirectoryChanged {
			mockSerialChangeIndexDirectoryFlag.Changed = true
			if tc.serialChangeIndexDirectoryErr == "" {
				mockSerialChangeIndexDirectoryFlag.Value = &testFlagValue{value: "sciDir", flagType: "string"}
			} else {
				mockSerialChangeIndexDirectoryFlag.Value = &testFlagValue{value: "true", flagType: "bool"}
			}
		} else {
			mockSerialChangeIndexDirectoryFlag.Changed = false
			if tc.generateReverseLookupZonesErr == "" && tc.generateSerialErr == "" && tc.pluginsDirectoryErr == "" {
				mockEnvironment.EXPECT().SerialChangeIndexDirectory().Return(tc.want.serialChangeIndexDirectory)
			}
		}

		if err := InitPluginContext(mockFlags); err != nil {
			if tc.generateReverseLookupZonesErr != "" {
				if err.Error() != tc.generateReverseLookupZonesErr {
					t.Errorf("incorrect error: %s, want %s", err, tc.generateReverseLookupZonesErr)
				}
			} else if tc.generateSerialErr != "" {
				if err.Error() != tc.generateSerialErr {
					t.Errorf("incorrect error: %s, want %s", err, tc.generateSerialErr)
				}
			} else if tc.pluginsDirectoryErr != "" {
				if err.Error() != tc.pluginsDirectoryErr {
					t.Errorf("incorrect error: %s, want %s", err, tc.pluginsDirectoryErr)
				}
			} else if tc.serialChangeIndexDirectoryErr != "" {
				if err.Error() != tc.serialChangeIndexDirectoryErr {
					t.Errorf("incorrect error: %s, want %s", err, tc.serialChangeIndexDirectoryErr)
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if C().GenerateReverseLookupZones() != tc.want.generateReverseLookupZones {
				t.Errorf("GenerateReverseLookupZones=%t, want %t", C().GenerateReverseLookupZones(), tc.want.generateReverseLookupZones)
			}
			if C().GenerateSerial() != tc.want.generateSerial {
				t.Errorf("GenerateSerial=%t, want %t", C().GenerateSerial(), tc.want.generateSerial)
			}
			if C().PluginsDirectory() != tc.want.pluginsDirectory {
				t.Errorf("PluginsDirectory=%s, want %s", C().PluginsDirectory(), tc.want.pluginsDirectory)
			}
			if C().SerialChangeIndexDirectory() != tc.want.serialChangeIndexDirectory {
				t.Errorf("SerialChangeIndexDirectory=%s, want %s", C().SerialChangeIndexDirectory(), tc.want.serialChangeIndexDirectory)
			}
		}
	}
}

type testFlagValue struct {
	pflag.Value
	value    string
	flagType string
}

func (f *testFlagValue) String() string {
	return f.value
}

func (f *testFlagValue) Set(string) error {
	return fmt.Errorf("Set not implemented on test flags")
}

func (f *testFlagValue) Type() string {
	return f.flagType
}
