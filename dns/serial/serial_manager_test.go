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

	"go.uber.org/mock/gomock"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/utils"
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

var _ utils.YamlFile[*models.SerialIndex] = &TestYamlFile{}

const testZoneName = "testing"

func newFsm(t *testing.T) fileSerialManager {
	t.Helper()
	return fileSerialManager{changeIndexDirectory: t.TempDir(), indexFile: &TestYamlFile{}}
}

func serialFilePath(fsm fileSerialManager) string {
	return filepath.Join(fsm.changeIndexDirectory, fmt.Sprintf("%s.%s", testZoneName, changeIndexFileExtension))
}

func setupGenerator(t *testing.T) (*gomock.Controller, *MockGenerator) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockGen := NewMockGenerator(ctrl)
	generator = mockGen
	t.Cleanup(func() { generator = &TimeBasedGenerator{} })
	return ctrl, mockGen
}

func TestNext(t *testing.T) {
	t.Run("success-new-file", func(t *testing.T) {
		_, mockGen := setupGenerator(t)
		fsm := newFsm(t)
		fs = &utils.FileSystem{}
		defer func() { fs = &utils.FileSystem{} }()

		mockGen.EXPECT().GenerateBase().Return(toUint32Ptr(12345678), nil)
		mockGen.EXPECT().FromSerialIndex(&models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(1)}).Return("1234567801", nil)

		serial, err := fsm.Next(testZoneName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if serial != "1234567801" {
			t.Errorf("got serial %q, want %q", serial, "1234567801")
		}
	})

	t.Run("success-existing-file", func(t *testing.T) {
		_, mockGen := setupGenerator(t)
		fsm := newFsm(t)
		fs = &utils.FileSystem{}
		defer func() { fs = &utils.FileSystem{} }()

		// Pre-create the serial file so Exists returns true
		filePath := serialFilePath(fsm)
		if err := os.WriteFile(filePath, []byte{}, 0600); err != nil {
			t.Fatalf("failed to create serial file: %v", err)
		}

		mockGen.EXPECT().GenerateBase().Return(toUint32Ptr(12345678), nil)
		mockGen.EXPECT().FromSerialIndex(&models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(33)}).Return("1234567833", nil)

		serial, err := fsm.Next(testZoneName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if serial != "1234567833" {
			t.Errorf("got serial %q, want %q", serial, "1234567833")
		}
	})

	t.Run("mkdir-error", func(t *testing.T) {
		fs = &utils.FileSystem{}
		defer func() { fs = &utils.FileSystem{} }()

		// Create a file at the directory path so MkdirAll fails
		conflictPath := filepath.Join(t.TempDir(), "conflict")
		if err := os.WriteFile(conflictPath, []byte("conflict"), 0600); err != nil {
			t.Fatalf("failed to create conflict file: %v", err)
		}
		fsm := fileSerialManager{changeIndexDirectory: conflictPath, indexFile: &TestYamlFile{}}

		_, err := fsm.Next(testZoneName)
		if err == nil {
			t.Error("expected an error, got none")
		}
	})

	t.Run("flock-error-new-file", func(t *testing.T) {
		fs = &utils.FileSystem{}
		defer func() { fs = &utils.FileSystem{} }()

		// Make the directory read-only so Flock can't create the file
		readOnlyDir := t.TempDir()
		if err := os.Chmod(readOnlyDir, 0555); err != nil {
			t.Fatalf("failed to chmod: %v", err)
		}
		t.Cleanup(func() { os.Chmod(readOnlyDir, 0755) })
		fsm := fileSerialManager{changeIndexDirectory: readOnlyDir, indexFile: &TestYamlFile{}}

		_, err := fsm.Next(testZoneName)
		if err == nil {
			t.Error("expected an error, got none")
		}
	})

	t.Run("flock-error-existing-file", func(t *testing.T) {
		fs = &utils.FileSystem{}
		defer func() { fs = &utils.FileSystem{} }()

		fsm := newFsm(t)
		// Create a directory at the serial file path so Exists returns true and Flock fails
		dirAtFilePath := serialFilePath(fsm)
		if err := os.MkdirAll(dirAtFilePath, 0755); err != nil {
			t.Fatalf("failed to create dir at file path: %v", err)
		}

		_, err := fsm.Next(testZoneName)
		if err == nil {
			t.Error("expected an error, got none")
		}
	})

	t.Run("generate-error", func(t *testing.T) {
		_, mockGen := setupGenerator(t)
		fsm := newFsm(t)
		fs = &utils.FileSystem{}
		defer func() { fs = &utils.FileSystem{} }()

		mockGen.EXPECT().GenerateBase().Return(toUint32Ptr(12345678), nil)
		mockGen.EXPECT().FromSerialIndex(gomock.Any()).Return("", errors.New("generateErr"))

		_, err := fsm.Next(testZoneName)
		if err == nil || err.Error() != "generateErr" {
			t.Errorf("got error %v, want %q", err, "generateErr")
		}
	})
}

func TestInitFile(t *testing.T) {
	testCases := []struct {
		name        string
		flockErr    bool
		generateErr bool
		writeErr    bool
	}{
		{name: "success"},
		{name: "flock-error", flockErr: true},
		{name: "generate-error", generateErr: true},
		{name: "write-error", writeErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, mockGen := setupGenerator(t)
			defer ctrl.Finish()
			fs = &utils.FileSystem{}
			defer func() { fs = &utils.FileSystem{} }()

			var path string
			if tc.flockErr {
				// Create a directory at the path to make Flock fail
				dirPath := filepath.Join(t.TempDir(), "test.serial")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}
				path = dirPath
			} else {
				path = filepath.Join(t.TempDir(), "test.serial")
			}

			fsm := fileSerialManager{changeIndexDirectory: t.TempDir(), indexFile: &TestYamlFile{}}
			if !tc.flockErr {
				call := mockGen.EXPECT().GenerateBase()
				if tc.generateErr {
					call.Return(nil, errors.New("generateErr"))
				} else {
					call.Return(toUint32Ptr(12345678), nil)
					if tc.writeErr {
						fsm.indexFile = &TestYamlFile{writeErr: true}
					}
				}
			}

			si, err := fsm.initFile(path)
			if err != nil {
				if tc.flockErr || tc.generateErr {
					return // expected error, any message is fine for filesystem errors
				}
				if tc.writeErr {
					if err.Error() != "writeErr" {
						t.Errorf("got error %q, want %q", err.Error(), "writeErr")
					}
					return
				}
				t.Errorf("unexpected error: %v", err)
			} else {
				if tc.flockErr || tc.generateErr || tc.writeErr {
					t.Error("expected an error, got none")
				} else {
					want := &models.SerialIndex{Base: toUint32Ptr(12345678), ChangeIndex: toUint32Ptr(1)}
					if !cmp.Equal(si, want) {
						t.Errorf("incorrect result:\n%s", cmp.Diff(si, want))
					}
				}
			}
		})
	}
}

func TestIncrementAndUpdate(t *testing.T) {
	testCases := []struct {
		name        string
		flockErr    bool
		readErr     bool
		generateErr bool
		diffBase    bool
		writeErr    bool
	}{
		{name: "success"},
		{name: "different-base", diffBase: true},
		{name: "flock-error", flockErr: true},
		{name: "read-error", readErr: true},
		{name: "generate-error", generateErr: true},
		{name: "write-error", writeErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, mockGen := setupGenerator(t)
			defer ctrl.Finish()
			fs = &utils.FileSystem{}
			defer func() { fs = &utils.FileSystem{} }()

			var path string
			if tc.flockErr {
				dirPath := filepath.Join(t.TempDir(), "test.serial")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}
				path = dirPath
			} else {
				path = filepath.Join(t.TempDir(), "test.serial")
			}

			indexFile := &TestYamlFile{}
			if tc.readErr {
				indexFile = &TestYamlFile{readErr: true}
			}
			fsm := fileSerialManager{changeIndexDirectory: t.TempDir(), indexFile: indexFile}

			if !tc.flockErr && !tc.readErr {
				call := mockGen.EXPECT().GenerateBase()
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
					}
				}
			}

			si, err := fsm.incrementAndUpdate(path)
			if err != nil {
				if tc.flockErr || tc.generateErr {
					return // expected, OS-specific message
				}
				if tc.readErr {
					if err.Error() != "readErr" {
						t.Errorf("got error %q, want %q", err.Error(), "readErr")
					}
					return
				}
				if tc.writeErr {
					if err.Error() != "writeErr" {
						t.Errorf("got error %q, want %q", err.Error(), "writeErr")
					}
					return
				}
				t.Errorf("unexpected error: %v", err)
			} else {
				if tc.flockErr || tc.generateErr || tc.writeErr {
					t.Error("expected an error, got none")
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
		})
	}
}
