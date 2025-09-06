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
type Type string

const (
	A          Type = Type("A")
	A6         Type = Type("A6")
	AAAA       Type = Type("AAAA")
	AFSDB      Type = Type("AFSDB")
	ALIAS      Type = Type("ALIAS")
	ANAME      Type = Type("ANAME")
	APL        Type = Type("APL")
	ATMA       Type = Type("ATMA")
	AXFR       Type = Type("AXFR")
	CAA        Type = Type("CAA")
	CDNSKEY    Type = Type("CDNSKEY")
	CDS        Type = Type("CDS")
	CERT       Type = Type("CERT")
	CNAME      Type = Type("CNAME")
	CSYNC      Type = Type("CSYNC")
	DHCID      Type = Type("DHCID")
	DLV        Type = Type("DLV")
	DNAME      Type = Type("DNAME")
	DNSKEY     Type = Type("DNSKEY")
	DOA        Type = Type("DOA")
	DS         Type = Type("DS")
	EID        Type = Type("EID")
	EUI48      Type = Type("EUI48")
	EUI64      Type = Type("EUI64")
	GID        Type = Type("GID")
	GPOS       Type = Type("GPOS")
	HINFO      Type = Type("HINFO")
	HIP        Type = Type("HIP")
	HTTPS      Type = Type("HTTPS")
	IPSECKEY   Type = Type("IPSECKEY")
	ISDN       Type = Type("ISDN")
	IXFR       Type = Type("IXFR")
	KEY        Type = Type("KEY")
	KX         Type = Type("KX")
	L32        Type = Type("L32")
	L64        Type = Type("L64")
	LOC        Type = Type("LOC")
	LP         Type = Type("LP")
	MAILA      Type = Type("MAILA")
	MAILB      Type = Type("MAILB")
	MB         Type = Type("MB")
	MD         Type = Type("MD")
	MF         Type = Type("MF")
	MG         Type = Type("MG")
	MINFO      Type = Type("MINFO")
	MR         Type = Type("MR")
	MX         Type = Type("MX")
	NAPTR      Type = Type("NAPTR")
	NB         Type = Type("NB")
	NBSTAT     Type = Type("NBSTAT")
	NID        Type = Type("NID")
	NIMLOC     Type = Type("NIMLOC")
	NINFO      Type = Type("NINFO")
	NS         Type = Type("NS")
	NSAP       Type = Type("NSAP")
	NSAP_PTR   Type = Type("NSAP_PTR")
	NSEC       Type = Type("NSEC")
	NSEC3      Type = Type("NSEC3")
	NSEC3PARAM Type = Type("NSEC3PARAM")
	NULL       Type = Type("NULL")
	NXT        Type = Type("NXT")
	OPENPGPKEY Type = Type("OPENPGPKEY")
	OPT        Type = Type("OPT")
	PTR        Type = Type("PTR")
	PX         Type = Type("PX")
	RKEY       Type = Type("RKEY")
	RP         Type = Type("RP")
	RT         Type = Type("RT")
	RRSIG      Type = Type("RRSIG")
	SIG        Type = Type("SIG")
	SINK       Type = Type("SINK")
	SMIMEA     Type = Type("SMIMEA")
	SOA        Type = Type("SOA")
	SPF        Type = Type("SPF")
	SRV        Type = Type("SRV")
	SSHFP      Type = Type("SSHFP")
	SVCB       Type = Type("SVCB")
	TA         Type = Type("TA")
	TALINK     Type = Type("TALINK")
	TKEY       Type = Type("TKEY")
	TLSA       Type = Type("TLSA")
	TSIG       Type = Type("TSIG")
	TXT        Type = Type("TXT")
	UID        Type = Type("UID")
	UINFO      Type = Type("UINFO")
	UNSPEC     Type = Type("UNSPEC")
	URI        Type = Type("URI")
	WKS        Type = Type("WKS")
	X25        Type = Type("X25")
	ZONEMD     Type = Type("ZONEMD")
)
