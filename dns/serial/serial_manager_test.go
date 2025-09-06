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

package serial

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/gofrs/flock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

type TestYamlFile struct {
	readErr  bool
	writeErr bool
}

func (t *TestYamlFile) Read(path string) (*models.SerialIndex, error) {
	if t.readErr {
		return nil, errors.New("readErr")
	}
	return &models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(32)}, nil
}

func (t *TestYamlFile) Write(path string, content *models.SerialIndex) error {
	if t.writeErr {
		return errors.New("writeErr")
	}
	return nil
}

var (
	_                     utils.YamlFile[*models.SerialIndex] = &TestYamlFile{}
	mockController        *gomock.Controller
	mockGenerator         *MockGenerator
	mockFs                *utils.MockFileSystemOperations
	fsm                   fileSerialManager
	testDir               string
	testZoneName          = "testing"
	serialChangeIndexFile string
)

func setup_SerialManager(t *testing.T) {
	mockController = gomock.NewController(t)
	mockGenerator = NewMockGenerator(mockController)
	mockFs = utils.NewMockFileSystemOperations(mockController)

	generator = mockGenerator
	fs = mockFs

	testDir, _ = os.MkdirTemp("", t.Name())
	serialChangeIndexFile = filepath.Join(testDir, fmt.Sprintf("%s.%s", testZoneName, changeIndexFileExtension))
	fsm = fileSerialManager{changeIndexDirectory: testDir, indexFile: &TestYamlFile{}}
}

func teardown_SerialManager(_ *testing.T) {
	generator = &TimeBasedGenerator{}
	fs = &utils.FileSystem{}
	mockController.Finish()
}

func TestNext(t *testing.T) {
	setup_SerialManager(t)
	defer teardown_SerialManager(t)
	testCases := []struct {
		name                  string
		mkdirErr              bool
		exists                bool
		incrementAndUpdateErr bool
		initFileErr           bool
		generateErr           bool
	}{
		{name: "success"},
		{name: "succes-exists", exists: true},
		{name: "mkdirErr", mkdirErr: true},
		{name: "incrementAndUpdateErr", incrementAndUpdateErr: true, exists: true},
		{name: "initFileErr", initFileErr: true},
		{name: "generateErr", generateErr: true},
	}

	for _, tc := range testCases {
		call := mockFs.EXPECT().MkdirAll(fsm.changeIndexDirectory, os.FileMode(0750))
		if tc.mkdirErr {
			call.Return(errors.New("mkdirErr"))
		} else {
			// Keep track of any errors we caused in the exists block
			errored := false
			if tc.exists {
				mockFs.EXPECT().Exists(serialChangeIndexFile).Return(true)
				if tc.incrementAndUpdateErr {
					mockFs.EXPECT().Flock(serialChangeIndexFile).Return(nil, errors.New("incrementAndUpdateErr"))
					errored = true
				}
			} else {
				mockFs.EXPECT().Exists(serialChangeIndexFile).Return(false)
				if tc.initFileErr {
					mockFs.EXPECT().Flock(serialChangeIndexFile).Return(nil, errors.New("initFileErr"))
					errored = true
				}
			}

			// If we've errored, the method is over an none of the below will ever happen
			if !errored {
				// Now we have our serialIndex and can continue on the happy path
				testFlock := flock.New(serialChangeIndexFile)
				mockFs.EXPECT().Flock(serialChangeIndexFile).Return(testFlock, nil)
				mockGenerator.EXPECT().GenerateBase().Return(toUint32Ptr(12345678), nil)

				if tc.generateErr {
					mockGenerator.EXPECT().FromSerialIndex(gomock.Any()).Return("", errors.New("generateErr"))
				} else {
					if tc.exists {
						mockGenerator.EXPECT().FromSerialIndex(&models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(33)}).Return("1234567833", nil)
					} else {
						mockGenerator.EXPECT().FromSerialIndex(&models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(1)}).Return("1234567801", nil)
					}
				}
			}
		}
		serial, err := fsm.Next(testZoneName)
		if err != nil {
			want := ""
			if tc.mkdirErr {
				want = "mkdirErr"
			} else if tc.incrementAndUpdateErr {
				want = "incrementAndUpdateErr"
			} else if tc.initFileErr {
				// This will come from the TestYamlFile impl
				want = "initFileErr"
			} else if tc.generateErr {
				want = "generateErr"
			}

			if err.Error() != want {
				t.Errorf("%s - incorrect error: '%s', want: '%s'", tc.name, err, want)
			}
		} else {
			if tc.mkdirErr || tc.incrementAndUpdateErr || tc.initFileErr {
				t.Errorf("%s - expected an error, found none", tc.name)
			} else {
				want := "1234567801"
				if tc.exists {
					want = "1234567833"
				}
				if serial != want {
					t.Errorf("%s - incorrect serial: '%s', want: '%s'", tc.name, serial, want)
				}
			}
		}
	}
}

func TestInitFile(t *testing.T) {
	setup_SerialManager(t)
	defer teardown_SerialManager(t)
	testCases := []struct {
		flockErr    bool
		generateErr bool
		writeErr    bool
	}{
		{},
		{flockErr: true},
		{generateErr: true},
		{writeErr: true},
	}

	for _, tc := range testCases {
		call := mockFs.EXPECT().Flock("testing")
		if tc.flockErr {
			call.Return(nil, errors.New("flockErr"))
		} else {
			flock := flock.New("testing")
			call.Return(flock, nil)

			call = mockGenerator.EXPECT().GenerateBase()
			if tc.generateErr {
				call.Return(nil, errors.New("generateErr"))
			} else {
				call.Return(toUint32Ptr(12345678), nil)
				if tc.writeErr {
					fsm.indexFile = &TestYamlFile{writeErr: true}
				}
			}
		}
		si, err := fsm.initFile("testing")
		if err != nil {
			want := ""
			if tc.flockErr {
				want = "flockErr"
			} else if tc.generateErr {
				want = "generateErr"
			} else if tc.writeErr {
				want = "writeErr"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.flockErr || tc.generateErr || tc.writeErr {
				t.Error("expected an error, found none")
			} else {
				if !cmp.Equal(si, &models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(1)}) {
					t.Errorf("incorrect result:\n%s", cmp.Diff(si, &models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(1)}))
				}
			}
		}
	}
}

func TestIncrementAndUpdate(t *testing.T) {
	setup_SerialManager(t)
	defer teardown_SerialManager(t)
	testCases := []struct {
		flockErr    bool
		readErr     bool
		generateErr bool
		diffBase    bool
		writeErr    bool
	}{
		{},
		{diffBase: true},
		{flockErr: true},
		{readErr: true},
		{generateErr: true},
		{writeErr: true},
	}

	for _, tc := range testCases {
		call := mockFs.EXPECT().Flock("testing")
		if tc.flockErr {
			call.Return(nil, errors.New("flockErr"))
		} else {
			flock := flock.New("testing")
			call.Return(flock, nil)

			if tc.readErr {
				fsm.indexFile = &TestYamlFile{readErr: true}
			} else {
				fsm.indexFile = &TestYamlFile{}
				call = mockGenerator.EXPECT().GenerateBase()
				if tc.generateErr {
					call.Return(nil, errors.New("generateErr"))
				} else {
					if tc.diffBase {
						call.Return(toUint32Ptr(123456789), nil)
					} else {
						call.Return(toUint32Ptr(12345678), nil)
					}
					if tc.writeErr {
						fsm.indexFile = &TestYamlFile{writeErr: true}
					} else {
						si := &models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(33)}
						if tc.diffBase {
							si.Base = toUint32Ptr(123456789)
							si.ChangeIndex = toUint32Ptr(1)
						}
					}
				}
			}
		}
		si, err := fsm.incrementAndUpdate("testing")
		if err != nil {
			want := ""
			if tc.flockErr {
				want = "flockErr"
			} else if tc.readErr {
				want = "readErr"
			} else if tc.generateErr {
				want = "generateErr"
			} else if tc.writeErr {
				want = "writeErr"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.flockErr || tc.generateErr || tc.writeErr {
				t.Error("expected an error, found none")
			} else {
				want := &models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(33)}
				if tc.diffBase {
					want.Base = toUint32Ptr(123456789)
					want.ChangeIndex = toUint32Ptr(1)
				}
				if !cmp.Equal(si, want) {
					t.Errorf("incorrect result:\n%s", cmp.Diff(si, want))
				}
			}
		}
	}
}
