package cli

import (
	"os"

	"github.com/codegangsta/cli"
)

func init() {
	app := cli.NewApp()
	app.Name = "auth-service"
	app.Usage = "..."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "environment, e",
			Value:  "development",
			Usage:  "environment in wich application is running",
			EnvVar: "AUTH_SERVICE_ENV",
		},
	}
	app.Commands = []cli.Command{
		runCommand,
		initDBCommand,
	}

	app.Run(os.Args)
}
