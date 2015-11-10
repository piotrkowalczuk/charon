package cli

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/piotrkowalczuk/charon/charond/controller/web"
	"github.com/piotrkowalczuk/charon/charond/service"
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
	service.InitMailers(service.Config.Mailer, service.TplManager)
	service.InitPasswordRecoverer(
		service.Logger,
		service.PasswordHasher,
		service.RepositoryManager.User,
		service.RepositoryManager.PasswordRecovery,
		service.PasswordRecoveryMailer,
	)

	router := httprouter.New()

	setupNotFoundRoute(router)
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
		ConfirmationMailer: service.RegistrationConfirmationMailer,
		TemplateManager:    service.TplManager,
		Routes:             service.Routes,
		URLGenerator:       service.URLGenerator,
		Mnemosyne:          service.Mnemosyne,
		PasswordRecoverer:  service.PasswordRecoverer,
	}

	handlers := []*web.Handler{
		web.NewHandler(web.HandlerOpts{
			Name:   "registration",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).RegistrationSuccess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "registration",
			Method: "POST",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).RegistrationProcess,
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
			Name:   "registration_confirmation",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).RegistrationConfirmation,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "logout",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).LogoutIndex,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "login",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).LoginIndex,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "login",
			Method: "POST",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).LoginProcess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "password_recovery",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).PasswordRecoveryIndex,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "password_recovery",
			Method: "POST",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).PasswordRecoveryProcess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "password_recovery_success",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).PasswordRecoverySuccess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "password_recovery_confirmation",
			Method: "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).PasswordRecoveryConfirmationIndex,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "password_recovery_confirmation",
			Method: "POST",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).PasswordRecoveryConfirmationProcess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:     "password_recovery_confirmation_success",
			Template: "password_recovery_success",
			Method:   "GET",
			Middlewares: web.NewMiddlewares(
				(*web.Handler).PasswordRecoveryConfirmationSuccess,
			),
			Container: container,
		}),
		web.NewHandler(web.HandlerOpts{
			Name:   "dashboard",
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

func setupNotFoundRoute(router *httprouter.Router) {
	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
	})
}
