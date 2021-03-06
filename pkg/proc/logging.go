package proc

import (
	"github.com/cybriq/gotiny"
	"github.com/cybriq/qu"
	"go.uber.org/atomic"
)

func LogConsume(
	quit qu.C, handler func(ent *Entry) (e error),
	filter func(pkg string) (out bool), args ...string,
) *Worker {
	log.D.Ln("starting log consumer")
	return Consume(
		quit, func(b []byte) (e error) {
			// we are only listening for entries
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "entr":
					var ent Entry
					n := gotiny.Unmarshal(b, &ent)
					log.D.Ln("consume", n)
					if filter(ent.Package) {
						// if the worker filter is out of sync this stops it printing
						return
					}
					switch ent.Level {
					case Fatal:
					case Error:
					case Warn:
					case Info:
					case Check:
					case Debug:
					case Trace:
					default:
						log.D.Ln("got an empty log entry")
						return
					}
					if e = handler(&ent); log.E.Chk(e) {
					}
				}
			}
			return
		}, args...,
	)
}

func Start(w *Worker) {
	log.D.Ln("sending start signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("run ")); n < 1 || log.E.Chk(e) {
		log.D.Ln("failed to write", w.Args)
	}
}

// Stop running the worker
func Stop(w *Worker) {
	log.D.Ln("sending stop signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("stop")); n < 1 || log.E.Chk(e) {
		log.D.Ln("failed to write", w.Args)
	}
}

// Kill sends a kill signal via the pipe logger
func Kill(w *Worker) {
	var e error
	if w == nil {
		log.D.Ln("asked to kill worker that is already nil")
		return
	}
	var n int
	log.D.Ln("sending kill signal")
	if n, e = w.StdConn.Write([]byte("kill")); n < 1 || log.E.Chk(e) {
		log.D.Ln("failed to write")
		return
	}
	if e = w.Cmd.Wait(); log.E.Chk(e) {
	}
	log.D.Ln("sent kill signal")
}

// SetLevel sets the level of logging from the worker
func SetLevel(w *Worker, level string) {
	if w == nil {
		return
	}
	log.D.Ln("sending set level", level)
	lvl := 0
	for i := range Levels {
		if level == Levels[i] {
			lvl = i
		}
	}
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("slvl" + string(byte(lvl)))); n < 1 ||
		log.E.Chk(e) {
		log.D.Ln("failed to write")
	}
}

// LogServe starts up a handler to listen to logs from the child process worker
func LogServe(quit qu.C, appName string) {
	log.D.Ln("starting log server")
	lc := AddLogChan()
	var logOn atomic.Bool
	logOn.Store(false)
	p := Serve(
		quit, func(b []byte) (e error) {
			// listen for commands to enable/disable logging
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "run ":
					log.D.Ln("setting to run")
					logOn.Store(true)
				case "stop":
					log.D.Ln("stopping")
					logOn.Store(false)
				case "slvl":
					log.D.Ln("setting level", Levels[b[4]])
					SetLogLevel(Levels[b[4]])
				case "kill":
					log.D.Ln(
						"received kill signal from pipe, shutting down",
						appName,
					)
					Request()
					quit.Q()
				}
			}
			return
		},
	)
	go func() {
	out:
		for {
			select {
			case <-quit.Wait():
				if !LogChanDisabled.Load() {
					LogChanDisabled.Store(true)
				}
				log.D.Ln("quitting pipe logger")
				Request()
				logOn.Store(false)
			out2:
				// drain log channel
				for {
					select {
					case <-lc:
						break
					default:
						break out2
					}
				}
				break out
			case ent := <-lc:
				if !logOn.Load() {
					break out
				}
				var n int
				var e error
				if n, e = p.Write(gotiny.Marshal(&ent)); !log.E.Chk(e) {
					if n < 1 {
						log.E.Ln("short write")
					}
				} else {
					break out
				}
			}
		}
		<-HandlersDone
		log.D.Ln("finished pipe logger")
	}()
}

// FilterNone is a filter that doesn't
func FilterNone(string) bool {
	return false
}

// SimpleLog is a very simple log printer
func SimpleLog(name string) func(ent *Entry) (e error) {
	return func(ent *Entry) (e error) {
		log.D.F(
			"%s[%s] %s %s",
			name,
			ent.Level,
			// ent.Time.Format(time.RFC3339),
			ent.Text,
			ent.CodeLocation,
		)
		return
	}
}
