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
	"github.com/bcurnow/zonemgr/schema"
	goplugin "github.com/hashicorp/go-plugin"
)

// Defines the types of plugins which can be support
type PluginType string

const (
	// Define the record handling plugin types in terms of the record type
	RecordA     PluginType = PluginType(schema.A)
	RecordCNAME PluginType = PluginType(schema.CNAME)
	RecordHINFO PluginType = PluginType(schema.HINFO)
	RecordMB    PluginType = PluginType(schema.MB)
	RecordMD    PluginType = PluginType(schema.MD)
	RecordMF    PluginType = PluginType(schema.MF)
	RecordMG    PluginType = PluginType(schema.MG)
	RecordMINFO PluginType = PluginType(schema.MINFO)
	RecordMR    PluginType = PluginType(schema.MR)
	RecordMX    PluginType = PluginType(schema.MX)
	RecordNS    PluginType = PluginType(schema.NS)
	RecordNULL  PluginType = PluginType(schema.NULL)
	RecordPTR   PluginType = PluginType(schema.PTR)
	RecordSOA   PluginType = PluginType(schema.SOA)
	RecordTXT   PluginType = PluginType(schema.TXT)
	RecordWKS   PluginType = PluginType(schema.WKS)
)

// This is the go-plugin handshake information that needs to be used for all plugins
var HandshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ZONEMGR_PLUGIN",
	MagicCookieValue: "BEA0CA21-AAC6-4EA8-BB29-4B6B2E39B1AE",
}
