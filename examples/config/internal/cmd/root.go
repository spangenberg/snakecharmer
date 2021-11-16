package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spangenberg/snakecharmer"
	"github.com/spangenberg/snakecharmer/examples/config/internal"
)

func Execute() {
	snakecharmer.ExecuteWithConfig(NewCmdRoot(), "/etc/snakecharmer", "snakecharmer")
}

func NewCmdRoot() *cobra.Command {
	cfg := new(internal.Config)
	cmd := &cobra.Command{
		Use: "flags",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.Unmarshal(&cfg)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%#v", cfg)
		},
	}
	snakecharmer.GenerateFlags(cmd, cfg)
	return cmd
}
