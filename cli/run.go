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
	service.InitRepositoryManager(service.DBPool)
	service.InitPasswordHasher(service.Config.PasswordHasher)
	service.InitTranslation(service.Config.Translation)
	service.InitRouting(service.Config.Routing)
	service.InitTemplates(service.Config.Templates, service.URLGenerator)
	service.InitMailer(service.Config.Mailer, service.MailTemplates)

	router := httprouter.New()

	setupStaticRoutes(router)
	setupWebRoutes(router)

	listenOn := service.Config.Server.Host + ":" + service.Config.Server.Port
	service.Logger.Fatal(http.ListenAndServe(listenOn, router))
}

func setupWebRoutes(router *httprouter.Router) {
	container := web.ServiceContainer{
		Logger:             service.Logger,
		DB:                 service.DBPool,
		RM:                 service.RepositoryManager,
		PasswordHasher:     service.PasswordHasher,
		ConfirmationMailer: service.ConfirmationMailer,
		Templates:          service.WebTemplates,
		Routes:             service.Routes,
		URLGenerator:       service.URLGenerator,
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
				(*web.Handler).RegistrationCreate,
			),
			Container: container,
		}),
	}

	for _, handler := range handlers {
		method := handler.Method
		routeName := handler.RouteName()
		pattern := service.Routes.GetPattern(handler.RouteName()).String()

		router.Handler(method, pattern, handler)

		service.Logger.WithFields(logrus.Fields{
			"name":    routeName,
			"method":  method,
			"pattern": pattern,
		}).Info("Route has been initialized successfully.")
	}
}

func setupStaticRoutes(router *httprouter.Router) {
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
}
