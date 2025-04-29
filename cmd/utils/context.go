package utils

import (
	"context"

	"github.com/spf13/cobra"
)

func FindContext(cmd *cobra.Command) context.Context {
	ctx := cmd.Context()
	if ctx == nil {
		if cmd.Parent() != nil {
			return FindContext(cmd.Parent())
		} else {
			return nil
		}
	}

	return ctx
}
