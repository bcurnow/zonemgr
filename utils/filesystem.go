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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
)

// This package wraps standard methods from the filepath and os packages
// It is intended to isolate the logic to a single place and to make testing easier

var (
	// Store methods from the filepath and os packages as variables
	// This will allow us to override these methods when we are testing
	abs     = filepath.Abs
	walkDir = filepath.WalkDir
	create  = os.Create
	chmod   = os.Chmod

	// Make sure we implement the interface
	_ FileSystem = &fileSystem{}
)

type FileSystem interface {
	// Creates the path specified, sets the mode and calls the contentFn to generate the file content
	CreateFile(path string, mode os.FileMode, contentFn func() ([]byte, error)) error
	// Takes a path name and returns the absolute path value
	// This method is similar to filepath.Abs but also handles
	// paths that start with '~' and will automatically expand this to the
	// home directory of the current user
	ToAbsoluteFilePath(path string) (string, error)
	// Walks the specified path and returns all executable files found
	// If includeSubDirs is true, will find all executables in all subDirs as well
	WalkExecutables(root string, includeSubDirs bool) (map[string]string, error)
}

type fileSystem struct {
}

func FS() FileSystem {
	return &fileSystem{}
}

func (fs *fileSystem) CreateFile(path string, mode os.FileMode, contentFn func() ([]byte, error)) error {
	outputFile, err := create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", path, err)
	}
	defer outputFile.Close()

	// Set the mode on the file
	if err := chmod(path, mode); err != nil {
		return fmt.Errorf("error with chmod of '%s' to '%o'", path, mode)
	}

	content, err := contentFn()
	if err != nil {
		return fmt.Errorf("error generating content for output file '%s': %w", path, err)
	}

	bytesWritten, err := outputFile.Write(content)
	if err != nil {
		return fmt.Errorf("error writing content for output file '%s': %w", path, err)
	} else {
		hclog.L().Trace("Wrote content to output file", "outputFile", path, "bytesWritten", bytesWritten)
	}

	return nil
}

func (fs *fileSystem) ToAbsoluteFilePath(path string) (string, error) {
	//go doesn't automatically handle the ~ expansion, do this manually
	if strings.HasPrefix(path, "~") {
		path = filepath.Join(HomeDir, path[1:])
	}

	absPath, err := abs(path)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Could not convert '%s' into an absolute path", path))
		return "", err
	}
	return absPath, nil
}

func (fs *fileSystem) WalkExecutables(root string, includeSubDirs bool) (map[string]string, error) {
	executables := make(map[string]string)
	err := walkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				hclog.L().Trace("Could not walk path", "path", path)
				return nil
			}
			return walkErr
		}
		hclog.L().Trace("Processing path", "path", path)

		// Don't traverse sub-directories, this is arbitrary but we are keeping it simple
		if d.IsDir() && path != root && !includeSubDirs {
			hclog.L().Trace("Subdirectories are not supported, skipping", "path", path)
			return filepath.SkipDir
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		// Check if this is a file and if the file is executable
		if info.Mode().IsRegular() {
			// 0111 checks for the execute bit to be set
			if info.Mode()&0111 == 0 {
				hclog.L().Trace("Skipping non-executable file", "path", path)
				return nil
			}

			// Get the absolute path of the file so we can provide the best debugging information
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			hclog.L().Trace("Adding executable", "executable", absPath)
			executables[filepath.Base(path)] = absPath
		}
		return nil
	})
	return executables, err
}
