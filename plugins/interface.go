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
)

type ZoneMgrPlugin interface {
	// Returns the version of plugin
	PluginVersion() (string, error)
	// Returns the set of plugin types that this plugin supports
	PluginTypes() ([]PluginType, error)
	// Allows for configuration of the plugin, this will be called once for each zone in the file
	Configure(config *schema.Config) error
	// Allows for validation and normalization/defaulting for the resource record
	Normalize(identifier string, rr *schema.ResourceRecord) error
	// Allows for validation of the entire normalized zone
	// This enables checks such as all CNAME records properly referencing a defined A record
	// This allows validation only, no defaulting
	ValidateZone(name string, zone *schema.Zone) error
	// Converts the resource record into a string to be writting out to a file
	Render(identifier string, rr *schema.ResourceRecord) (string, error)
}

type Validator interface {
	// Performs the standard validations for resource records
	// This includes:
	//   - Validation that the resource record is of the specified type - This is not case insensitive but the the type will be normalized to uppercase
	//   - Validation of the class - An empty class will be considered valid, any defaulting or enforcement beyond that is the responsiblity of the individual plugins
	//   - Validation that only Value or Values is populated
	//   - Validation that only Comment or Values is populated
	StandardValidations(identifier string, rr *schema.ResourceRecord, supportedTypes ...PluginType) error
	// Checks if the supplied resource record matches one of the support plugin types
	IsSupportedPluginType(identifier string, rrType schema.ResourceRecordType, supportedTypes ...PluginType) error
	// Validates that the name provided matches the RFC1035 regex for valid names according to RFC1035
	// and is less then or equal to 255 total characters
	IsValidRFC1035Name(identifier string, name string, rrType schema.ResourceRecordType) error
	// Checks if the name provide is either the wildcard ('@') or is a valid name
	IsValidNameOrWildcard(identifier string, name string, rrType schema.ResourceRecordType) error
	// Formats and email address according to RFC1035
	FormatEmail(identifier string, email string, rrType schema.ResourceRecordType) (string, error)
	// Most DNS names in a zone file need to be fully qualified domain names, while we can't validate if the entire name itself is valid,
	// we can ensure that it is a valid name and ends with a trailing dot
	IsFullyQualified(identifier string, name string, rrType schema.ResourceRecordType) error
	// Ensure that the string passed in ends with a trailing dot
	EnsureTrailingDot(name string) string
	// Ensure that the string is an IP address
	EnsureIP(identifier string, s string, rrType schema.ResourceRecordType) error
	// Ensure that the string is NOT an IP address
	EnsureNotIP(identifier string, s string, rrType schema.ResourceRecordType) error
	//  Ensurre that the string is a 32-bit integer which is positive and greater than zero
	IsPositive(identifier string, s string, fieldName string, rrType schema.ResourceRecordType) error
}
