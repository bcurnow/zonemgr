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

import "github.com/bcurnow/zonemgr/schema"

// Defines the types of plugins which can be support
type PluginType string

const (
	A          PluginType = PluginType(schema.A)
	A6         PluginType = PluginType(schema.A6)
	AAAA       PluginType = PluginType(schema.AAAA)
	AFSDB      PluginType = PluginType(schema.AFSDB)
	ALIAS      PluginType = PluginType(schema.ALIAS)
	ANAME      PluginType = PluginType(schema.ANAME)
	APL        PluginType = PluginType(schema.APL)
	ATMA       PluginType = PluginType(schema.ATMA)
	AXFR       PluginType = PluginType(schema.AXFR)
	CAA        PluginType = PluginType(schema.CAA)
	CDNSKEY    PluginType = PluginType(schema.CDNSKEY)
	CDS        PluginType = PluginType(schema.CDS)
	CERT       PluginType = PluginType(schema.CERT)
	CNAME      PluginType = PluginType(schema.CNAME)
	CSYNC      PluginType = PluginType(schema.CSYNC)
	DHCID      PluginType = PluginType(schema.DHCID)
	DLV        PluginType = PluginType(schema.DLV)
	DNAME      PluginType = PluginType(schema.DNAME)
	DNSKEY     PluginType = PluginType(schema.DNSKEY)
	DOA        PluginType = PluginType(schema.DOA)
	DS         PluginType = PluginType(schema.DS)
	EID        PluginType = PluginType(schema.EID)
	EUI48      PluginType = PluginType(schema.EUI48)
	EUI64      PluginType = PluginType(schema.EUI64)
	GID        PluginType = PluginType(schema.GID)
	GPOS       PluginType = PluginType(schema.GPOS)
	HINFO      PluginType = PluginType(schema.HINFO)
	HIP        PluginType = PluginType(schema.HIP)
	HTTPS      PluginType = PluginType(schema.HTTPS)
	IPSECKEY   PluginType = PluginType(schema.IPSECKEY)
	ISDN       PluginType = PluginType(schema.ISDN)
	IXFR       PluginType = PluginType(schema.IXFR)
	KEY        PluginType = PluginType(schema.KEY)
	KX         PluginType = PluginType(schema.KX)
	L32        PluginType = PluginType(schema.L32)
	L64        PluginType = PluginType(schema.L64)
	LOC        PluginType = PluginType(schema.LOC)
	LP         PluginType = PluginType(schema.LP)
	MAILA      PluginType = PluginType(schema.MAILA)
	MAILB      PluginType = PluginType(schema.MAILB)
	MB         PluginType = PluginType(schema.MB)
	MD         PluginType = PluginType(schema.MD)
	MF         PluginType = PluginType(schema.MF)
	MG         PluginType = PluginType(schema.MG)
	MINFO      PluginType = PluginType(schema.MINFO)
	MR         PluginType = PluginType(schema.MR)
	MX         PluginType = PluginType(schema.MX)
	NAPTR      PluginType = PluginType(schema.NAPTR)
	NB         PluginType = PluginType(schema.NB)
	NBSTAT     PluginType = PluginType(schema.NBSTAT)
	NID        PluginType = PluginType(schema.NID)
	NIMLOC     PluginType = PluginType(schema.NIMLOC)
	NINFO      PluginType = PluginType(schema.NINFO)
	NS         PluginType = PluginType(schema.NS)
	NSAP       PluginType = PluginType(schema.NSAP)
	NSAP_PTR   PluginType = PluginType(schema.NSAP_PTR)
	NSEC       PluginType = PluginType(schema.NSEC)
	NSEC3      PluginType = PluginType(schema.NSEC3)
	NSEC3PARAM PluginType = PluginType(schema.NSEC3PARAM)
	NULL       PluginType = PluginType(schema.NULL)
	NXT        PluginType = PluginType(schema.NXT)
	OPENPGPKEY PluginType = PluginType(schema.OPENPGPKEY)
	OPT        PluginType = PluginType(schema.OPT)
	PTR        PluginType = PluginType(schema.PTR)
	PX         PluginType = PluginType(schema.PX)
	RKEY       PluginType = PluginType(schema.RKEY)
	RP         PluginType = PluginType(schema.RP)
	RT         PluginType = PluginType(schema.RT)
	RRSIG      PluginType = PluginType(schema.RRSIG)
	SIG        PluginType = PluginType(schema.SIG)
	SINK       PluginType = PluginType(schema.SINK)
	SMIMEA     PluginType = PluginType(schema.SMIMEA)
	SOA        PluginType = PluginType(schema.SOA)
	SPF        PluginType = PluginType(schema.SPF)
	SRV        PluginType = PluginType(schema.SRV)
	SSHFP      PluginType = PluginType(schema.SSHFP)
	SVCB       PluginType = PluginType(schema.SVCB)
	TA         PluginType = PluginType(schema.TA)
	TALINK     PluginType = PluginType(schema.TALINK)
	TKEY       PluginType = PluginType(schema.TKEY)
	TLSA       PluginType = PluginType(schema.TLSA)
	TSIG       PluginType = PluginType(schema.TSIG)
	TXT        PluginType = PluginType(schema.TXT)
	UID        PluginType = PluginType(schema.UID)
	UINFO      PluginType = PluginType(schema.UINFO)
	UNSPEC     PluginType = PluginType(schema.UNSPEC)
	URI        PluginType = PluginType(schema.URI)
	WKS        PluginType = PluginType(schema.WKS)
	X25        PluginType = PluginType(schema.X25)
	ZONEMD     PluginType = PluginType(schema.ZONEMD)
)
