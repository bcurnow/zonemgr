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
	"testing"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/plugins"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

type testNormalize struct {
	plugin     plugins.ZoneMgrPlugin
	pluginType plugins.PluginType
	rrType     models.ResourceRecordType
	expects    func(identifier string, rr *models.ResourceRecord)
}

var (
	mockController *gomock.Controller
	mockValidator  *plugins.MockValidator
)

func setup(t *testing.T) {
	mockController = gomock.NewController(t)
	mockValidator = plugins.NewMockValidator(mockController)
	validations = mockValidator
}

func teardown(_ *testing.T) {
	mockController.Finish()
	validations = plugins.V()
}

func testPluginVersion(t *testing.T, p plugins.ZoneMgrPlugin) {
	ver, err := p.PluginVersion()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if ver != utils.Version() {
		t.Errorf("incorrect version %s, want %s", ver, utils.Version())
	}
}

func testPluginTypes(t *testing.T, p plugins.ZoneMgrPlugin, want ...plugins.PluginType) {
	pluginTypes, err := p.PluginTypes()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !cmp.Equal(pluginTypes, want) {
		t.Errorf("unexpected plugin types %s, want %s", pluginTypes, want)
	}
}

func testConfigure(t *testing.T, p plugins.ZoneMgrPlugin) {
	if err := p.Configure(nil); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func testValidateZone(t *testing.T, p plugins.ZoneMgrPlugin) {
	if err := p.ValidateZone("noop", &models.Zone{}); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

// Runs the testinng necessary for any resouurce type where the name needs to be a valid name or wildcard
// NOTE: setup needs to be called before calling this method!
func testNormalizeValidNameAndDefaulting(t *testing.T, tn *testNormalize) {
	testCases := []struct {
		rr         *models.ResourceRecord
		identifier string
	}{
		{
			identifier: "ValidRecordWithOutName",
			rr: &models.ResourceRecord{
				Type:  tn.rrType,
				Value: "value.example.com.",
			},
		},
		{
			identifier: "Valid record with a name",
			rr: &models.ResourceRecord{
				Type:  tn.rrType,
				Value: "value.example.com.",
				Name:  "name",
			},
		},
	}

	for _, tc := range testCases {
		tn.expects(tc.identifier, tc.rr)
		err := tn.plugin.Normalize(tc.identifier, tc.rr)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}
}

// Runs the testinng necessary for any resouurce type where the name needs to be a valid name or wildcard
// NOTE: setup needs to be called before calling this method!
func testNormalizeInvalidName(t *testing.T, tn *testNormalize) {
	identifier := "Invalid name"
	// Make a name that has more than 63 characters to voilate the spec
	rr := &models.ResourceRecord{Type: tn.rrType, Name: "invalidname" + strings.Repeat("e", 63)}
	tn.expects(identifier, rr)
	err := tn.plugin.Normalize(identifier, rr)
	if err == nil {
		t.Errorf("expected error")
	} else {
		if err.Error() != "not a valid name" {
			t.Errorf("unexpected error: %s, want not a valid name", err)
		}
	}
}

// Runs the testinng necessary for any resouurce type where the Value needs to be an IP address
// NOTE: setup needs to be called before calling this method!
func testNormalizeValueIsIP(t *testing.T, tn *testNormalize) {
	identifier := "name"
	rr := &models.ResourceRecord{Type: tn.rrType, Value: "not an IP", Name: identifier}

	tn.expects(identifier, rr)
	err := tn.plugin.Normalize(identifier, rr)
	if err == nil {
		t.Errorf("expected error")
	} else {
		if err.Error() != "is not IP" {
			t.Errorf("unexpected error: '%s', want is not IP", err.Error())
		}
	}
}

// Runs the testinng necessary for any resouurce type where the Value needs to NOT be an IP address
// NOTE: setup needs to be called before calling this method!
func testNormalizeValueNotIP(t *testing.T, plugin plugins.ZoneMgrPlugin, pluginType plugins.PluginType, rrType models.ResourceRecordType) {
	identifier := "name"
	rr := &models.ResourceRecord{Type: rrType, Value: "not an IP", Name: identifier}

	mockValidator.EXPECT().StandardValidations(identifier, rr, pluginType)
	mockValidator.EXPECT().IsValidNameOrWildcard(identifier, identifier, rrType)
	mockValidator.EXPECT().EnsureNotIP(identifier, rr.Value, rrType).Return(fmt.Errorf("is an IP"))

	err := plugin.Normalize(identifier, rr)
	if err == nil {
		t.Errorf("expected error")
	} else {
		if err.Error() != "is an IP" {
			t.Errorf("unexpected error: %s, want is an IP", err)
		}
	}
}

func testNormalizeValueNotFullyQualified(t *testing.T, tn *testNormalize) {
	identifier := "name"
	rr := &models.ResourceRecord{Type: tn.rrType, Name: "name", Value: "value.example.com."}
	tn.expects(identifier, rr)
	err := tn.plugin.Normalize(identifier, rr)
	if err == nil {
		t.Errorf("expected error")
	} else {
		if err.Error() != "not fully qualified" {
			t.Errorf("unexpected error: %s, want not fully qualified", err)
		}
	}
}
