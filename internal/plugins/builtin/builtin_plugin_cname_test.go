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
	"errors"
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
)

func TestNormalize_CNAMEPlugin(t *testing.T) {
	setup(t)
	defer teardown(t)
	tc := &testConfig{
		plugin:     &BuiltinPluginCNAME{},
		pluginType: plugins.CNAME,
		rrType:     models.CNAME,
	}
	rr := &models.ResourceRecord{
		Type:  tc.rrType,
		Name:  "cnameplugin",
		Value: "notanip",
	}
	tc.expects = normalizeExpects_CNAMEPlugin(true, false, false)
	testCommonValidations(t, tc, rr)
	tc.expects = normalizeExpects_CNAMEPlugin(false, true, false)
	testIsValidNameOrWildcard(t, tc, rr)
	tc.expects = normalizeExpects_CNAMEPlugin(false, false, true)
	testEnsureNotIP(t, tc, rr)

	// Now test the defaulting logic
	identifier := "testing-name-defaulting"
	noName := &models.ResourceRecord{
		Type:  tc.rrType,
		Value: "notanip",
	}

	mockValidator.EXPECT().CommonValidations(identifier, noName, tc.pluginType)
	// make sure the name defaulted
	mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, rr.Type)
	mockValidator.EXPECT().EnsureNotIP(identifier, noName.RetrieveSingleValue(), rr.Type)
	if err := tc.plugin.Normalize(identifier, noName); err != nil {
		t.Errorf("unexpected error:\n'%s'", err)
	}
}

func TestValidateZone_CNAMEPlugin(t *testing.T) {
	plugin := &BuiltinPluginCNAME{}
	testCases := []struct {
		zone *models.Zone
		err  error
	}{
		{
			zone: &models.Zone{ResourceRecords: map[string]*models.ResourceRecord{}},
		},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"lonecname": {
						Name: "lonecname",
						Type: models.CNAME,
					},
				},
			},
			err: errors.New("found CNAME records but there are no A records present, all CNAMES must reference an A record name, zone: 'testing'"),
		},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"cname": {
						Name:  "badcname",
						Type:  models.CNAME,
						Value: "notgoingtofindme",
					},
					"arecord": {
						Name:  "nottherightone",
						Type:  models.A,
						Value: "1.2.3.4",
					},
				},
			},
			err: errors.New("invalid CNAME record, 'cname' has a value of 'notgoingtofindme' which does not match any defined A record name, zone: 'testing'"),
		},
		{
			zone: &models.Zone{
				ResourceRecords: map[string]*models.ResourceRecord{
					"arecord": {
						Name:  "arecord",
						Type:  models.A,
						Value: "1.2.3.4",
					},
					"cname": {
						Name:  "cname",
						Type:  models.CNAME,
						Value: "arecord",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		if err := plugin.ValidateZone("testing", tc.zone); err != nil {
			handleCustomError(t, err, tc.err)
		}
	}
}

func TestRender_CNAMEPlugin(t *testing.T) {
	setup(t)
	defer teardown(t)
	rr := &models.ResourceRecord{
		Type: models.CNAME,
		Name: "render",
	}
	plugin := &BuiltinPluginCNAME{}
	pluginType := plugins.CNAME
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

func normalizeExpects_CNAMEPlugin(commonValidationsErr bool, isValidNameOrWildcardErr bool, ensureNotIPErr bool) func(identifier string, rr *models.ResourceRecord, err bool) {
	return func(identifier string, rr *models.ResourceRecord, err bool) {
		call := mockValidator.EXPECT().CommonValidations(identifier, rr, plugins.CNAME)
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
		if ensureNotIPErr && err {
			call.Return(testingError)
		}
	}
}
