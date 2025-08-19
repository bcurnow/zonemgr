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
	"testing"

	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/schema"
)

func TestNSPluginVersion(t *testing.T) {
	testPluginVersion(t, &NSPlugin{})
}

func TestNSPluginTypes(t *testing.T) {
	testPluginTypes(t, &NSPlugin{}, plugins.NS)
}

func TestNSConfigure(t *testing.T) {
	testConfigure(t, &NSPlugin{})
}

func TestNSNormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &NSPlugin{}
	testNormalizeValidNameAndDefaulting(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     schema.NS,
		expects: func(identifier string, rr *schema.ResourceRecord) {
			mockValidations.EXPECT().StandardValidations(identifier, rr, plugins.NS)

			if rr.Name == "" {
				mockValidations.EXPECT().IsValidNameOrWildcard(identifier, "@", rr.Type)
			} else {
				mockValidations.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
			}

			if rr.RetrieveSingleValue() == "" {
				mockValidations.EXPECT().EnsureIP(identifier, identifier, rr.Type)
				mockValidations.EXPECT().IsFullyQualified(identifier, identifier, rr.Type)
			} else {
				mockValidations.EXPECT().EnsureIP(identifier, rr.RetrieveSingleValue(), rr.Type)
				mockValidations.EXPECT().IsFullyQualified(identifier, rr.RetrieveSingleValue(), rr.Type)
			}

		},
	})
	testNormalizeInvalidName(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     schema.NS,
		expects: func(identifier string, rr *schema.ResourceRecord) {
			mockValidations.EXPECT().StandardValidations(identifier, rr, plugins.NS)
			mockValidations.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, schema.NS).Return(fmt.Errorf("not a valid name"))
		},
	})
	testNormalizeValueIsIP(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     schema.NS,
		expects: func(identifier string, rr *schema.ResourceRecord) {
			mockValidations.EXPECT().StandardValidations(identifier, rr, plugins.NS)
			mockValidations.EXPECT().IsValidNameOrWildcard(identifier, identifier, schema.NS)
			mockValidations.EXPECT().EnsureIP(identifier, rr.Value, schema.NS).Return(fmt.Errorf("is not IP"))

		},
	})
	testNormalizeValueNotFullyQualified(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     schema.NS,
		expects: func(identifier string, rr *schema.ResourceRecord) {
			mockValidations.EXPECT().StandardValidations(identifier, rr, plugins.NS)
			mockValidations.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, schema.NS)
			mockValidations.EXPECT().EnsureIP(identifier, rr.RetrieveSingleValue(), schema.NS)
			mockValidations.EXPECT().IsFullyQualified(identifier, rr.RetrieveSingleValue(), schema.NS).Return(fmt.Errorf("not fully qualified"))
		},
	})
}

func TestNSValidateZone(t *testing.T) {
	testValidateZone(t, &NSPlugin{})
}

func TestNSRender(t *testing.T) {
	setup(t)
	defer teardown(t)
	//Render uses the standard method so we're going to cheat
	mockValidations.EXPECT().IsSupportedPluginType("testing", schema.NS, plugins.NS)
	plugin := &NSPlugin{}
	_, err := plugin.Render("testing", &schema.ResourceRecord{Type: schema.NS})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
