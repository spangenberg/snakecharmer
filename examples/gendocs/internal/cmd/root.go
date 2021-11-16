package cmd

import (
	"github.com/spf13/cobra"

	"github.com/spangenberg/snakecharmer"
)

func Execute() {
	snakecharmer.Execute(NewCmdRoot())
}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gendocs",
	}

	return cmd
}
