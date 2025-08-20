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

func TestCNAMEPluginVersion(t *testing.T) {
	testPluginVersion(t, &BuiltinPluginCNAME{})
}

func TestCNAMEPluginTypes(t *testing.T) {
	testPluginTypes(t, &BuiltinPluginCNAME{}, CNAME)
}

func TestCNAMEConfigure(t *testing.T) {
	testConfigure(t, &BuiltinPluginCNAME{})
}

func TestCNAMENormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &BuiltinPluginCNAME{}
	testNormalizeValidNameAndDefaulting(t, &testNormalize{
		plugin:     plugin,
		pluginType: CNAME,
		rrType:     models.CNAME,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, CNAME)
			if rr.Name == "" {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, rr.Type)
			} else {
				mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
			}

			mockValidator.EXPECT().EnsureNotIP(identifier, rr.RetrieveSingleValue(), rr.Type)
		},
	})
	testNormalizeInvalidName(t, &testNormalize{
		plugin:     plugin,
		pluginType: CNAME,
		rrType:     models.CNAME,
		expects: func(identifier string, rr *models.ResourceRecord) {
			mockValidator.EXPECT().StandardValidations(identifier, rr, CNAME)
			mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, models.CNAME).Return(fmt.Errorf("not a valid name"))
		},
	})
	testNormalizeValueNotIP(t, plugin, CNAME, models.CNAME)
}

func TestCNAMEValidateZone(t *testing.T) {
	testValidateZone(t, &BuiltinPluginCNAME{})
}

func TestCNAMERender(t *testing.T) {
	setup(t)
	defer teardown(t)
	//Render uses the standard method so we're going to cheat
	mockValidator.EXPECT().IsSupportedPluginType("testing", models.CNAME, CNAME)
	plugin := &BuiltinPluginCNAME{}
	_, err := plugin.Render("testing", &models.ResourceRecord{Type: models.CNAME})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
