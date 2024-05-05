package main

import (
	"context"
	"os"
	"os/signal"
)

type exitCode int

const (
	ExitCodeOK exitCode = iota
	ExitCodeError
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	code := run(ctx)
	defer exit(code)
}

func run(_ context.Context) exitCode {
	// Do something
	return ExitCodeOK
}

// exit is a wrapper of os.Exit.
func exit[T ~int](code T) {
	os.Exit(int(code))
}
