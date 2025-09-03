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

package utils

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRead_ZoneYamlFile(t *testing.T) {

	readFile = func(name string) ([]byte, error) { return []byte("testing"), nil }
	unmarshal = func(in []byte, out interface{}) (err error) { return nil }

	f := (&ZoneYamlFile{})

	zone, err := f.Read("testing")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if zone != nil {
		t.Errorf("incorrect result: %s, expected nil", zone)
	}
}

func TestWrite_ZoneYamlFile(t *testing.T) {
	if err := (&ZoneYamlFile{}).Write("testing", nil); err == nil {
		t.Error("expected an error, found none")
	} else {
		if err.Error() != "not implemented" {
			t.Errorf("incorrect error: '%s', want: 'not implemented'", err)
		}
	}
}

func TestRead_SerialIndexYamlFile(t *testing.T) {

	readFile = func(_ string) ([]byte, error) { return []byte("testing"), nil }
	unmarshal = func(_ []byte, _ interface{}) (err error) { return nil }

	f := (&SerialIndexYamlFile{})

	serialIndex, err := f.Read("testing")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if serialIndex != nil {
		t.Errorf("incorrect result: %v, expected nil", serialIndex)
	}
}

func TestWrite_SerialIndexYamlFile(t *testing.T) {
	createTemp(t)
	defer tempTeardown(t)
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) { return testFile, nil }
	marshal = func(_ interface{}) (out []byte, err error) { return []byte("testing"), nil }

	if err := (&SerialIndexYamlFile{}).Write("testing", nil); err != nil {
		t.Errorf("unexpected err: %s", err)
	} else {
		content, err := os.ReadFile(testFile.Name())
		if err != nil {
			t.Errorf("unable to read test file '%s': %s", testFile.Name(), err)
		} else {
			if !slices.Equal(content, []byte("testing")) {
				t.Errorf("incorrect file contents: '%s', want: 'testing'", content)
			}
		}
	}
}

func TestUnmarshalYaml(t *testing.T) {
	testCases := []struct {
		readFileErr  bool
		unmarshalErr bool
	}{
		{},
		{readFileErr: true},
		{unmarshalErr: true},
	}

	for _, tc := range testCases {
		if tc.readFileErr {
			readFile = func(_ string) ([]byte, error) { return nil, errors.New("readFileErr") }
		} else {
			readFile = func(_ string) ([]byte, error) { return []byte("testing"), nil }
		}

		if tc.unmarshalErr {
			unmarshal = func(in []byte, _ interface{}) (err error) { return errors.New("unmarshalErr") }
		} else {
			unmarshal = func(in []byte, _ interface{}) (err error) {
				if !slices.Equal(in, []byte("testing")) {
					return fmt.Errorf("unexpected input bytes: '%s', want: 'testing'", in)
				}
				return nil
			}
		}

		o, err := unmarshalYaml[string]("testing")
		if err != nil {
			want := ""
			if tc.readFileErr {
				want = "failed to open 'testing': readFileErr"
			} else if tc.unmarshalErr {
				want = "failed to parse input YAML: unmarshalErr"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.readFileErr || tc.unmarshalErr {
				t.Error("expected an error, found none")
			}

			if o != "" {
				t.Errorf("incorrect result: '%s', expected empty string", o)
			}
		}
	}
}

func TestMarshalYaml(t *testing.T) {
	defer tempTeardown(t)
	testCases := []struct {
		openErr    bool
		marshalErr bool
		writeErr   bool
	}{
		{},
		{openErr: true},
		{marshalErr: true},
		{writeErr: true},
	}

	for _, tc := range testCases {
		createTemp(t)
		defer tempTeardown(t)

		if tc.openErr {
			openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) { return nil, errors.New("openErr") }
		} else {
			openFile = os.OpenFile
		}

		if tc.marshalErr {
			marshal = func(_ interface{}) (out []byte, err error) { return nil, errors.New("marshalErr") }
		} else {
			marshal = yaml.Marshal
		}

		if tc.writeErr {
			// To simulate a write error, we're going to open the file for read-only
			marshalFileMode = os.O_RDONLY | os.O_CREATE | os.O_APPEND
		}

		if err := marshalYaml(testFile.Name(), "content"); err != nil {
			want := ""
			if tc.openErr {
				want = fmt.Sprintf("failed to open '%s': openErr", testFile.Name())
			} else if tc.marshalErr {
				want = "failed to marshal 'content': marshalErr"
			} else if tc.writeErr {
				want = fmt.Sprintf("failed to write to '%s': write %s: bad file descriptor", testFile.Name(), testFile.Name())
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.openErr || tc.marshalErr || tc.writeErr {
				t.Error("expected an error, found none")
			} else {
				content, err := os.ReadFile(testFile.Name())
				if err != nil {
					t.Errorf("unable to read test file '%s': %s", testFile.Name(), err)
				} else {
					// yaml.Marshal adds a newline to the end
					if !slices.Equal(content, []byte("content\n")) {
						t.Errorf("incorrect file contents: '%s', want: 'content'", content)
					}
				}
			}
		}
	}
}
