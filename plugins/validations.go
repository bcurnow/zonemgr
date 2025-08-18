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
package plugins

import (
	"fmt"
	"net/mail"
	"regexp"
	"slices"
	"strings"

	"github.com/bcurnow/zonemgr/schema"
)

// This regex is based on RFC1035 and allows for:
//   - Label: Up to 63 characters a-z, A-Z, 0-9, and hyphen which does not start (?!-) or end (?<=\w) with a hypen
//   - Domain: Any number of sub-domain entries separate by "." which follow the same rules as the label
const dnsNameRegexRFC1035String = `^([A-Za-z])([A-Za-z0-9-]{1,62})(\.[A-Za-z0-9-]{1,63})*\.{0,1}$`

var dnsNameRegexRFC1035 = regexp.MustCompile(dnsNameRegexRFC1035String)

// Performs the standard validations for resource records
// This includes:
//   - Validation that the resource record is of the specified type - This is not case insensitive but the the type will be normalized to uppercase
//   - Validation of the class - An empty class will be considered valid, any defaulting or enforcement beyond that is the responsiblity of the individual plugins
//   - Validation that only Value or Values is populated
//   - Validation that only Comment or Values is populated
func StandardValidations(identifier string, rr *schema.ResourceRecord, supportedTypes []PluginType) error {
	// Validate that this resource record is of the supported type
	if err := IsSupportedPluginType(identifier, rr.Type, supportedTypes); err != nil {
		return err
	}

	// Validate the class
	if !rr.Class.IsValid() {
		return fmt.Errorf("invalid %s record, '%s' is not a valid class, identifier: '%s'", rr.Type, rr.Class, identifier)
	}

	// Validate that there is only one Value set
	if err := rr.IsValueSetInOnePlace(identifier); err != nil {
		return err
	}

	// Validate that there is only one Comment set
	if err := rr.IsCommentSetInOnePlace(identifier); err != nil {
		return err
	}

	return nil
}

// Checks if the supplied resource record matches one of the support plugin types
func IsSupportedPluginType(identifier string, rrType schema.ResourceRecordType, supportedTypes []PluginType) error {
	if !slices.Contains(supportedTypes, PluginType(rrType)) {
		return fmt.Errorf("this plugin does not handle resource records of type '%s' only '%s', identifier: '%s'", rrType, supportedTypes, identifier)
	}
	return nil
}

// Validates that the name provided matches the RFC1035 regex for valid names according to RFC1035
// and is less then or equal to 255 total characters
func IsValidRFC1035Name(identifier string, name string, rrType schema.ResourceRecordType) error {
	if len(name) > 255 {
		return fmt.Errorf("invalid %s record, must be less than 255 characters: '%s', identifier=%s", rrType, name, identifier)
	}

	if !dnsNameRegexRFC1035.MatchString(name) {
		return fmt.Errorf("invalid %s record, does not match regexp '%s': '%s', identifier=%s", rrType, dnsNameRegexRFC1035String, name, identifier)
	}

	//Split the domain at each part (".") and then run some additional validations
	parts := strings.Split(name, ".")
	for _, part := range parts {
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return fmt.Errorf("invalid %s record, can not start or end with a hyphen (-): '%s', identifier=%s", rrType, name, identifier)
		}
	}
	return nil
}

// Checks if the name provide is either the wildcard ('@') or is a valid name
func IsValidNameOrWildcard(identifier string, name string, rrType schema.ResourceRecordType) error {
	// Check if the name matches the regex or is a wildcard
	if name == "@" {
		return nil
	}
	return IsValidRFC1035Name(identifier, name, rrType)
}

// Formats and email address according to RFC1035
func FormatEmail(identifier string, email string, rrType schema.ResourceRecordType) (string, error) {
	if strings.Contains(email, "@") {
		// Assume this is a standard email address that will be parseable
		address, err := mail.ParseAddress(email)
		if err != nil {
			return "", fmt.Errorf("invalid %s record, invalid email address: %s, identifier=%s", rrType, email, identifier)
		}

		// Get the username portion of the email (<username>@<domain>), keep in mind that valid usernames can continue '@'
		emailSeperator := strings.LastIndex(address.Address, "@")
		username := email[:emailSeperator]
		domain := email[emailSeperator+1:]

		// Escape any dots (.) in the username as these are special characters in a zonefile
		username = strings.ReplaceAll(username, ".", "\\.")

		// Recombinee the user and domain with a dot (.) to conform to RFC1035
		email = username + "." + domain
		// Replace the @ with a dot to follow RFC
	}

	// At this point, assume that email address is a properly formatted RFC1035 string, there's only so much we can do to parse at this point
	return EnsureTrailingDot(email), nil
}

// Most DNS names in a zone file need to be fully qualified domain names, while we can't validate if the entire name itself is valid,
// we can ensure that it is a valid name and ends with a trailing dot
func IsFullyQualified(identifier string, name string, rrType schema.ResourceRecordType) error {
	if err := IsValidRFC1035Name(identifier, name, rrType); err != nil {
		return err
	}

	if !hasTrailingDot(name) {
		return fmt.Errorf("invalid %s record, must end with a trailing dot: '%s', identifier=%s", rrType, name, identifier)
	}

	// Count the full stops ('.') in the name, there must be at least two (one for the root and one for the domain)
	if strings.Count(name, ".") < 2 {
		return fmt.Errorf("invalid %s record, must be fully qualified with at least two dots: '%s', identifier=%s", rrType, name, identifier)
	}
	return nil
}

// Ensure that the string passed in ends with a trailing dot
func EnsureTrailingDot(name string) string {
	if !hasTrailingDot(name) {
		return name + "."
	}
	return name
}

func hasTrailingDot(name string) bool {
	return len(name) > 0 && name[len(name)-1] == '.'
}
