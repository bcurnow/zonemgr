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

// Defines the types of plugins which can be support
type PluginType string

const (
	A          PluginType = PluginType("A")
	A6         PluginType = PluginType("A6")
	AAAA       PluginType = PluginType("AAAA")
	AFSDB      PluginType = PluginType("AFSDB")
	ALIAS      PluginType = PluginType("ALIAS")
	ANAME      PluginType = PluginType("ANAME")
	APL        PluginType = PluginType("APL")
	ATMA       PluginType = PluginType("ATMA")
	AXFR       PluginType = PluginType("AXFR")
	CAA        PluginType = PluginType("CAA")
	CDNSKEY    PluginType = PluginType("CDNSKEY")
	CDS        PluginType = PluginType("CDS")
	CERT       PluginType = PluginType("CERT")
	CNAME      PluginType = PluginType("CNAME")
	CSYNC      PluginType = PluginType("CSYNC")
	DHCID      PluginType = PluginType("DHCID")
	DLV        PluginType = PluginType("DLV")
	DNAME      PluginType = PluginType("DNAME")
	DNSKEY     PluginType = PluginType("DNSKEY")
	DOA        PluginType = PluginType("DOA")
	DS         PluginType = PluginType("DS")
	EID        PluginType = PluginType("EID")
	EUI48      PluginType = PluginType("EUI48")
	EUI64      PluginType = PluginType("EUI64")
	GID        PluginType = PluginType("GID")
	GPOS       PluginType = PluginType("GPOS")
	HINFO      PluginType = PluginType("HINFO")
	HIP        PluginType = PluginType("HIP")
	HTTPS      PluginType = PluginType("HTTPS")
	IPSECKEY   PluginType = PluginType("IPSECKEY")
	ISDN       PluginType = PluginType("ISDN")
	IXFR       PluginType = PluginType("IXFR")
	KEY        PluginType = PluginType("KEY")
	KX         PluginType = PluginType("KX")
	L32        PluginType = PluginType("L32")
	L64        PluginType = PluginType("L64")
	LOC        PluginType = PluginType("LOC")
	LP         PluginType = PluginType("LP")
	MAILA      PluginType = PluginType("MAILA")
	MAILB      PluginType = PluginType("MAILB")
	MB         PluginType = PluginType("MB")
	MD         PluginType = PluginType("MD")
	MF         PluginType = PluginType("MF")
	MG         PluginType = PluginType("MG")
	MINFO      PluginType = PluginType("MINFO")
	MR         PluginType = PluginType("MR")
	MX         PluginType = PluginType("MX")
	NAPTR      PluginType = PluginType("NAPTR")
	NB         PluginType = PluginType("NB")
	NBSTAT     PluginType = PluginType("NBSTAT")
	NID        PluginType = PluginType("NID")
	NIMLOC     PluginType = PluginType("NIMLOC")
	NINFO      PluginType = PluginType("NINFO")
	NS         PluginType = PluginType("NS")
	NSAP       PluginType = PluginType("NSAP")
	NSAP_PTR   PluginType = PluginType("NSAP_PTR")
	NSEC       PluginType = PluginType("NSEC")
	NSEC3      PluginType = PluginType("NSEC3")
	NSEC3PARAM PluginType = PluginType("NSEC3PARAM")
	NULL       PluginType = PluginType("NULL")
	NXT        PluginType = PluginType("NXT")
	OPENPGPKEY PluginType = PluginType("OPENPGPKEY")
	OPT        PluginType = PluginType("OPT")
	PTR        PluginType = PluginType("PTR")
	PX         PluginType = PluginType("PX")
	RKEY       PluginType = PluginType("RKEY")
	RP         PluginType = PluginType("RP")
	RT         PluginType = PluginType("RT")
	RRSIG      PluginType = PluginType("RRSIG")
	SIG        PluginType = PluginType("SIG")
	SINK       PluginType = PluginType("SINK")
	SMIMEA     PluginType = PluginType("SMIMEA")
	SOA        PluginType = PluginType("SOA")
	SPF        PluginType = PluginType("SPF")
	SRV        PluginType = PluginType("SRV")
	SSHFP      PluginType = PluginType("SSHFP")
	SVCB       PluginType = PluginType("SVCB")
	TA         PluginType = PluginType("TA")
	TALINK     PluginType = PluginType("TALINK")
	TKEY       PluginType = PluginType("TKEY")
	TLSA       PluginType = PluginType("TLSA")
	TSIG       PluginType = PluginType("TSIG")
	TXT        PluginType = PluginType("TXT")
	UID        PluginType = PluginType("UID")
	UINFO      PluginType = PluginType("UINFO")
	UNSPEC     PluginType = PluginType("UNSPEC")
	URI        PluginType = PluginType("URI")
	WKS        PluginType = PluginType("WKS")
	X25        PluginType = PluginType("X25")
	ZONEMD     PluginType = PluginType("ZONEMD")
)
