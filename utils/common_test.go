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
	"os"
	"testing"
)

var (
	testDir  string
	testFile *os.File
)

func createTemp(t *testing.T) {
	// Create a temp directory
	dir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Errorf("error creating temp directory name: %s", err)
	}
	testDir = dir

	file, err := os.CreateTemp(testDir, t.Name())
	if err != nil {
		t.Errorf("error creating temp file in '%s': %s", testDir, err)
	}
	testFile = file
}

func tempTeardown(_ *testing.T) {
	defer testFile.Close()
	defer os.RemoveAll(testDir)
}
