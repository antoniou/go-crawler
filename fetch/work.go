package fetch

// Possible Worker states
const (
	WAITING uint8 = 0
	STOPPED uint8 = 1
	RUNNING uint8 = 2
)

// Worker is an interface that can be used
// to manage agents perform work in a different
// thread.
type Worker interface {
	// Run starts the Asynchronous worker
	Run() error

	// Returns worker name
	// Example names are:
	// - Fetcher
	// - Parser
	// - Tracker
	// - Sitemapper
	Type() string

	// State returns the state the worker is in:
	// RUNNING - processing work
	// WAITING - Waits for work
	// STOPPED - Not running
	State() uint8
	SetState(state uint8)
}

// AsyncWorker implements the worker interface
// It is meant to be embedded in another struct,
// like AsyncHttpFetcher
type AsyncWorker struct {
	RunFunc func() error

	state uint8
	quit  chan uint8
	Name  string
}

// Run calls the encapsulating
func (w *AsyncWorker) Run() error {
	return w.RunFunc()
}

// Stop notifies the quit channel
// The encapsulating struct's RunFunc
// needs to receive from the quit channel
// in order to stop.
func (w *AsyncWorker) Stop() {
	w.quit <- 0
}

// State getter (See interface definition)
func (w *AsyncWorker) State() uint8 {
	return w.state
}

// SetState setter (See interface definition)
func (w *AsyncWorker) SetState(state uint8) {
	w.state = state
}

// Type returns the Name given to the Worker
// in initialisation
func (w *AsyncWorker) Type() string {
	return w.Name
}
