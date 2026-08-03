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

package builtin

import (
	"fmt"
	"strings"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
)

var _ plugins.ZoneMgrPlugin = &BuiltinPluginTXT{}

// RFC1035 3.3: a <character-string> is a single length octet followed by that number of characters,
// limiting each one to 255 characters. The octet length, not the rendered/escaped text length, is what's
// bound by this limit, so all length checks below operate on raw (unescaped) byte counts.
const txtMaxCharacterStringLength = 255

type BuiltinPluginTXT struct {
	plugins.ZoneMgrPlugin
}

func (p *BuiltinPluginTXT) PluginVersion() (string, error) {
	return utils.Version(), nil
}

func (p *BuiltinPluginTXT) PluginTypes() ([]plugins.Type, error) {
	return plugins.PluginTypes(plugins.TXT), nil
}

func (p *BuiltinPluginTXT) Configure(config *models.Config) error {
	// no config
	return nil
}

func (p *BuiltinPluginTXT) Normalize(identifier string, rr *models.ResourceRecord) error {
	if err := validations.CommonValidations(identifier, rr, plugins.TXT); err != nil {
		return err
	}

	if rr.Name == "" {
		rr.Name = identifier
	}

	if err := validations.EnsureValidNameOrWildcard(identifier, rr.Name, rr.Type); err != nil {
		return err
	}

	if _, err := txtValues(identifier, rr); err != nil {
		return err
	}

	return nil
}

func (p *BuiltinPluginTXT) ValidateZone(name string, zone *models.Zone) error {
	// no-op
	return nil
}

func (p *BuiltinPluginTXT) Render(identifier string, rr *models.ResourceRecord) (string, error) {
	if err := validations.EnsureSupportedPluginType(identifier, rr.Type, plugins.TXT); err != nil {
		return "", err
	}

	values, err := txtValues(identifier, rr)
	if err != nil {
		return "", err
	}

	var record strings.Builder
	record.WriteString(rr.RenderResourceWithoutValue())

	quotedValues := make([]string, len(values))
	for i, value := range values {
		quotedValues[i] = quoteTXTValue(value)
	}
	record.WriteString(strings.Join(quotedValues, " "))

	if comment := rr.RetrieveSingleComment(); comment != "" {
		record.WriteString(" ;")
		record.WriteString(comment)
	}

	return record.String(), nil
}

// Returns the RFC1035 <character-string> chunks for the resource record, in order. When Values is populated,
// each entry is treated as an already-split character-string and must not exceed txtMaxCharacterStringLength
// characters; an oversized entry is a validation error rather than being further split. When Values is empty,
// the single Value shortcut is automatically split into txtMaxCharacterStringLength sized chunks.
func txtValues(identifier string, rr *models.ResourceRecord) ([]string, error) {
	if len(rr.Values) == 0 {
		return splitTXTValue(rr.Value), nil
	}

	values := make([]string, len(rr.Values))
	for i, v := range rr.Values {
		if len(v.Value) > txtMaxCharacterStringLength {
			return nil, fmt.Errorf("invalid %s record, character-string exceeds %d characters: '%s', identifier: '%s'", rr.Type, txtMaxCharacterStringLength, v.Value, identifier)
		}
		values[i] = v.Value
	}
	return values, nil
}

// Wraps an already RFC1035 <character-string> sized chunk in double quotes, escaping its content per the
// master file <character-string> syntax defined in RFC1035 5.1.
func quoteTXTValue(value string) string {
	return `"` + escapeTXTValue(value) + `"`
}

// Escapes a bare, unescaped double quote so it can't prematurely end the quoted <character-string> we render
// it inside, per RFC1035 5.1. The value is expected to already carry any escaping the caller wants (\" for a
// literal quote, \\ for a literal backslash, \DDD for an arbitrary byte); a quote is "already escaped" when
// it's immediately preceded by an odd (unpaired) run of backslashes, e.g. \" is left alone but \\" has its
// quote escaped, since the leading \\ there is itself a complete, self-contained escaped backslash.
func escapeTXTValue(value string) string {
	var escaped strings.Builder
	precedingBackslash := false
	for i := 0; i < len(value); i++ {
		b := value[i]
		if b == '"' && !precedingBackslash {
			escaped.WriteByte('\\')
		}
		escaped.WriteByte(b)
		precedingBackslash = b == '\\' && !precedingBackslash
	}
	return escaped.String()
}

// Splits a value into pieces of at most txtMaxCharacterStringLength characters, always returning at least
// one (possibly empty) chunk.
func splitTXTValue(value string) []string {
	if len(value) <= txtMaxCharacterStringLength {
		return []string{value}
	}

	var chunks []string
	for i := 0; i < len(value); i += txtMaxCharacterStringLength {
		chunks = append(chunks, value[i:min(i+txtMaxCharacterStringLength, len(value))])
	}
	return chunks
}

func init() {
	registerBuiltIn(plugins.TXT, &BuiltinPluginTXT{})
}
