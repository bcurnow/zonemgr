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
)

var HomeDir string

// This will exit the entire program if we can't get this but this generaly shouldn't happen
// unless we're being run in a very strange way
func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to retrieve current user's home directory: %s\n", err)
		workingDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not determine current working directory, returning empty string")
			homeDir = ""
		} else {
			fmt.Fprintf(os.Stderr, "Defaulting to current working directory: %s\n", workingDir)
			homeDir = workingDir
		}
	}

	HomeDir = homeDir
}
