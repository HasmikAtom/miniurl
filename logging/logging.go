package logging

import (
	"log"
	"os"
)

var (
	// Failsafe defaults to avoid panic when logging is used before InitResources
	Trace          = log.New(os.Stdout, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info           = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning        = log.New(os.Stderr, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error          = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Debugging      = log.New(os.Stdout, "[DEBUGGING] ", log.Ldate|log.Ltime|log.Lshortfile)
	TracingEnabled bool
)
