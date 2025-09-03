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
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/gofrs/flock"
)

type infoErrorDirEntry struct {
}

var _ os.DirEntry = &infoErrorDirEntry{}

func (d *infoErrorDirEntry) Name() string {
	return "testing"
}

func (d *infoErrorDirEntry) IsDir() bool {
	return false
}

func (d *infoErrorDirEntry) Type() fs.FileMode {
	return 0755
}

func (d *infoErrorDirEntry) Info() (fs.FileInfo, error) {
	return nil, errors.New("dinfoErr")
}

func TestCreateFile(t *testing.T) {
	defer tempTeardown(t)
	testCases := []struct {
		createErr  bool
		chmodErr   bool
		contentErr bool
		writeErr   bool
	}{
		{},
		{createErr: true},
		{chmodErr: true},
		{contentErr: true},
		{writeErr: true},
	}

	for _, tc := range testCases {
		createTemp(t)
		defer tempTeardown(t)
		if tc.createErr {
			create = func(name string) (*os.File, error) {
				if name != testFile.Name() {
					return nil, errors.New("wrong name: " + name)
				} else {
					return nil, errors.New("createErr")
				}
			}
		} else {
			create = func(_ string) (*os.File, error) { return testFile, nil }
		}

		if tc.chmodErr {
			chmod = func(name string, mode os.FileMode) error {
				if name != testFile.Name() && mode != 0755 {
					return fmt.Errorf("incorrect chmod: name=%s, mode=%o", name, mode)
				}
				return errors.New("chmodErr")
			}
		} else {
			chmod = func(_ string, _ os.FileMode) error { return nil }
		}

		contentFn := func() ([]byte, error) {
			return []byte("content"), nil
		}

		if tc.contentErr {
			contentFn = func() ([]byte, error) { return nil, errors.New("contentErr") }
		}

		if tc.writeErr {
			if err := testFile.Close(); err != nil {
				t.Errorf("error closing file to create a writeErr: %s", err)
			}
		}

		if err := (&FileSystem{}).CreateFile(testFile.Name(), 0755, contentFn); err != nil {
			want := ""
			if tc.createErr {
				want = fmt.Sprintf("failed to create output file '%s': createErr", testFile.Name())
			} else if tc.chmodErr {
				want = fmt.Sprintf("error with chmod of '%s' to '755': chmodErr", testFile.Name())
			} else if tc.contentErr {
				want = fmt.Sprintf("error generating content for output file '%s': contentErr", testFile.Name())
			} else if tc.writeErr {
				want = fmt.Sprintf("error writing content for output file '%s': write %s: file already closed", testFile.Name(), testFile.Name())
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.createErr || tc.chmodErr || tc.contentErr || tc.writeErr {
				t.Error("expected an error, found none")
			} else {
				content, err := os.ReadFile(testFile.Name())
				if err != nil {
					t.Errorf("Unable to read test file '%s': %s", testFile.Name(), err)
				} else {
					if !slices.Equal(content, []byte("content")) {
						t.Errorf("incorrect content in test file '%s': '%s', want '%s'", testFile.Name(), content, "content")
					}
				}
			}
		}
	}
}

func TestExists(t *testing.T) {
	testCases := []struct {
		want bool
	}{
		{true},
		{false},
	}

	for _, tc := range testCases {
		stat = func(_ string) (os.FileInfo, error) {
			if tc.want {
				return nil, nil
			} else {
				return nil, os.ErrNotExist
			}
		}

		if tc.want != (&FileSystem{}).Exists("testing") {
			t.Errorf("incorrect result %v, wanted %v", (&FileSystem{}).Exists("testing"), tc.want)
		}
	}
}

func TestFlock(t *testing.T) {
	defer tempTeardown(t)
	testCases := []struct {
		tryErr      bool
		isNotLocked bool
	}{
		{},
		{true, false},
		{false, true},
	}

	for _, tc := range testCases {
		createTemp(t)
		defer tempTeardown(t)

		fileLock := flock.New(testFile.Name())
		defer fileLock.Close()
		if tc.tryErr {
			// remove the permissions to the file
			os.Chmod(testFile.Name(), 0000)
		} else if tc.isNotLocked {
			if err := fileLock.Lock(); err != nil {
				t.Errorf("error locking file for testing: %s", err)
			}
		}

		_, err := (&FileSystem{}).Flock(testFile.Name())
		fileLock.Close()
		if err != nil {
			want := ""
			if tc.tryErr {
				want = fmt.Sprintf("open %s: permission denied", testFile.Name())
			} else if tc.isNotLocked {
				want = fmt.Sprintf("unable to acquire exclusive lock on '%s'", testFile.Name())
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.tryErr || tc.isNotLocked {
				t.Error("expected an error, found none")
			}
		}
	}
}

func TestHomeDir(t *testing.T) {
	originalHomeDir := homeDir
	defer func() { homeDir = originalHomeDir }()
	want := "bogus"
	homeDir = want

	actual := (&FileSystem{}).HomeDir()
	if actual != want {
		t.Errorf("incorrect value: '%s', want: '%s'", actual, want)
	}
}

func TestMkdirAll(t *testing.T) {
	defer func() { mkdirAll = os.MkdirAll }()
	mkdirAll = func(path string, perm os.FileMode) error {
		if path != "testing" && perm != 0755 {
			return fmt.Errorf("wrong values: path=%s, perm=%o", path, perm)
		}
		return nil
	}

	if err := (&FileSystem{}).MkdirAll("testing", 0755); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestToAbsoluteFilePath(t *testing.T) {
	testCases := []struct {
		hasPrefix bool
		absErr    bool
	}{
		{},
		{hasPrefix: true},
		{absErr: true},
	}

	for _, tc := range testCases {
		originalHomeDir := homeDir
		defer func() { homeDir = originalHomeDir }()
		defer func() { abs = filepath.Abs }()
		path := "testing"
		if tc.hasPrefix {
			path = "~/testing"
			homeDir = "homeDir"
		}

		if tc.absErr {
			abs = func(path string) (string, error) {
				return "", errors.New("absErr")
			}
		} else {
			abs = func(path string) (string, error) {
				return path, nil
			}
		}

		absPath, err := (&FileSystem{}).ToAbsoluteFilePath(path)
		if err != nil {
			if tc.absErr {
				if err.Error() != "absErr" {
					t.Errorf("incorrect error: '%s', want: '%s'", err, "absErr")
				}
			} else {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			if tc.absErr {
				t.Error("expected an error, found none")
			}
			want := "testing"
			if tc.hasPrefix {
				want = "homeDir/testing"
			}
			if absPath != want {
				t.Errorf("incorrect path: '%s', want: '%s'", absPath, want)
			}
		}
	}
}

func TestWalkExecutables(t *testing.T) {
	defer tempTeardown(t)
	testCases := []struct {
		notExist       bool
		walkErr        bool
		absErr         bool
		includeSubDirs bool
		dInfoErr       bool
	}{
		{},
		{notExist: true},
		{walkErr: true},
		{absErr: true},
		{includeSubDirs: true},
		{dInfoErr: true},
	}

	for _, tc := range testCases {
		createTemp(t)
		defer tempTeardown(t)

		// We've now got a directory with a single file (that is not executable) in it
		// Let's setup some other cases

		// Create a directory with an executable in it
		if err := os.MkdirAll(filepath.Join(testDir, "sub-dir-with-executable"), 0755); err != nil {
			t.Errorf("unable to setup test, can't create sub-dir-with-executable")
		}

		testExecutable, err := os.Create(filepath.Join(testDir, "sub-dir-with-executable", "executable"))
		if err != nil {
			t.Errorf("unable to setup test, can't create executable in sub-dir-with-executable: %s", err)
		} else {
			testExecutable.Close()
			if err := os.Chmod(testExecutable.Name(), 0755); err != nil {
				t.Errorf("unable to setup test, can't change permissions on %s: %s", testExecutable.Name(), err)
			}
		}

		// Create an executable in the testDir
		testExecutable, err = os.Create(filepath.Join(testDir, "executable"))
		if err != nil {
			t.Errorf("unable to setup test, can't create executable in sub-dir-with-executable: %s", err)
		} else {
			testExecutable.Close()
			if err := os.Chmod(testExecutable.Name(), 0755); err != nil {
				t.Errorf("unable to setup test, can't change permissions on %s: %s", testExecutable.Name(), err)
			}
		}

		if tc.walkErr {
			if err := os.Chmod(testDir, 0000); err != nil {
				t.Errorf("unable to setup test, can't remove perms from %s: %s", testDir, err)
			}
		}

		if tc.dInfoErr {
			// We need to actually replace the core waltDir function for this
			// as we need to pass in a fake os.DirEntry in order to trigger the error
			walkDir = func(root string, fn fs.WalkDirFunc) error {
				return fn(root, &infoErrorDirEntry{}, nil)
			}
		}

		if tc.absErr {
			abs = func(path string) (string, error) { return "", errors.New("absErr") }
		} else {
			abs = func(path string) (string, error) { return path, nil }
		}

		dir := testDir
		if tc.notExist {
			dir = "does-not-exist"
		}
		executables, err := (&FileSystem{}).WalkExecutables(dir, tc.includeSubDirs)
		if err != nil {
			want := ""
			if tc.absErr {
				want = "absErr"
			} else if tc.walkErr {
				want = fmt.Sprintf("open %s: permission denied", testDir)
			} else if tc.dInfoErr {
				want = "dinfoErr"
			}

			if err.Error() != want {
				t.Errorf("incorrect error: '%s', want: '%s'", err, want)
			}
		} else {
			if tc.absErr || tc.walkErr || tc.dInfoErr {
				t.Error("expected an error, found none")
			}
			wantedExecutables := 1

			if tc.includeSubDirs {
				wantedExecutables = 2
			} else if tc.notExist {
				wantedExecutables = 0
			}

			if len(executables) != wantedExecutables {
				t.Errorf("incorrect number of executables: %d, wanted: %d", len(executables), wantedExecutables)
			}
		}
	}
}

func TestDetermineHomeDir(t *testing.T) {
	defer func() { userHomeDir = os.UserHomeDir; getWd = os.Getwd }()
	testCases := []struct {
		userHomeDirErr bool
		getWdErr       bool
	}{
		{},
		{userHomeDirErr: true},
		// To get to the getWd call, you have to get a userHomeDirErr first
		{userHomeDirErr: true, getWdErr: true},
	}

	for _, tc := range testCases {
		if tc.userHomeDirErr {
			userHomeDir = func() (string, error) { return "", errors.New("userHomeDirErr") }
		} else {
			userHomeDir = func() (string, error) { return "testing-homedir", nil }
		}

		if tc.getWdErr {
			getWd = func() (dir string, err error) { return "", errors.New("getWdErr") }
		} else {
			getWd = func() (dir string, err error) { return "testing-wd", nil }
		}

		determineHomeDir()

		// These need to be done in a specific order
		// Check for getWdErr first as userHomeDirErr must be set to true as well
		want := "testing-homedir"
		if tc.getWdErr {
			want = ""
		} else if tc.userHomeDirErr {
			want = "testing-wd"
		}

		if homeDir != want {
			t.Errorf("incorrect homeDir: '%s', want: '%s'", homeDir, want)
		}
	}
}
