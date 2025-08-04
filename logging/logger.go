package logging

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

var logger hclog.Logger

// Gets the main logger for the application
func Logger() hclog.Logger {
	return logger
}

func init() {
	// Create a logger, this wil print to stderr and use standard logging formats (e.g. timestamp)
	logger = hclog.New(&hclog.LoggerOptions{
		Name:   "zonemgr",
		Output: os.Stderr,
		Level:  hclog.Trace,
	})
}
