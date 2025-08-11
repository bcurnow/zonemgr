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

package schema

// Defines the types of plugins which can be support
type ResourceRecordType string

const (
	A     ResourceRecordType = "A"
	CNAME ResourceRecordType = "CNAME"
	HINFO ResourceRecordType = "HINFO"
	MB    ResourceRecordType = "MB"
	MD    ResourceRecordType = "MD"
	MF    ResourceRecordType = "MF"
	MG    ResourceRecordType = "MG"
	MINFO ResourceRecordType = "MINFO"
	MR    ResourceRecordType = "MR"
	MX    ResourceRecordType = "MX"
	NS    ResourceRecordType = "NS"
	NULL  ResourceRecordType = "NULL"
	PTR   ResourceRecordType = "PTR"
	SOA   ResourceRecordType = "SOA"
	TXT   ResourceRecordType = "TXT"
	WKS   ResourceRecordType = "WKS"
)

const (
	ResourceRecordNameFormatString             = "%-40s"
	ResourceRecordTypeFormatString             = "%-6s "
	ResourceRecordMultivalueIndentFormatString = "%4s"
)
