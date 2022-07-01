package main

import (
	"github.com/cybriq/glint/pkg/proc"

	"github.com/cybriq/glint/version"
)

var F, E, W, I, D, T = proc.GetLogPrinterSet(proc.AddLoggerSubsystem(version.PathBase))
