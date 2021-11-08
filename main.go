package main

import (
	"context"
	"flag"
	"fmt"
	"goliveview-starter/config"
	"goliveview-starter/views"
	"goliveview-starter/views/accounts"
	"goliveview-starter/views/app"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	glvc := glv.Websocket(
		&name,
		glv.EnableHTMLFormatting(),
		glv.DisableTemplateCache(),
		glv.EnableDebugLog(),
		glv.EnableWatch(),
	)

	// 404 and /
	r.NotFound(glvc.Handle(&views.NotfoundView{}))
	r.Handle("/", glvc.Handle(&views.LandingView{Auth: authnAPI}))
	// unauthenticated views
	accountViews := accounts.Views{Auth: authnAPI}
	r.Handle("/signup", glvc.Handle(accountViews.Signup()))
	r.Handle("/confirm/{token}", glvc.Handle(accountViews.Confirm()))
	r.Handle("/login", glvc.Handle(accountViews.Login()))
	r.Handle("/magic-login/{token}", glvc.Handle(accountViews.ConfirmMagic()))
	r.Handle("/forgot", glvc.Handle(accountViews.Forgot()))
	r.Handle("/reset/{token}", glvc.Handle(accountViews.Reset()))
	// authenticated routes
	r.Route("/account", func(r chi.Router) {
		r.Use(authnAPI.IsAuthenticated)
		r.Handle("/", glvc.Handle(accountViews.Settings()))
		r.Handle("/email/change/{token}", glvc.Handle(accountViews.ConfirmEmailChange()))
	})
	r.Route("/app", func(r chi.Router) {
		r.Use(authnAPI.IsAuthenticated)
		r.Handle("/", glvc.Handle(&app.DashboardView{Auth: authnAPI}))
	})

	// third party auth provider routes
	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		err := authnAPI.LoginWithProvider(w, r)
		if err != nil {
			log.Printf("LoginWithProvider err %v\n", err)
			http.Error(w, "not found", 404)
			return
		}
		redirectTo := "/app"
		from := r.URL.Query().Get("from")
		if from != "" {
			redirectTo = from
		}

		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	})

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		err := authnAPI.LoginProviderCallback(w, r, nil)
		if err != nil {
			log.Printf("LoginProviderCallback err %v\n", err)
			http.Error(w, "not found", 404)
			return
		}
		redirectTo := "/app"
		from := r.URL.Query().Get("from")
		if from != "" {
			redirectTo = from
		}

		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		acc, err := authnAPI.CurrentAccount(r)
		if err != nil {
			log.Println("err logging out ", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		acc.Logout(w, r)
	})

	// setup static assets handler
	workDir, _ := os.Getwd()
	public := http.Dir(filepath.Join(workDir, "./", "public", "assets"))
	staticHandler(r, "/static", public)

	// others
	staticFileHandler(r, "/robots.txt", filepath.Join(workDir, "./", "public", "robots.txt"))
	staticFileHandler(r, "/favicon.ico", filepath.Join(workDir, "./", "public", "favicon.ico"))

	// server
	fmt.Printf("listening on http://localhost:%d\n", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
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

func staticFileHandler(r chi.Router, pattern string, filename string) {
	r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}
