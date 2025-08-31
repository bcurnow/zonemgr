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

	"github.com/bcurnow/zonemgr/utils"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/models/testingutils"
	"github.com/bcurnow/zonemgr/plugins"
)

type testConfig struct {
	plugin     plugins.ZoneMgrPlugin
	pluginType plugins.PluginType
	rrType     models.ResourceRecordType
	expects    func(identifier string, rr *models.ResourceRecord, err bool)
}

type testCase struct {
	identifier string
	err        bool
}

var (
	mockValidator           *plugins.MockValidator
	mockSerialIndexManager  *utils.MockSerialIndexManager
	mockSoaValuesNormalizer *plugins.MockSOAValuesNormalizer
	testingError            error
)

func setup(t *testing.T) {
	testingutils.Setup(t)
	mockValidator = plugins.NewMockValidator(testingutils.MockController)
	validations = mockValidator
	mockSerialIndexManager = utils.NewMockSerialIndexManager(testingutils.MockController)
	serialIndexManager = mockSerialIndexManager
	mockSoaValuesNormalizer = plugins.NewMockSOAValuesNormalizer(testingutils.MockController)
	soaValuesNormalizer = mockSoaValuesNormalizer
	testingError = errors.New("testing error")
}

func teardown(t *testing.T) {
	testingutils.Teardown(t)
	validations = plugins.V()
	serialIndexManager = nil
}

// Runs the testing necessary for any resource type where the name needs to be a valid name or wildcard
// NOTE: setup needs to be called before calling this method!
func testCommonValidations(t *testing.T, testConf *testConfig, rr *models.ResourceRecord) {
	testCases := []testCase{
		{
			identifier: "CommonValidations-Valid",
		},
		{
			identifier: "CommonValidations-Error",
			err:        true,
		},
	}

	for _, tc := range testCases {
		testConf.expects(tc.identifier, rr, tc.err)
		err := testConf.plugin.Normalize(tc.identifier, rr)
		handleError(t, err, tc.err)
	}
}

func testIsValidNameOrWildcard(t *testing.T, testConf *testConfig, rr *models.ResourceRecord) {
	testCases := []testCase{
		{
			identifier: "IsValidNameOrWildcard-Valid",
		},
		{
			identifier: "IsValidNameOrWildcard-Error",
			err:        true,
		},
	}

	for _, tc := range testCases {
		testConf.expects(tc.identifier, rr, tc.err)
		err := testConf.plugin.Normalize(tc.identifier, rr)
		handleError(t, err, tc.err)
	}
}

func testEnsureIP(t *testing.T, testConf *testConfig, rr *models.ResourceRecord) {
	testCases := []testCase{
		{
			identifier: "EnsureIP-Valid",
		},
		{
			identifier: "EnsureIP-Error",
			err:        true,
		},
	}

	for _, tc := range testCases {
		testConf.expects(tc.identifier, rr, tc.err)
		err := testConf.plugin.Normalize(tc.identifier, rr)
		handleError(t, err, tc.err)
	}
}

func testEnsureNotIP(t *testing.T, testConf *testConfig, rr *models.ResourceRecord) {
	testCases := []testCase{
		{
			identifier: "EnsureNotIP-Valid",
		},
		{
			identifier: "EnsureNotIP-Error",
			err:        true,
		},
	}

	for _, tc := range testCases {
		testConf.expects(tc.identifier, rr, tc.err)
		err := testConf.plugin.Normalize(tc.identifier, rr)
		handleError(t, err, tc.err)
	}
}

func testRender(t *testing.T, testConf testConfig, rr *models.ResourceRecord) {
	testCases := []testCase{
		{
			identifier: "Render-Valid",
		},
		{
			identifier: "Render-WrongPluginType",
			err:        true,
		},
	}
	for _, tc := range testCases {
		testConf.expects(tc.identifier, rr, tc.err)
		_, err := testConf.plugin.Render(tc.identifier, rr)
		handleError(t, err, tc.err)
	}
}

func testIsFullyQualified(t *testing.T, testConf *testConfig, rr *models.ResourceRecord) {
	testCases := []testCase{
		{
			identifier: "IsFullyQualified-Valid",
		},
		{
			identifier: "IsFullyQualified-Error",
			err:        true,
		},
	}
	for _, tc := range testCases {
		testConf.expects(tc.identifier, rr, tc.err)
		err := testConf.plugin.Normalize(tc.identifier, rr)
		handleError(t, err, tc.err)
	}
}

func handleError(t *testing.T, err error, wantError bool) {
	if wantError {
		handleCustomError(t, err, testingError)
	} else {
		handleCustomError(t, err, nil)
	}
}

func handleCustomError(t *testing.T, err error, wanted error) {
	if wanted != nil && err == nil {
		t.Errorf("expected error")
	} else if wanted != nil && err != nil {
		if err.Error() != wanted.Error() {
			t.Errorf("Unexpected error:\n'%s'\nwant\n'%s'", err, wanted)
		}
	} else if wanted == nil && err != nil {
		t.Errorf("unexpected error:\n'%s'", err)
	}
}
