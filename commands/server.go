package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
	"github.com/pahoa/pahoa/handlers"
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
	config := viper.New()

	config.AutomaticEnv()
	config.SetEnvPrefix("pahoa")

	config.SetConfigFile(serverOptions.cfgFile)

	if err := config.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to load configuration file: %s",
			serverOptions.cfgFile)
	}

	var board core.Board
	if err := config.UnmarshalKey("board", &board); err != nil {
		return fmt.Errorf("Failed to load board configuration")
	}

	model := &core.Model{}

	executor := core.NewExecutor(&core.NewExecutorOptions{
		Model:    model,
		Handlers: handlers.GetHandlers(),
		Config:   config,
	})
	executor.Start()

	s := server.NewServer(&server.NewServerOptions{
		Board:    &board,
		Model:    model,
		Executor: executor,
	})

	log.Printf("Starting server at: http://%s", serverOptions.address)
	return http.ListenAndServe(serverOptions.address, s)
}
