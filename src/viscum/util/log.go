package util

import (
  "fmt"
  "os"
)

// Prints the arguments to a given stream.
func log(file *os.File, args ...interface{}) {
  fmt.Fprintln(file, args...)
}

// Prints the arguments to stdout.
func Info(args ...interface{}) {
  log(os.Stdout, args...)
}

// Prints the arguments to stderr.
func Error(args ...interface{}) {
  log(os.Stderr, args...)
}

// Prints the arguments to stderr and exits with error code 1.
func Fatal(args ...interface{}) {
  Error(args...)
  os.Exit(1)
}
