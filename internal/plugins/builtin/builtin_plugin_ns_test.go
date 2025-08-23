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

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
)

func TestNSPluginVersion(t *testing.T) {
	testPluginVersion(t, &BuiltinPluginNS{})
}

func TestNSPluginTypes(t *testing.T) {
	testPluginTypes(t, &BuiltinPluginNS{}, plugins.NS)
}

func TestNSConfigure(t *testing.T) {
	testConfigure(t, &BuiltinPluginNS{})
}

func TestNSNormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &BuiltinPluginNS{}
	testNormalizeValidNameAndDefaulting(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     models.NS,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.NS)

			if rr.Name == "" {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, "@", rr.Type)
			} else {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
			}

			if rr.RetrieveSingleValue() == "" {
				mockValidator.EXPECT().EnsureIP(identifier, identifier, rr.Type)
				mockValidator.EXPECT().IsFullyQualified(identifier, identifier, rr.Type)
			} else {
				mockValidator.EXPECT().EnsureIP(identifier, rr.RetrieveSingleValue(), rr.Type)
				mockValidator.EXPECT().IsFullyQualified(identifier, rr.RetrieveSingleValue(), rr.Type)
			}

		},
	})
	testNormalizeInvalidName(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     models.NS,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.NS)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, models.NS).Return(fmt.Errorf("not a valid name"))
		},
	})
	testNormalizeValueIsIP(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     models.NS,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.NS)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, models.NS)
			mockValidator.EXPECT().EnsureIP(identifier, rr.Value, models.NS).Return(fmt.Errorf("is not IP"))

		},
	})
	testNormalizeValueNotFullyQualified(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.NS,
		rrType:     models.NS,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.NS)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, models.NS)
			mockValidator.EXPECT().EnsureIP(identifier, rr.RetrieveSingleValue(), models.NS)
			mockValidator.EXPECT().IsFullyQualified(identifier, rr.RetrieveSingleValue(), models.NS).Return(fmt.Errorf("not fully qualified"))
		},
	})
}

func TestNSValidateZone(t *testing.T) {
	testValidateZone(t, &BuiltinPluginNS{})
}

func TestNSRender(t *testing.T) {
	setup(t)
	defer teardown(t)
	//Render uses the standard method so we're going to cheat
	mockValidator.EXPECT().IsSupportedPluginType("testing", models.NS, plugins.NS)
	plugin := &BuiltinPluginNS{}
	_, err := plugin.Render("testing", &models.ResourceRecord{Type: models.NS})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
