//go:build !windows
// +build !windows

package proc

import (
	"syscall"
)

// Pause sends a signal to the worker process to stop
func (w *Worker) Pause() (e error) {
	if e = w.Cmd.Process.Signal(syscall.SIGSTOP); !log.E.Chk(e) {
		log.D.Ln("paused")
	}
	return
}

// Continue sends a signal to a worker process to resume work
func (w *Worker) Continue() (e error) {
	if e = w.Cmd.Process.Signal(syscall.SIGCONT); !log.E.Chk(e) {
		log.D.Ln("resumed")
	}
	return
}
