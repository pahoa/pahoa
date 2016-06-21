package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pahoaCmd = &cobra.Command{Use: "pahoa"}

func init() {
	pahoaCmd.AddCommand(addCmd)
	pahoaCmd.AddCommand(moveCmd)
	pahoaCmd.AddCommand(listCmd)
	pahoaCmd.AddCommand(serverCmd)
}

func Execute() error {
	return pahoaCmd.Execute()
}

func initClientConfig(config *viper.Viper, cmd *cobra.Command) {
	config.AutomaticEnv()
	config.SetEnvPrefix("pahoa")

	config.SetConfigName(".pahoa")
	config.AddConfigPath("$HOME")
	config.AddConfigPath(".")

	cmd.PersistentPreRunE = func(*cobra.Command, []string) error {
		if err := config.ReadInConfig(); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("Failed to load configuration file: %#v", err)
		}

		if config.Get("endpoint") == "" {
			return fmt.Errorf("endpoint is required")
		}

		return nil
	}

	ps := cmd.PersistentFlags()
	ps.StringP("endpoint", "e", "", "pahoa url endpoint")

	config.BindPFlag("endpoint", ps.Lookup("endpoint"))
}
