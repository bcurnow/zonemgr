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

	"go.uber.org/mock/gomock"

	"github.com/bcurnow/zonemgr/dns/serial"
	"github.com/bcurnow/zonemgr/plugins"
)

var (
	mockController         *gomock.Controller
	mockSerialIndexManager *serial.MockSerialManager
	errTesting             = errors.New("testing error")
)

func setup(t *testing.T) {
	t.Helper()
	mockController = gomock.NewController(t)
	mockSerialIndexManager = serial.NewMockSerialManager(mockController)
	serialIndexManager = mockSerialIndexManager
}

func teardown(_ *testing.T) {
	validations = plugins.V()
	serialIndexManager = nil
	mockController.Finish()
}

func checkErr(t *testing.T, err error, wantErr string) {
	t.Helper()
	if wantErr != "" {
		if err == nil {
			t.Errorf("expected error %q, got nil", wantErr)
		} else if err.Error() != wantErr {
			t.Errorf("got error %q, want %q", err.Error(), wantErr)
		}
	} else if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
