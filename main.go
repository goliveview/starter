package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"starter/config"
	"starter/views"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/adnaan/authn"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"

	glv "github.com/goliveview/controller"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()
	// load config
	configFile := flag.String("config", "env.local", "path to config file")
	envPrefix := os.Getenv("ENV_PREFIX")
	if envPrefix == "" {
		envPrefix = "app"
	}
	flag.Parse()
	cfg, err := config.Load(*configFile, envPrefix)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(cfg)

	// setup authn api
	authnAPI := authn.New(ctx, authn.Config{
		Driver:        cfg.Driver,
		Datasource:    cfg.DataSource,
		SessionSecret: cfg.SessionSecret,
		SendMail:      config.SendEmailFunc(cfg),
		GothProviders: []goth.Provider{
			google.New(
				cfg.GoogleClientID,
				cfg.GoogleSecret,
				fmt.Sprintf("%s/auth/callback?provider=google", cfg.Domain),
				"email", "profile",
			),
		},
	})

	// setup router
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat(cfg.HealthPath))
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	// create liveview controller and set routes
	name := "goliveview-starter"
	glvc := glv.Websocket(&name,
		glv.EnableHTMLFormatting(),
		glv.DisableTemplateCache(),
		glv.EnableDebugLog(),
		glv.EnableWatch(),
	)

	r.NotFound(glvc.NewView("./templates/404.html",
		glv.WithLayout("./templates/layouts/error.html")))

	landingLayout := glv.WithLayout("./templates/layouts/landing.html")

	r.Handle("/", glvc.NewView(
		"./templates/views/landing",
		landingLayout,
		glv.WithViewHandler(&views.HandlerLandingView{Auth: authnAPI})))

	r.Handle("/signup", glvc.NewView(
		"./templates/views/accounts/signup",
		landingLayout,
		glv.WithViewHandler(&views.HandlerSignupView{Auth: authnAPI})))

	r.Handle("/confirm/{token}",
		glvc.NewView("./templates/views/accounts/confirm", landingLayout,
			glv.WithViewHandler(&views.HandlerConfirmView{Auth: authnAPI})))

	loginView := &views.HandlerLoginView{Auth: authnAPI}
	r.Handle("/login", glvc.NewView("./templates/views/accounts/login",
		landingLayout,
		glv.WithOnPost(loginView.OnPost),
		glv.WithViewHandler(loginView)))

	// setup static assets handler
	workDir, _ := os.Getwd()
	public := http.Dir(filepath.Join(workDir, "./", "public", "assets"))
	staticHandler(r, "/static", public)

	fmt.Println("listening on http://localhost:4000")
	err = http.ListenAndServe(":4000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func staticHandler(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
