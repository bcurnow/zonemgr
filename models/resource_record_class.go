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

// Defines the types of classes available in a zone file
type ResourceRecordClass string

const (
	INTERNET ResourceRecordClass = "IN"
	CSNET    ResourceRecordClass = "CS"
	CHAOS    ResourceRecordClass = "CH"
	HESIOD   ResourceRecordClass = "HS"
)

func (rrc ResourceRecordClass) IsValid() bool {
	switch rrc {

	case INTERNET, CSNET, CHAOS, HESIOD, "": //It's always valid for the class to be empty
		return true
	default:
		return false
	}
}
