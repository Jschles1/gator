package main

import "fmt"

type command struct {
	name      string
	arguments []string
}

type commands struct {
	entries map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.entries[name]; !ok {
		c.entries[name] = f
	}
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.entries[cmd.name]; !ok {
		return fmt.Errorf("error: invalid command")
	}
	if err := c.entries[cmd.name](s, cmd); err != nil {
		return err
	}
	return nil
}
