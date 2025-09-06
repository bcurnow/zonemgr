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

	"github.com/spf13/viper"
)

func TestRun_Env(t *testing.T) {
	v = viper.New()
	v.Set("test2", "testing")
	v.Set("test1", "testing")

	// The entire function is essentially logging so capture stdout (this doesn't log to stderr)
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("could not create pipe to capture stdout")
	}
	os.Stdout = w

	envCmd.Run(pluginsCmd, []string{})
	w.Close()

	wanted := "test1=\"testing\"\ntest2=\"testing\"\n"

	var buf bytes.Buffer
	buf.ReadFrom(r)
	r.Close()

	if buf.String() != wanted {
		t.Errorf("Invalid plugin output:\n%s\nwanted:\n%s", buf.String(), wanted)
	}
}
