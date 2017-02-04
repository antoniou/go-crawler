package client

import (
	"log"

	"github.com/urfave/cli"
)

//Client represents a Mistri command line client
type Client struct {
	app      *cli.App
	commands map[string]Command
}

func (c *Client) RegisterCommand(name string, cmd Command) error {
	_, registered := c.commands[name]
	if registered {
		return clientError{s: "Command %s already registered"}
	}
	c.commands[name] = cmd
	return nil
}

func (c *Client) Run(arguments []string) error {
	return c.app.Run(arguments)
}

func (c *Client) RegisterCommands() error {
	if c.commands == nil {
		c.commands = make(map[string]Command)
	}
	err := c.RegisterCommand("crawl", NewCrawlCommand())
	if err != nil {
		log.Print(err)
	}
	return err
}

func New() (c *Client) {
	c = &Client{}
	c.RegisterCommands()

	app := cli.NewApp()
	app.Name = "go-crawler"
	app.Usage = ""
	app.Commands = []cli.Command{}
	for name, cmd := range c.commands {
		app.Commands = append(app.Commands, cli.Command{
			Name: name,
			Action: func(c *cli.Context) error {
				return cmd.Run(c.Args())
			},
		})
	}
	c.app = app
	return c
}

// clientError is an error used to signal different error situations in command handling.
type clientError struct {
	s         string
	userError bool
}

func (c clientError) Error() string {
	return c.s
}
