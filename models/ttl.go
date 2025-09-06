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

package models

import (
	"fmt"
)

type TTL struct {
	Value   *int32 `yaml:"value"` // The use of a pointer to an int32 allows us to handle missing (nil) values more easily
	Comment string `yaml:"comment"`
}

func (ttl *TTL) String() string {
	return fmt.Sprintf("TTL{ Value: %s, Comment: %s }", int32ToString(ttl.Value), ttl.Comment)
}

func (t *TTL) Render() string {
	if t.Value != nil {
		comment := t.Comment
		if comment != "" {
			comment = " ;" + comment
		}
		return fmt.Sprintf("$TTL %d%s", *t.Value, comment)
	}
	return ""
}
