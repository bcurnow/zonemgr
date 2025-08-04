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
const dnsNameRegexRFC1035String = `^([A-Za-z0-9-]{1,63})(\.[A-Za-z0-9-]{1,63})*\.{0,1}$`

var dnsNameRegexRFC1035 = regexp.MustCompile(dnsNameRegexRFC1035String)

// Useful shortcut for plugin functions which return (schema.ResourceRecord, error) when returning an error
func NilResourceRecord() schema.ResourceRecord {
	return schema.ResourceRecord{}
}

// Performs the standard validations for resource records
// This includes:
//   - Validation that the resource record is of the specified type - This is not case insensitive but the the type will be normalized to uppercase
//   - Validation of the class - An empty class will be considered valid, any defaulting or enforcement beyond that is the responsiblity of the individual plugins
//   - Validation that only Value or Values is populated
//   - Validation that only Comment or Values is populated
func StandardValidations(identifier string, rr *schema.ResourceRecord, supportedTypes []PluginType) error {
	NormalizeType(rr)

	// Validate that this resource record is of the supported type
	if err := IsSupportedPluginType(identifier, rr, supportedTypes); err != nil {
		return err
	}

	// Validate the class
	if !IsValidClass(rr) {
		return fmt.Errorf("%s record invalid, '%s' is not a valid class, identifier: '%s'", rr.Type, rr.Class, identifier)
	}

	// Validate that there is only one Value set
	if err := IsValueSetInOnePlace(identifier, rr); err != nil {
		return err
	}

	// Validate that there is only one Comment set
	if err := IsCommentSetInOnePlace(identifier, rr); err != nil {
		return err
	}

	return nil
}

// Ensures that the resource record type is a consistent value
func NormalizeType(rr *schema.ResourceRecord) {
	rr.Type = strings.ToUpper(rr.Type)
}

// Checks if the supplied resource record matches one of the support plugin types
func IsSupportedPluginType(identifier string, rr *schema.ResourceRecord, supportedTypes []PluginType) error {
	if !slices.Contains(supportedTypes, PluginType(rr.Type)) {
		return fmt.Errorf("This plugin does not handle resource records of type '%s' only '%s', identifier: '%s'", rr.Type, supportedTypes, identifier)
	}
	return nil
}

// Validates that the specified class, if set, is a valid ResourceRecordClass
func IsValidClass(rr *schema.ResourceRecord) bool {
	// Class is always optional
	if rr.Class == "" {
		return true
	}

	_, err := schema.ResourceRecordClassFromString(rr.Class)
	if err != nil {
		return false
	}
	return true
}

// Validates that either Values has more than one element or Value is set, not both
// Allows for Value to be blank and does not check Values[*].Value at all
func IsValueSetInOnePlace(identifier string, rr *schema.ResourceRecord) error {
	if len(rr.Values) > 0 && rr.Value != "" {
		return fmt.Errorf("%s record invalid, can not specify both value and values, identifier: '%s'", rr.Type, identifier)
	}
	return nil
}

// Validates that either Values has more than one element or Comment is set, not both
// Allows for Comment to be blank and does not check Values[*].Comment at all
func IsCommentSetInOnePlace(identifier string, rr *schema.ResourceRecord) error {
	if len(rr.Values) > 0 && rr.Comment != "" {
		return fmt.Errorf("%s record invalid, can not specify both comment and values, identifier: '%s'", rr.Type, identifier)
	}
	return nil
}

// Validates that the name provided matches the RFC1035 regex for valid names according to RFC1035
// and is less then or equal to 255 total characters
func IsValidRFC1035Name(name string) error {
	if len(name) > 255 {
		return fmt.Errorf("Must be less than 255 characters: '%s'", name)
	}

	if !dnsNameRegexRFC1035.MatchString(name) {
		return fmt.Errorf("Does not match regexp '%s': '%s'", dnsNameRegexRFC1035String, name)
	}

	//Split the domain at each part (".") and then run some additional validations
	parts := strings.Split(name, ".")
	for _, part := range parts {
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return fmt.Errorf("Can not start or end with a hyphen (-): '%s'", name)
		}
	}
	return nil
}

// Checks if the name provide is either the wildcard ('@') or is a valid name
func IsValidNameOrWildcard(name string) error {
	// Check if the name matches the regex or is a wildcard
	if name == "@" {
		return nil
	}
	return IsValidRFC1035Name(name)
}

// Formats and email address according to RFC1035
func FormatEmail(email string) (string, error) {
	if strings.Contains(email, "@") {
		// Assume this is a standard email address that will be parseable
		address, err := mail.ParseAddress(email)
		if err != nil {
			return "", fmt.Errorf("Invalid email address: %s", email)
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
func IsFullyQualified(name string) error {
	if err := IsValidRFC1035Name(name); err != nil {
		return err
	}

	if !hasTrailingDot(name) {
		return fmt.Errorf("Must end with a trailing dot: '%s'", name)
	}

	// Count the full stops ('.') in the name, there must be at least two (one for the root and one for the domain)
	if strings.Count(name, ".") < 2 {
		return fmt.Errorf("Must be fully qualified with at least two dots: '%s'", name)
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
