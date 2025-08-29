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
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
)

func ToAbsoluteFilePath(path string, name string) (string, error) {
	//go doesn't automatically handle the ~ expansion, do this manually
	if strings.HasPrefix(path, "~") {
		path = filepath.Join(HomeDir, path[1:])
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Could not convert %s value '%s' into an absolute path", name, path))
		return "", err
	}
	return absPath, nil
}
