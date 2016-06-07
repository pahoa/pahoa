package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/server"
)

var serverCmd = &cobra.Command{
	Use:  "server",
	RunE: serverRun,
}

var serverOptions = struct {
	cfgFile string
	address string
}{}

func init() {
	serverCmd.PersistentFlags().StringVarP(
		&serverOptions.cfgFile,
		"config",
		"c",
		"./pahoa.yaml",
		"path to config file")
	serverCmd.PersistentFlags().StringVarP(
		&serverOptions.address,
		"bind",
		"b",
		"127.0.0.1:5544",
		"interface and port to which the server will bind")
}

func serverRun(cmd *cobra.Command, args []string) error {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("pahoa")

	v.SetConfigFile(serverOptions.cfgFile)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to load configuration file: %s",
			serverOptions.cfgFile)
	}

	s := server.NewServer()

	log.Printf("Starting server at: http://%s", serverOptions.address)
	return http.ListenAndServe(serverOptions.address, s)
}
