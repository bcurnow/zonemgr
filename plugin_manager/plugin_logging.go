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

package plugin_manager

import (
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
)

func PluginLogger() hclog.Logger {
	if pluginDebug {
		return hclog.L().Named("plugin")
	}
	return hclog.New(&hclog.LoggerOptions{
		Name:  "plugin",
		Level: hclog.Off,
	})
}

func PluginStdout() io.Writer {
	if pluginDebug {
		return os.Stdout
	}
	return io.Discard
}

func PluginStderr() io.Writer {
	if pluginDebug {
		return os.Stderr
	}
	return io.Discard
}
