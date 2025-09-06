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

package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/bcurnow/zonemgr/utils"
)

func TestRun_Version(t *testing.T) {
	// The entire function is essentially logging so capture stdout (this doesn't log to stderr)
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("could not create pipe to capture stdout")
	}
	os.Stdout = w

	versionCmd.Run(pluginsCmd, []string{})
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	r.Close()

	if buf.String() != utils.Version()+"\n" {
		t.Errorf("Invalid plugin output:\n%s\nwanted:\n%s", buf.String(), utils.Version()+"\n")
	}
}
