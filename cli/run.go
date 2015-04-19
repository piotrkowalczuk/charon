package cli

import (
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/julienschmidt/httprouter"
	"github.com/piotrkowalczuk/auth-service/controller/web"
	"github.com/piotrkowalczuk/auth-service/service"
)

var (
	runCommand = cli.Command{
		Name:   "run",
		Usage:  "starts server",
		Action: runCommandAction,
	}
)

func runCommandAction(context *cli.Context) {
	service.InitConfig(context.GlobalString("environment"))
	service.InitLogger(service.Config.Logger)
	service.InitDB(service.Config.DB)

	registrationGET := &web.Handler{
		Logger: service.Logger,
		DB:     service.DBPool,
		Middlewares: web.NewMiddlewares(
			(*web.Handler).RegistrationGET,
		),
	}
	router := httprouter.New()
	router.Handler("GET", "/registration", registrationGET)

	listenOn := service.Config.Server.Host + ":" + service.Config.Server.Port
	service.Logger.Fatal(http.ListenAndServe(listenOn, router))
}
