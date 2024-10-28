package supervisor

import (
	"log/slog"
)

// `Serve` runs `serve` for an indeterminate amount of time, in a goroutine, restarting
// it whenever it exits. The goroutine itself `recover`s from any recoverable `panic`s,
// logging an error.
// Caller is expected to pass any parameters by capture, including input and output
// channels it may need.
func Serve(serve func(), serviceName string) {
	// `defer` only runs when the function is exiting, so we need to wrap the service
	// to allow its caller (the supervisor) to keep running after recovery
	supervise := func() {
		// Recovery logic needs to run on `defer`, as it's the only code that runs during
		// a `panic`
		defer func() {
			panicked := recover()
			if panicked != nil {
				slog.Error(serviceName+": recovered from panic", "message", panicked)
			}
		}()
		serve()
	}
	// Now we fire up the goroutine with our service managed by a supervisor
	go func() {
		for {
			slog.Info(serviceName+": starting")
			supervise()
		}
	}()
}

// `OneShot` runs `task` only once under a supervised goroutine.
// If the function `panic`s the `panic` value will be reported on the `panicCh` channel
// and it will be logged.
// Otherwise, the channel will be closed on exit.
// Caller is expected to always wait on `panicCh` to ensure termination.
// Caller is expected to pass any parameters by capture, including channels for output
// of non-`panic` errors and results.
// If the caller doesn't care about the result, `panicCh` may be `nil`. The `panic`s
// will still be caught and logged, but not sent.
// TODO: we might still want to give `task` some kind of id.
func OneShot(task func(), panicCh chan<- any) {
	// Now we fire up the goroutine with our service managed by a supervisor
	go func() {
		// Recovery logic needs to run on `defer`, as it's the only code that runs during
		// a `panic`
		defer func() {
			panicked := recover()
			if panicked == nil {
				if panicCh != nil {
					close(panicCh)
				}
				return
			}
			slog.Error("oneshot task: recovered from panic", "message", panicked)
			if panicCh != nil {
				panicCh <- panicked
				close(panicCh)
			}
		}()
		task()
	}()
}
