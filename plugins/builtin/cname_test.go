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

func TestCNAMEPluginVersion(t *testing.T) {
	testPluginVersion(t, &CNAMEPlugin{})
}

func TestCNAMEPluginTypes(t *testing.T) {
	testPluginTypes(t, &CNAMEPlugin{}, plugins.CNAME)
}

func TestCNAMEConfigure(t *testing.T) {
	testConfigure(t, &CNAMEPlugin{})
}

func TestCNAMENormalize(t *testing.T) {
	setup(t)
	defer teardown(t)
	plugin := &CNAMEPlugin{}
	testNormalizeValidNameAndDefaulting(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.CNAME,
		rrType:     schema.CNAME,
		expects: func(identifier string, rr *schema.ResourceRecord) {
			mockValidations.EXPECT().StandardValidations(identifier, rr, plugins.CNAME)
			if rr.Name == "" {
				mockValidations.EXPECT().IsValidNameOrWildcard(identifier, identifier, rr.Type)
			} else {
				mockValidations.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, rr.Type)
			}

			mockValidations.EXPECT().EnsureNotIP(identifier, rr.RetrieveSingleValue(), rr.Type)
		},
	})
	testNormalizeInvalidName(t, &testNormalize{
		plugin:     plugin,
		pluginType: plugins.CNAME,
		rrType:     schema.CNAME,
		expects: func(identifier string, rr *schema.ResourceRecord) {
			mockValidations.EXPECT().StandardValidations(identifier, rr, plugins.CNAME)
			mockValidations.EXPECT().IsValidNameOrWildcard(identifier, rr.Name, schema.CNAME).Return(fmt.Errorf("not a valid name"))
		},
	})
	testNormalizeValueNotIP(t, plugin, plugins.CNAME, schema.CNAME)
}

func TestCNAMEValidateZone(t *testing.T) {
	testValidateZone(t, &CNAMEPlugin{})
}

func TestCNAMERender(t *testing.T) {
	setup(t)
	defer teardown(t)
	//Render uses the standard method so we're going to cheat
	mockValidations.EXPECT().IsSupportedPluginType("testing", schema.CNAME, plugins.CNAME)
	plugin := &CNAMEPlugin{}
	_, err := plugin.Render("testing", &schema.ResourceRecord{Type: schema.CNAME})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
