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
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
)

func TestNormalize_APlugin(t *testing.T) {
	setup(t)
	defer teardown(t)
	tc := &testConfig{
		plugin:     &BuiltinPluginA{},
		pluginType: plugins.A,
		rrType:     models.A,
	}
	rr := &models.ResourceRecord{
		Type:  tc.rrType,
		Name:  "aplugin",
		Value: "1.2.3.4",
	}
	tc.expects = normalizeExpects_APlugin(true, false, false)
	testCommonValidations(t, tc, rr)
	tc.expects = normalizeExpects_APlugin(false, true, false)
	testIsValidNameOrWildcard(t, tc, rr)
	tc.expects = normalizeExpects_APlugin(false, false, true)
	testEnsureIP(t, tc, rr)

	// Test name defaulting
	identifier := "testing-name-defaulting"
	noName := &models.ResourceRecord{
		Type:  tc.rrType,
		Value: "1.2.3.4",
	}
	mockValidator.EXPECT().CommonValidations(identifier, noName, tc.pluginType)
	// Make sure the name defaulted
	mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, rr.Type)
	mockValidator.EXPECT().EnsureIP(identifier, noName.RetrieveSingleValue(), rr.Type)

	if err := tc.plugin.Normalize(identifier, noName); err != nil {
		t.Errorf("unexpected error:\n'%s'", err)
	}
}

func TestRender_APlugin(t *testing.T) {
	setup(t)
	defer teardown(t)
	rr := &models.ResourceRecord{
		Type: models.A,
		Name: "render",
	}
	plugin := &BuiltinPluginA{}
	pluginType := plugins.A
	testRender(t, testConfig{
		plugin:     plugin,
		pluginType: pluginType,
		rrType:     rr.Type,
		expects: func(identifier string, rr *models.ResourceRecord, err bool) {
			call := mockValidator.EXPECT().IsSupportedPluginType(identifier, rr.Type, pluginType)
			if err {
				call.Return(testingError)
			}
		},
	}, rr)
	//Render uses the standard method so we're going to cheat
	mockValidator.EXPECT().IsSupportedPluginType("testing", rr.Type, pluginType)
	_, err := plugin.Render("testing", rr)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func normalizeExpects_APlugin(commonValidationsErr bool, isValidNameOrWildcardErr bool, ensureIPErr bool) func(identifier string, rr *models.ResourceRecord, err bool) {
	return func(identifier string, rr *models.ResourceRecord, err bool) {
		call := mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.A)
		if commonValidationsErr && err {
			call.Return(testingError)
			return
		}
		call = mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
		if isValidNameOrWildcardErr && err {
			call.Return(testingError)
			return
		}
		call = mockValidator.EXPECT().EnsureIP(identifier, rr.RetrieveSingleValue(), rr.Type)
		if ensureIPErr && err {
			call.Return(testingError)
		}
	}
}
