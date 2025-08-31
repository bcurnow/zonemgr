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

func TestNSNormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	tc := &testConfig{
		plugin:     &BuiltinPluginNS{},
		pluginType: plugins.NS,
		rrType:     models.NS,
	}
	rr := &models.ResourceRecord{
		Type:  tc.rrType,
		Name:  "nsrecord",
		Value: "ns1.example.com.",
	}
	tc.expects = normalizeExpects_NSPlugin(true, false, false, false)
	testCommonValidations(t, tc, rr)
	tc.expects = normalizeExpects_NSPlugin(false, true, false, false)
	testIsValidNameOrWildcard(t, tc, rr)
	tc.expects = normalizeExpects_NSPlugin(false, false, true, false)
	testEnsureNotIP(t, tc, rr)
	tc.expects = normalizeExpects_NSPlugin(false, false, false, true)
	testIsFullyQualified(t, tc, rr)

	// Test value defaulting
	identifier := "testing-value-defaulting"
	noValue := &models.ResourceRecord{
		Name: "valuedefaulting",
		Type: tc.rrType,
	}
	mockValidator.EXPECT().CommonValidations(identifier, noValue, plugins.NS)
	mockValidator.EXPECT().IsValidNameOrWildcard(identifier, noValue.Name, noValue.Type)
	// Make sure the name defaulted
	mockValidator.EXPECT().EnsureNotIP(identifier, identifier, noValue.Type)
	mockValidator.EXPECT().IsFullyQualified(identifier, identifier, noValue.Type)

	if err := tc.plugin.Normalize(identifier, noValue); err != nil {
		t.Errorf("unexpected error:\n'%s'", err)
	}

	// Test name defaulting
	identifier = "testing-name-defaulting"
	noName := &models.ResourceRecord{
		Type:  tc.rrType,
		Value: "ns1.example.com.",
	}
	mockValidator.EXPECT().CommonValidations(identifier, noName, plugins.NS)
	// Make sure the name defaulted
	mockValidator.EXPECT().IsValidNameOrWildcard(identifier, "@", noName.Type)
	mockValidator.EXPECT().EnsureNotIP(identifier, noName.RetrieveSingleValue(), noName.Type)
	mockValidator.EXPECT().IsFullyQualified(identifier, noName.RetrieveSingleValue(), noName.Type)

	if err := tc.plugin.Normalize(identifier, noName); err != nil {
		t.Errorf("unexpected error:\n'%s'", err)
	}
}

func TestNSRender(t *testing.T) {
	setup(t)
	defer teardown(t)
	rr := &models.ResourceRecord{
		Type: models.NS,
		Name: "render",
	}
	plugin := &BuiltinPluginNS{}
	pluginType := plugins.NS
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

func normalizeExpects_NSPlugin(commonValidationsErr bool, isValidNameOrWildcardErr bool, ensureIPErr bool, isFullyQualifiedErr bool) func(identifier string, rr *models.ResourceRecord, err bool) {
	return func(identifier string, rr *models.ResourceRecord, err bool) {
		call := mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.NS)
		if commonValidationsErr && err {
			call.Return(testingError)
			return
		}
		call = mockValidator.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
		if isValidNameOrWildcardErr && err {
			call.Return(testingError)
			return
		}
		call = mockValidator.EXPECT().EnsureNotIP(identifier, rr.RetrieveSingleValue(), rr.Type)
		if ensureIPErr && err {
			call.Return(testingError)
			return
		}
		call = mockValidator.EXPECT().IsFullyQualified(identifier, rr.RetrieveSingleValue(), rr.Type)
		if isFullyQualifiedErr && err {
			call.Return(testingError)
			return
		}
	}
}
