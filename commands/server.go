package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/mattes/migrate/driver/sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pahoa/pahoa/core"
	"github.com/pahoa/pahoa/handlers"
	"github.com/pahoa/pahoa/server"
)

var serverCmdConfig = viper.New()

var serverCmd = &cobra.Command{
	Use:  "server",
	RunE: serverRun,
}

var serverCmdFile string

func init() {
	pf := serverCmd.PersistentFlags()

	pf.StringVarP(&serverCmdFile, "config", "c", "./pahoa.yaml", "path to config file")
	pf.StringP("bind", "b", "127.0.0.1:5544",
		"interface and port to which the server will bind")

	initServerConfig(serverCmdConfig, serverCmd)

	serverCmdConfig.BindPFlag("bind", pf.Lookup("bind"))
}

func serverRun(cmd *cobra.Command, args []string) error {
	var board core.Board
	if err := serverCmdConfig.UnmarshalKey("board", &board); err != nil {
		return fmt.Errorf("Failed to load board configuration")
	}

	model, err := core.NewSQLModel("sqlite3", "./pahoa.db", "./migrations")
	if err != nil {
		return fmt.Errorf("Failed to initialize model: %#v", err)
	}

	executor := core.NewExecutor(&core.NewExecutorOptions{
		Model:    model,
		Handlers: handlers.GetHandlers(),
		Config:   serverCmdConfig,
	})
	executor.Start()

	s := server.NewServer(&server.NewServerOptions{
		Board:    &board,
		Model:    model,
		Executor: executor,
	})

	log.Printf("Starting server at: http://%s", serverCmdConfig.GetString("bind"))
	return http.ListenAndServe(serverCmdConfig.GetString("bind"), s)
}
