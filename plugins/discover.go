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

package plugins

import (
	"os"
	"path/filepath"
)

// Walks the specified directory looking for plugins, returns an array of all the executables found
// Until goplugin.Discover is updated to check for the executable bit, this is our own implementation
func discoverPlugins(dir string) (map[string]string, error) {
	var executables = make(map[string]string)

	logger.Trace("Walking plugins dir", "dir", dir)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Don't traverse sub-directories, this is arbitrary but we are keeping it simple
		if d.IsDir() && path != dir {
			logger.Info("Subdirectories are not supported", "Subdirectory", path)
			return filepath.SkipDir
		}

		// Because we're using WalkDir, we need to get the FileInfo from the DirEntry
		info, err := d.Info()
		if err != nil {
			return err
		}

		// Check if this is a file and if the file is executable
		if info.Mode().IsRegular() {
			// 0111 checks for the execute bit to be set
			if info.Mode()&0111 == 0 {
				logger.Info("Skipping non-executable file", "File", path)
				return nil
			}

			// Get the absolute path of the file so we can provide the best debugging information
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			executables[path] = absPath
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return executables, nil
}
