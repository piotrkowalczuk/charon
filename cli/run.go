package cli

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/go-soa/charon/controller/web"
	"github.com/go-soa/charon/service"
	"github.com/julienschmidt/httprouter"
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
	service.InitMnemosyne(service.Config.Mnemosyne)
	service.InitRepositoryManager(service.DBPool)
	service.InitPasswordHasher(service.Config.PasswordHasher)
	service.InitTranslation(service.Config.Translation)
	service.InitRouting(service.Config.Routing)
	service.InitTemplateManager(service.Config.Templates)
	service.InitMailer(service.Config.Mailer, service.TplManager)

	router := httprouter.New()

	setupStaticRoutes(router)
	setupWebRoutes(router)

	host := service.Config.Server.Host
	port := service.Config.Server.Port
	tls := service.Config.Server.TLS

	service.Logger.WithFields(logrus.Fields{
		"tls":  tls,
		"host": host,
		"port": port,
	}).Info("HTTP(S) server is going to start.")
	if tls {
		service.Logger.Fatal(http.ListenAndServeTLS(
			host+":"+port,
			service.Config.Server.CertFile,
			service.Config.Server.KeyFile,
			router,
		))
	} else {
		service.Logger.Fatal(http.ListenAndServe(host+":"+port, router))
	}
}

func setupWebRoutes(router *httprouter.Router) {
	container := web.ServiceContainer{
		Config:             service.Config,
		Logger:             service.Logger,
		DB:                 service.DBPool,
		RM:                 service.RepositoryManager,
		PasswordHasher:     service.PasswordHasher,
		ConfirmationMailer: service.ConfirmationMailer,
		TemplateManager:    service.TplManager,
		Routes:             service.Routes,
		URLGenerator:       service.URLGenerator,
		Mnemosyne:          service.Mnemosyne,
	}

	handlers := []*web.Handler{
		web.NewHandler(web.HandlerOpts{
			Name:   "registration_index",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).RegistrationSuccess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "registration_success",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).RegistrationSuccess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "registration_index",
			Method: "POST",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).RegistrationProcess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "registration_confirmation",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).RegistrationConfirmation,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "login_index",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).LoginIndex,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "login_index",
			Method: "POST",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).LoginProcess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "dashboard_index",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).IsAuthenticatedMiddleware,
				(*web.Handler).DashboardIndex,
			),
			Container: container,
		}),
	}

	for _, handler := range handlers {
		handler.Register(router)
	}
}

func setupStaticRoutes(router *httprouter.Router) {
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
}
