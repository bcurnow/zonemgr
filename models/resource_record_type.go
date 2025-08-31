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

// Defines the t ypes of plugins which can be support
type ResourceRecordType string

// Source: https://en.wikipedia.org/wiki/List_of_DNS_record_types
const (
	A          ResourceRecordType = "A"
	A6         ResourceRecordType = "A6"
	AAAA       ResourceRecordType = "AAAA"
	AFSDB      ResourceRecordType = "AFSDB"
	ALIAS      ResourceRecordType = "ALIAS"
	ANAME      ResourceRecordType = "ANAME"
	APL        ResourceRecordType = "APL"
	ATMA       ResourceRecordType = "ATMA"
	AXFR       ResourceRecordType = "AXFR"
	CAA        ResourceRecordType = "CAA"
	CDNSKEY    ResourceRecordType = "CDNSKEY"
	CDS        ResourceRecordType = "CDS"
	CERT       ResourceRecordType = "CERT"
	CNAME      ResourceRecordType = "CNAME"
	CSYNC      ResourceRecordType = "CSYNC"
	DHCID      ResourceRecordType = "DHCID"
	DLV        ResourceRecordType = "DLV"
	DNAME      ResourceRecordType = "DNAME"
	DNSKEY     ResourceRecordType = "DNSKEY"
	DOA        ResourceRecordType = "DOA"
	DS         ResourceRecordType = "DS"
	EID        ResourceRecordType = "EID"
	EUI48      ResourceRecordType = "EUI48"
	EUI64      ResourceRecordType = "EUI64"
	GID        ResourceRecordType = "GID"
	GPOS       ResourceRecordType = "GPOS"
	HINFO      ResourceRecordType = "HINFO"
	HIP        ResourceRecordType = "HIP"
	HTTPS      ResourceRecordType = "HTTPS"
	IPSECKEY   ResourceRecordType = "IPSECKEY"
	ISDN       ResourceRecordType = "ISDN"
	IXFR       ResourceRecordType = "IXFR"
	KEY        ResourceRecordType = "KEY"
	KX         ResourceRecordType = "KX"
	L32        ResourceRecordType = "L32"
	L64        ResourceRecordType = "L64"
	LOC        ResourceRecordType = "LOC"
	LP         ResourceRecordType = "LP"
	MAILA      ResourceRecordType = "MAILA"
	MAILB      ResourceRecordType = "MAILB"
	MB         ResourceRecordType = "MB"
	MD         ResourceRecordType = "MD"
	MF         ResourceRecordType = "MF"
	MG         ResourceRecordType = "MG"
	MINFO      ResourceRecordType = "MINFO"
	MR         ResourceRecordType = "MR"
	MX         ResourceRecordType = "MX"
	NAPTR      ResourceRecordType = "NAPTR"
	NB         ResourceRecordType = "NB"
	NBSTAT     ResourceRecordType = "NBSTAT"
	NID        ResourceRecordType = "NID"
	NIMLOC     ResourceRecordType = "NIMLOC"
	NINFO      ResourceRecordType = "NINFO"
	NS         ResourceRecordType = "NS"
	NSAP       ResourceRecordType = "NSAP"
	NSAP_PTR   ResourceRecordType = "NSAP-PTR"
	NSEC       ResourceRecordType = "NSEC"
	NSEC3      ResourceRecordType = "NSEC3"
	NSEC3PARAM ResourceRecordType = "NSEC3PARAM"
	NULL       ResourceRecordType = "NULL"
	NXT        ResourceRecordType = "NXT"
	OPENPGPKEY ResourceRecordType = "OPENPGPKEY"
	OPT        ResourceRecordType = "OPT"
	PTR        ResourceRecordType = "PTR"
	PX         ResourceRecordType = "PX"
	RKEY       ResourceRecordType = "RKEY"
	RP         ResourceRecordType = "RP"
	RT         ResourceRecordType = "RT"
	RRSIG      ResourceRecordType = "RRSIG"
	SIG        ResourceRecordType = "SIG"
	SINK       ResourceRecordType = "SINK"
	SMIMEA     ResourceRecordType = "SMIMEA"
	SOA        ResourceRecordType = "SOA"
	SPF        ResourceRecordType = "SPF"
	SRV        ResourceRecordType = "SRV"
	SSHFP      ResourceRecordType = "SSHFP"
	SVCB       ResourceRecordType = "SVCB"
	TA         ResourceRecordType = "TA"
	TALINK     ResourceRecordType = "TALINK"
	TKEY       ResourceRecordType = "TKEY"
	TLSA       ResourceRecordType = "TLSA"
	TSIG       ResourceRecordType = "TSIG"
	TXT        ResourceRecordType = "TXT"
	UID        ResourceRecordType = "UID"
	UINFO      ResourceRecordType = "UINFO"
	UNSPEC     ResourceRecordType = "UNSPEC"
	URI        ResourceRecordType = "URI"
	WKS        ResourceRecordType = "WKS"
	X25        ResourceRecordType = "X25"
	ZONEMD     ResourceRecordType = "ZONEMD"
)
