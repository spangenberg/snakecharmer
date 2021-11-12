package snakecharmer

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ExecuteWithConfig(cmd *cobra.Command, configPath, prefix string, cfgFile *string, verbose *bool) {
	cobra.OnInitialize(initConfig(cmd, configPath, prefix, cfgFile, verbose))

	Execute(cmd)
}

func HandleError(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	}
}

func Validate(f func(cmd *cobra.Command, args []string) (Config, error)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		i, err := f(cmd, args)
		if err != nil {
			return err
		}
		return validate(i)
	}
}

func initConfig(cmd *cobra.Command, configPath, prefix string, cfgFile *string, verbose *bool) func() {
	return func() {
		if *cfgFile != "" {
			viper.SetConfigFile(*cfgFile)
		} else if configPath != "" {
			viper.AddConfigPath(configPath)
			viper.SetConfigName("config")
		}

		viper.AutomaticEnv()
		viper.SetEnvPrefix(prefix)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

		if err := viper.ReadInConfig(); err == nil && *verbose {
			cmd.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
