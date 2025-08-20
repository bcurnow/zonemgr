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
	"fmt"
	"testing"

	"github.com/bcurnow/zonemgr/models"
)

func TestPTRPluginVersion(t *testing.T) {
	testPluginVersion(t, &BuiltinPluginPTR{})
}

func TestPTRPluginTypes(t *testing.T) {
	testPluginTypes(t, &BuiltinPluginPTR{}, PTR)
}

func TestPTRConfigure(t *testing.T) {
	testConfigure(t, &BuiltinPluginPTR{})
}

func TestPTRNormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &BuiltinPluginPTR{}
	testNormalizeValidNameAndDefaulting(t, &testNormalize{
		plugin:     plugin,
		pluginType: PTR,
		rrType:     models.PTR,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, PTR)
			if rr.Name == "" {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, rr.Type)
			} else {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
			}

			mockValidator.EXPECT().IsFullyQualified(identifier, rr.RetrieveSingleValue(), rr.Type)
		},
	})
	testNormalizeInvalidName(t, &testNormalize{
		plugin:     plugin,
		pluginType: PTR,
		rrType:     models.PTR,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, PTR)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, models.PTR).Return(fmt.Errorf("not a valid name"))

		},
	})
	testNormalizeValueNotFullyQualified(t, &testNormalize{
		plugin:     plugin,
		pluginType: PTR,
		rrType:     models.PTR,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, PTR)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, models.PTR)
			mockValidator.EXPECT().IsFullyQualified(identifier, rr.RetrieveSingleValue(), models.PTR).Return(fmt.Errorf("not fully qualified"))
		},
	})
}

func TestPTRValidateZone(t *testing.T) {
	testValidateZone(t, &BuiltinPluginPTR{})
}

func TestPTRRender(t *testing.T) {
	setup(t)
	defer teardown(t)
	//Render uses the standard method so we're going to cheat
	mockValidator.EXPECT().IsSupportedPluginType("testing", models.PTR, PTR)
	plugin := &BuiltinPluginPTR{}
	_, err := plugin.Render("testing", &models.ResourceRecord{Type: models.PTR})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
