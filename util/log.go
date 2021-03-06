package util

import "log"

// CustomLog wraps around log.Logger
// that will only output if verbose
// is true
type CustomLog struct {
	verbose bool
}

var logger *CustomLog

//Logger returns the singleton logger
func Logger(verbose bool) *CustomLog {
	if logger == nil || logger.verbose != verbose {
		logger = &CustomLog{
			verbose: verbose,
		}
	}
	return logger
}

// Println wraps around log.Println but prints
// only when in verbose mode
func (c *CustomLog) Println(v ...interface{}) {
	if c.verbose {
		log.Println(v...)
	}
}

// Printf wraps around log.Printf but prints
// only when in verbose mode
func (c *CustomLog) Printf(format string, v ...interface{}) {
	if c.verbose {
		log.Printf(format, v...)
	}
}

// Println wraps around log.Println
// Prints only when verbose mode is true
func Println(v ...interface{}) {
	if logger == nil {
		Logger(false)
	}
	logger.Println(v...)
}

// Printf wraps around log.Printf
// Prints only when verbose mode is true
func Printf(format string, v ...interface{}) {
	if logger == nil {
		Logger(false)
	}
	logger.Printf(format, v...)
}
