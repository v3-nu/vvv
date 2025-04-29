package utils

import "github.com/spf13/cobra"

type CommandGroup struct {
	Command  *cobra.Command
	Children []func(*cobra.Command)
}

func (c *CommandGroup) RegisterChildren() {
	for _, child := range c.Children {
		child(c.Command)
	}
}

func (c *CommandGroup) Register(root *cobra.Command) {
	root.AddCommand(c.Command)
	c.RegisterChildren()
}
