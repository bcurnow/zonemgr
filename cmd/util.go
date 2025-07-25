/*
Copyright © 2025 Brian Curnow

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func toAbsoluteFilePath(path string, name string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve %s %s: %v\n", name, path, err)
		os.Exit(1)
	}
	return absPath
}

func templateContent(path string, name string, defaultContent string) string {
	if path != "" {
		path = toAbsoluteFilePath(path, name)
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to read %s %s: %v\n", name, path, err)
			os.Exit(1)
		}

		return string(contentBytes)
	}
	return defaultContent
}
