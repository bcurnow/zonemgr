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
package normalize

import (
	"fmt"
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

func formatEmail(email string) (string, error) {
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
