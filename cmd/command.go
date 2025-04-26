package cmd

import "github.com/urfave/cli/v2"

// Command defines the interface that all commands must implement
type Command interface {
	// Name returns the name of the command
	Name() string

	// Usage returns the usage description of the command
	Usage() string

	// Flags returns the command's flags
	Flags() []cli.Flag

	// Action executes the command with the given context
	Action(*cli.Context) error
}
