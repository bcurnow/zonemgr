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
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestPersistentPreRunE_Validate(t *testing.T) {
	setup(t)
	defer teardown(t)
	// We don't want to bother with the rootCmd persistentPreRunE so
	// we're just going to make sure it is called
	originalRootPPRE := rootCmd.PersistentPreRunE
	defer func() { rootCmd.PersistentPreRunE = originalRootPPRE }()

	testCases := []struct {
		absErr bool
	}{
		{},
		{absErr: true},
	}

	for _, tc := range testCases {
		rootPPRECalled := false
		rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			rootPPRECalled = true
			return nil
		}

		call := mockFs.EXPECT().ToAbsoluteFilePath("testing")
		if tc.absErr {
			call.Return("", errors.New("absErr"))
		} else {
			call.Return("testing", nil)
			mockPluginManager.EXPECT().Plugins().Return(testPlugins)
			mockPluginManager.EXPECT().Metadata().Return(testMetadata)
		}

		args := []string{"--input", "testing"}
		validateCmd.ParseFlags(args)
		if err := validateCmd.PersistentPreRunE(validateCmd, args); err != nil {
			if tc.absErr {
				want := "absErr"
				if err.Error() != want {
					t.Errorf("incorrect error: '%s', want: '%s'", err, want)
				}
			}
		} else {
			if !rootPPRECalled {
				t.Errorf("expected rootCmd PersistentPreRunE to be called, was not")
			}

			if inputFile != "testing" {
				t.Errorf("wrong value for input: '%s', want: '%s'", inputFile, "testing")
			}
		}
	}
}

func TestRunE_ValidateYaml(t *testing.T) {
	setup(t)
	defer teardown(t)
	testCases := []struct {
		parserErr bool
		want      string
	}{
		{want: "testing is valid\n"},
		{parserErr: true, want: "failed to parse input file testing: parserErr"},
	}

	for _, tc := range testCases {
		inputFile = "testing"
		originalStdout := os.Stdout
		defer func() { os.Stdout = originalStdout }()

		r, w, err := os.Pipe()
		if err != nil {
			t.Errorf("could not create pipe to capture stdout")
		}
		os.Stdout = w

		call := mockParser.EXPECT().Parse("testing")
		if tc.parserErr {
			call.Return(nil, errors.New("parserErr"))
		} else {
			call.Return(nil, nil)
		}

		err = validateYamlCmd.RunE(validateYamlCmd, []string{})
		w.Close()
		if err != nil {
			if tc.parserErr {
				if err.Error() != tc.want {
					t.Errorf("incorrect error: '%s', want: '%s'", err, tc.want)
				}
			}
		} else {
			if tc.parserErr {
				t.Error("expected an error, found none")
			} else {
				var buf bytes.Buffer
				buf.ReadFrom(r)
				r.Close()

				if buf.String() != tc.want {
					t.Errorf("incorrect output: '%s', want: '%s'", buf.String(), tc.want)
				}
			}
		}
	}
}
