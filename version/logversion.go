package version

import (
	"path/filepath"
	"runtime"

	"github.com/cybriq/glint/pkg/proc"
)

var F, E, W, I, D, T proc.LevelPrinter

func init() {
	_, file, _, _ := runtime.Caller(0)
	verPath := filepath.Dir(file) + "/"
	F, E, W, I, D, T = proc.GetLogPrinterSet(proc.AddLoggerSubsystem(verPath))
}
