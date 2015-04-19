package cli

import (
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/piotrkowalczuk/auth-service/service"
)

var (
	initDBCommand = cli.Command{
		Name:   "initdb",
		Usage:  "set up database",
		Action: initDBCommandAction,
	}
)

func initDBCommandAction(context *cli.Context) {
	service.InitConfig(context.GlobalString("environment"))
	service.InitLogger(service.Config.Logger)
	service.InitDB(service.Config.DB)

	queryBytes, err := ioutil.ReadFile("data/sql/schema_" + service.Config.DB.Adapter + ".sql")
	if err != nil {
		service.Logger.Fatal(err)
	}

	service.Logger.Info("Schema file opened successfully.")

	_, err = service.DBPool.Exec(string(queryBytes))
	if err != nil {
		service.Logger.Fatal(err)
	}

	service.Logger.Info("Database has been initialized successfully.")
}
