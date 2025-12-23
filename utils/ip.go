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
	"net/netip"
	"slices"
	"strings"
)

// This type is wrapper around a netit/Addr which provides some additional functionality.
type IP struct {
	// The underlying net.ip/Addr.
	ip netip.Addr
}

func ParseIP(ip string) (IP, error) {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return IP{}, err
	}
	return IP{ip: ipAddr}, nil
}

func (i IP) ReverseZoneName() string {
	if i.ip.Is4() {
		octets := i.ip.AsSlice()
		// Reverse zone names for ipv4 are the first three octets in reverse order.
		return fmt.Sprintf("%d.%d.%d.in-addr.arpa", octets[3], octets[2], octets[1])
	}

	// Reverse zone names for ipv6 are the first 64 bytes reversed and with a dot in between each value. For example:
	// AAAA Record Value: fdda:5cc1:23:4::1f
	// Without zero compression: fdda:5cc1:0023:0004:0000:0000:0000:001f
	// Reversed: f100:0000:0000:0000:4000:3200:1cc5:addf
	// Dotted notation: f.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.0.0.0.3.2.0.0.1.c.c.5.a.d.d.f
	// Reverse Zone Name: f.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa
	// PTR Record Value: 4.0.0.0.3.2.0.0.1.c.c.5.a.d.d.f
	return i.toReverseDottedNotation()[0:32] + "ip6.arpa"
}

func (i IP) PTRRecordValue() string {
	if i.ip.Is4() {
		octets := i.ip.AsSlice()
		// PTR record values for ipv4 are the last octet.
		return fmt.Sprintf("%d", octets[3])
	}

	// Reverse zone names for ipv6 are the first 64 bytes reversed and with a dot in between each value. For example:
	// AAAA Record Value: fdda:5cc1:23:4::1f
	// Without zero compression: fdda:5cc1:0023:0004:0000:0000:0000:001f
	// Reversed: f100:0000:0000:0000:4000:3200:1cc5:addf
	// Dotted notation: f.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.0.0.0.3.2.0.0.1.c.c.5.a.d.d.f
	// Reverse Zone Name: f.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa
	// PTR Record Value: 4.0.0.0.3.2.0.0.1.c.c.5.a.d.d.f
	return i.toReverseDottedNotation()[32:]
}

func (i IP) StringExpanded() string {
	return i.ip.StringExpanded()
}

func (i IP) String() string {
	return i.ip.String()
}

func (i IP) Is4() bool {
	return i.ip.Is4()
}

func (i IP) Is6() bool {
	return i.ip.Is6()
}

func (i IP) toReverseDottedNotation() string {
	expandedString := i.ip.StringExpanded()
	// Remove all colons from the string.
	expandedString = strings.ReplaceAll(expandedString, ":", "")
	// Split the string into a slice of strings, where each string is a single character.
	chars := strings.Split(expandedString, "")
	// Reverse the slice.
	slices.Reverse(chars)
	// Join the slice back into a string with dots in between each character.
	return strings.Join(chars, ".")
}
