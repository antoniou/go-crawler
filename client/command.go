package client

import "fmt"

// A Command is a runnable sub-command of a CLI.
type Command interface {
	Run(args []string) error
}

type BaseCommand struct {
	Name        string
	Description string
	Usage       string
	ArgsUsage   string
	Action      func(interface{}) error
}

func (c *BaseCommand) Run(args []string) error {
	return fmt.Errorf("Should not be implemented")
}
