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
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/hashicorp/go-hclog"
)

// This package wraps standard methods from the filepath and os packages
// It is intended to isolate the logic to a single place and to make testing easier

var (
	// Store methods from the filepath and os packages as variables
	// This will allow us to override these methods when we are testing
	abs      = filepath.Abs
	chmod    = os.Chmod
	create   = os.Create
	homeDir  string
	mkdirAll = os.MkdirAll
	stat     = os.Stat
	walkDir  = filepath.WalkDir

	// Make sure we implement the interface
	_ FileSystem = &fileSystem{}
)

type FileSystem interface {
	// Creates the path specified, sets the mode and calls the contentFn to generate the file content
	CreateFile(path string, mode os.FileMode, contentFn func() ([]byte, error)) error
	// Returns true if the path exists, false otherwise
	Exists(path string) bool
	// Gets a lock on a file
	Flock(path string) (*flock.Flock, error)
	// Returns the current user's home directory or "" if the current user can not be determined (e.g. when go-plugin executes a plugin)
	HomeDir() string
	// Creates the directory and any missing parents using the supplied permissions for any directory it creates, if the directory exists, this is a nop
	MkdirAll(path string, mode os.FileMode) error
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

func (fs *fileSystem) Exists(path string) bool {
	_, err := stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Will try and get the lock for 10 seconds, returns an error if it can't get the lock
func (fs *fileSystem) Flock(path string) (*flock.Flock, error) {
	// Create a file lock, this doesn't lock the file...yet
	fileLock := flock.New(path)
	// Setup a 10 second timer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hclog.L().Trace("Attempting to lock file, will try for 10 seconds", "file", path)
	// Try the lock every half second
	locked, err := fileLock.TryLockContext(ctx, 500*time.Millisecond)
	if err != nil {
		return nil, err
	}

	if locked {
		hclog.L().Trace("Locked file", "file", path)
		return fileLock, nil
	}

	// This should only happen if we fail to lock the file, this can happen because we timeout
	return nil, fmt.Errorf("unexpected error, expected '%s' to be locked", path)
}

func (fs *fileSystem) HomeDir() string {
	return homeDir
}

func (fs *fileSystem) MkdirAll(path string, mode os.FileMode) error {
	return mkdirAll(path, mode)
}

func (fs *fileSystem) ToAbsoluteFilePath(path string) (string, error) {
	//go doesn't automatically handle the ~ expansion, do this manually
	if strings.HasPrefix(path, "~") {
		path = filepath.Join(fs.HomeDir(), path[1:])
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

// This will exit the entire program if we can't get this but this generaly shouldn't happen
// unless we're being run in a very strange way
func init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to retrieve current user's home directory: %s\n", err)
		workingDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not determine current working directory, returning empty string")
			dir = ""
		} else {
			fmt.Fprintf(os.Stderr, "Defaulting to current working directory: %s\n", workingDir)
			dir = workingDir
		}
	}

	homeDir = dir
}
