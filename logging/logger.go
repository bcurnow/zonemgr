package logging

import (
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
)

var pluginDebug = false

func EnablePluginDebug() {
	pluginDebug = true
}

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

// This function updates the hclog.DefaultOptions. It should be called as early as possible because
// and hclog.Default() or hclog.L() calls made before this is called will get the defaults
func ConfigureLogging(level hclog.Level, jsonFormat bool, disableTime bool, logColor bool) {
	if logColor {
		hclog.DefaultOptions.Color = hclog.AutoColor
	} else {
		hclog.DefaultOptions.Color = hclog.ColorOff
	}

	hclog.DefaultOptions.Level = level
	hclog.DefaultOptions.JSONFormat = jsonFormat
	hclog.DefaultOptions.DisableTime = disableTime
}

func defaultLogging() {
	hclog.DefaultOptions.Name = "zonemgr"
	hclog.DefaultOptions.Output = os.Stderr
	hclog.DefaultOptions.Level = hclog.Info
	hclog.DefaultOptions.Color = hclog.AutoColor
	hclog.DefaultOptions.JSONFormat = false
	hclog.DefaultOptions.DisableTime = true
}

func init() {
	defaultLogging()
}
