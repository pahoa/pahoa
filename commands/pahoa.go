package commands

import (
	"github.com/spf13/cobra"
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

var clientOptions = struct {
	endpoint string
}{}

func initClientCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&clientOptions.endpoint,
		"endpoint",
		"e",
		"",
		"pahoa endpoint")
}
