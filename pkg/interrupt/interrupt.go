package interrupt

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var requested bool

// InterruptChan is used to receive SIGINT (Ctrl+C) signals.
var InterruptChan chan os.Signal

// InterruptSignals is the list of signals that cause the interrupt
var InterruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

// ShutdownRequestChan is a channel that can receive shutdown requests
var ShutdownRequestChan = make(chan struct{})

// AddHandlerChannel is used to add an interrupt handler to the list of handlers to be invoked on SIGINT (Ctrl+C) signals.
var AddHandlerChannel = make(chan func())

// HandlersDone is closed after all interrupt handlers run the first time an interrupt is signaled.
var HandlersDone = make(chan struct{})

// Listener listens for interrupt signals, registers interrupt callbacks, and responds to custom shutdown signals as required
func Listener(
	) {
	var interruptCallbacks []func()
	invokeCallbacks := func() {
		// run handlers in LIFO order.
		for i := range interruptCallbacks {
			idx := len(interruptCallbacks) - 1 - i
			interruptCallbacks[idx]()
		}
		close(HandlersDone)
	}

	for {
		select {
		case sig := <-InterruptChan:
			fmt.Printf("received signal (%s) - shutting down...\n", sig)
			_ = sig
			requested = true
			invokeCallbacks()
			return
		case <-ShutdownRequestChan:
			fmt.Println("received shutdown request - shutting down...")
			requested = true
			invokeCallbacks()
			return
		case handler := <-AddHandlerChannel:
			interruptCallbacks = append(interruptCallbacks, handler)
		}
	}

}

// AddHandler adds a handler to call when a SIGINT (Ctrl+C) is received.
func AddHandler(
	handler func()) {
	// Create the channel and start the main interrupt handler which invokes all other callbacks and exits if not already done.
	if InterruptChan == nil {
		InterruptChan = make(chan os.Signal, 1)
		signal.Notify(InterruptChan, InterruptSignals...)
		go Listener()
	}

	AddHandlerChannel <- handler
}

// Request programatically requests a shutdown
func Request(
	) {
	close(ShutdownRequestChan)
}

// Requested returns true if an interrupt has been requested
func Requested(
	) bool {
	return requested
}
