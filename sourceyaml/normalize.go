/*
Copyright © 2025 Brian Curnow

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package sourceyaml

import (
	"fmt"
	"net"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const nameRegex = `^[a-zA-Z0-9._-]+$`

var (
	allowedClasses = map[string]bool{
		"IN": true,
		"CS": true,
		"CH": true,
		"HS": true,
	}
)

// This function reviews various types of records to ensure that there are no errors that the YAML parser can't catch
// For example, not ending a value with a '.', using an IP vs a name in a CNAME record, etc.
func Normalize(zones map[string]*Zone) error {
	for name, zone := range zones {
		// Check the zone name
		err := isFullyQualified(name)
		if err != nil {
			return fmt.Errorf("Invalid zone name %s: %w", name, err)
		}

		// Check the remaining fields in the zone
		err = NormalizeZone(zone)
		if err != nil {
			return fmt.Errorf("Failed to normalize zone %s: %w", name, err)
		}

		zones[name] = zone
	}

	return nil
}

// Performs all the normalization on a single zone
func NormalizeZone(zone *Zone) error {
	// Check the SOA record fields
	if err := normalizeSOAData(zone); err != nil {
		return err
	}

	// Normalize the serial number
	if zone.GenerateSerial {
		index := zone.SerialChangeIndex
		if zone.SerialChangeIndex == 0 {
			index = uint32(1)
		}

		serial, err := generateSerial(index)
		if err != nil {
			return fmt.Errorf("Unable to generate serial: %w", err)
		}
		zone.Serial = serial
	} else {
		if zone.Serial == 0 {
			return fmt.Errorf("serial number must be provided when generate_serial is false")
		}
	}
	// Normalize the Resource Records
	for name, record := range zone.ResourceRecords {
		if err := normlizeResourceRecord(name, &record, zone.ResourceRecords); err != nil {
			return err
		}
	}

	return nil
}

func normalizeSOAData(zone *Zone) error {
	if zone.Class == "" {
		zone.Class = "IN" // Default class is IN
	} else if !isValidClass(zone.Class) {
		return fmt.Errorf("Invalid SOA class: %s", zone.Class)
	}

	err := isFullyQualified(zone.Nameserver)
	if err != nil {
		return fmt.Errorf("Invalid nameserver in SOA: %s: %w", zone.Nameserver, err)
	}

	// Check if the administrator is valid
	// We require an email but the format requires some special formatting
	formattedEmail, err := formatAdministratorEmail(zone.Administrator)
	if err != nil {
		return fmt.Errorf("Failed to format administrator email: %w", err)
	}
	zone.Administrator = formattedEmail

	// Check the values of the SOA fields
	if zone.Refresh.Value <= 0 || zone.Retry.Value <= 0 || zone.Expire.Value <= 0 || zone.Minimum.Value < 0 {
		return fmt.Errorf("Invalid SOA values for zone, Refresh: %d, Retry: %d, Expire: %d, Minimum: %d. Values must be greater than or equal to zero (0)",
			zone.Refresh.Value, zone.Retry.Value, zone.Expire.Value, zone.Minimum.Value)

	}

	return nil
}

// Handles all the different resource record types
func normlizeResourceRecord(name string, record *ResourceRecord, resourceRecords map[string]ResourceRecord) error {
	switch record.Type {
	case "NS":
		return normalizeNSRecord(name, record)
	case "A":
		return normalizeARecord(name, record)
	case "CNAME":
		return normalizeCNAMERecord(name, record, resourceRecords)
	default:
		return fmt.Errorf("Unsupported resource record type: %s", record.Type)
	}
}

func normalizeNSRecord(nameserver string, record *ResourceRecord) error {
	// Check the nameserver
	err := isFullyQualified(nameserver)
	if err != nil {
		return fmt.Errorf("Invalid nameserver %s: %w", nameserver, err)
	}

	// Check if the value is a valid name (not an IP address)
	if net.ParseIP(nameserver) != nil {
		return fmt.Errorf("NS record value cannot be an IP address: %s", nameserver)
	}

	// Check if the value is valid, in this case, the value" is the domain because this is optional in an NS record
	if record.Value != "" && !isValidNameOrWildcard(record.Value) {
		return fmt.Errorf("Invalid NS record value: %s", record.Value)
	}

	// Check if the class is valid
	if !isValidClass(record.Class) {
		return fmt.Errorf("Invalid NS record class: %s", record.Class)
	}

	return nil
}

func normalizeARecord(name string, record *ResourceRecord) error {
	// Check if the domain is valid
	if !isValidNameOrWildcard(name) {
		return fmt.Errorf("Invalid A record name: %s", name)
	}

	if !isValidClass(record.Class) {
		return fmt.Errorf("Invalid A record class: %s", record.Class)
	}

	// Check if the value is a valid IP address
	if net.ParseIP(record.Value) == nil {
		return fmt.Errorf("Invalid A record value: %s, must be a valid IP address", record.Value)
	}

	return nil
}

func normalizeCNAMERecord(name string, record *ResourceRecord, resourceRecords map[string]ResourceRecord) error {
	// Check if the name is valid
	if !isValidName(name) {
		return fmt.Errorf("Invalid CNAME record name: %s", name)
	}

	// Check if the value is a valid name (not an IP address)
	if net.ParseIP(record.Value) != nil {
		return fmt.Errorf("CNAME record value cannot be an IP address: %s", record.Value)
	}

	// CNAME values must also point to a valid name within the zone
	found, valid := resourceRecords[record.Value]
	if !valid || found.Type != "A" {
		return fmt.Errorf("CNAME record value: %s must point to a valid A record within the zone", record.Value)
	}

	return nil
}

func isValidClass(class string) bool {
	if class == "" {
		// Class is always optional
		return true
	}

	// Check if the class is one of the allowed classes
	_, valid := allowedClasses[class]
	return valid
}

func isValidName(name string) bool {
	// Check if the name matches the regex
	matched, _ := regexp.MatchString(nameRegex, name)
	return matched
}

func isValidNameOrWildcard(name string) bool {
	// Check if the name matches the regex or is a wildcard
	if name == "@" {
		return true
	}
	return isValidName(name)
}

func formatAdministratorEmail(email string) (string, error) {
	// Parse the email address
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("Invalid email address: %s", email)
	}

	return ensureTrailingDot(strings.ReplaceAll(email, "@", ".")), nil
}

// Most DNS names in a zone file need to be fully qualified domain names, while we can't validate if the entire name itself is valid, we can ensure that it ends with a trailing dot
func isFullyQualified(name string) error {
	if !isValidName(name) {
		return fmt.Errorf("Invalid characters")
	}

	if !hasTrailingDot(name) {
		return fmt.Errorf("Must end with a trailing dot")
	}

	// Count the full stops ('.') in the name, there must be at least two (one for the root and one for the domain)
	if strings.Count(name, ".") < 2 {
		return fmt.Errorf("Must be fully qualified with at least two dots")
	}
	return nil
}

func hasTrailingDot(name string) bool {
	return len(name) > 0 && name[len(name)-1] == '.'
}

// Ensure that the string passed in ends with a trailing dot
func ensureTrailingDot(name string) string {
	if !hasTrailingDot(name) {
		return name + "."
	}
	return name
}

func generateSerial(index uint32) (uint32, error) {
	t := time.Now()
	serialString := fmt.Sprintf("%02d%02d%04d%02d", t.Day(), t.Month(), t.Year(), index)

	parsedSerial, err := strconv.ParseUint(serialString, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Unable to generate a serial number from day: %d, month: %d, year: %d, changeIndex: %d: %w", t.Day(), t.Month(), t.Year(), index, err)
	}

	// Explicitly convert to a uint32
	return uint32(parsedSerial), nil
}
