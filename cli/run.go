package cli

import (
	"html/template"
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/go-soa/auth/controller/web"
	"github.com/go-soa/auth/service"
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
	service.InitMailer(service.Config.Mailer)
	service.InitPasswordHasher(service.Config.PasswordHasher)

	router := httprouter.New()

	setupStaticRoutes(router)
	setupWebRoutes(router)

	listenOn := service.Config.Server.Host + ":" + service.Config.Server.Port
	service.Logger.Fatal(http.ListenAndServe(listenOn, router))
}

func getTemplatePath(path string) string {
	return service.Config.Template.Path + "/" + path
}

func setupWebRoutes(router *httprouter.Router) {
	templates, err := template.ParseFiles(
		getTemplatePath("header.html"),
		getTemplatePath("footer.html"),
		getTemplatePath("registration/index.html"),
		getTemplatePath("registration/success.html"),
	)
	if err != nil {
		service.Logger.Fatal(err)
	}

	registrationIndex := &web.Handler{
		TmplName:       "registration_index",
		Tmpl:           templates,
		Logger:         service.Logger,
		DB:             service.DBPool,
		RM:             service.RepositoryManager,
		PasswordHasher: service.PasswordHasher,
		Middlewares: web.NewMiddlewares(
			(*web.Handler).RegistrationIndex,
		),
		Mailer: service.Mail,
	}

	registrationSuccess := &web.Handler{
		TmplName:       "registration_success",
		Tmpl:           templates,
		Logger:         service.Logger,
		DB:             service.DBPool,
		RM:             service.RepositoryManager,
		PasswordHasher: service.PasswordHasher,
		Middlewares: web.NewMiddlewares(
			(*web.Handler).RegistrationSuccess,
		),
		Mailer: service.Mail,
	}

	registrationCreate := &web.Handler{
		TmplName:       "registration_index",
		Tmpl:           templates,
		Logger:         service.Logger,
		DB:             service.DBPool,
		RM:             service.RepositoryManager,
		PasswordHasher: service.PasswordHasher,
		Middlewares: web.NewMiddlewares(
			(*web.Handler).RegistrationCreate,
		),
		Mailer: service.Mail,
	}

	router.Handler("GET", "/registration", registrationIndex)
	router.Handler("GET", "/registration/success", registrationSuccess)
	router.Handler("POST", "/registration", registrationCreate)
}

func setupStaticRoutes(router *httprouter.Router) {
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
}
