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

func TestAPluginVersion(t *testing.T) {
	// Make sure test helpers (like testPluginVersion) are not in plugins or schema packages.
	testPluginVersion(t, &APlugin{})
}

func TestAPluginTypes(t *testing.T) {
	testPluginTypes(t, &APlugin{}, plugins.A)
}

func TestAConfigure(t *testing.T) {
	testConfigure(t, &APlugin{})
}

func TestANormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &APlugin{}
	testNormalizeValidNameAndDefaulting(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.A,
		rrType:     models.A,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, plugins.A)
			if rr.Name == "" {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, rr.Type)
			} else {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
			}

			mockValidator.EXPECT().EnsureIP(identifier, rr.RetrieveSingleValue(), rr.Type)
		},
	})
	testNormalizeInvalidName(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.A,
		rrType:     models.A,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, plugins.A)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, models.A).Return(fmt.Errorf("not a valid name"))
		},
	})
	testNormalizeValueIsIP(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.A,
		rrType:     models.A,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, plugins.A)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, models.A)
			mockValidator.EXPECT().EnsureIP(identifier, rr.Value, models.A).Return(fmt.Errorf("is not IP"))
		},
	})
}

func TestAValidateZone(t *testing.T) {
	testValidateZone(t, &APlugin{})
}

func TestARender(t *testing.T) {
	setup(t)
	defer teardown(t)
	//Render uses the standard method so we're going to cheat
	mockValidator.EXPECT().IsSupportedPluginType("testing", models.A, plugins.A)
	plugin := &APlugin{}
	_, err := plugin.Render("testing", &models.ResourceRecord{Type: models.A})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
