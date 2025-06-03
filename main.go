package main

import (
	"fmt"
	"os"

	"github.com/charankamal20/gator/internal/config"
)

type state struct {
	conf *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	dict map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.dict[cmd.name]; exists {
		return handler(s, cmd)
	}

	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	if c.dict == nil {
		c.dict = make(map[string]func(*state, command) error)
	}

	c.dict[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}

	s.conf.SetUser(cmd.args[0])
	fmt.Printf("User set to %s\n", cmd.args[0])

	return nil
}

func main() {
	conf := config.Read()
	currState := &state{
		conf: &conf,
	}

	allCommands := &commands{
		dict: make(map[string]func(*state, command) error, 0),
	}

	allCommands.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) < 2 {
		os.Exit(1)
	}

	cmd := &command{
		name: args[0],
		args: args[1:],
	}

	allCommands.run(currState, *cmd)
}
