package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"starter/views"
	"strings"

	glv "github.com/goliveview/controller"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//ctx := context.Background()
	//cfg := config.Config{}

	// setup authn api
	//authnConfig := authn.Config{
	//	Driver:        cfg.Driver,
	//	Datasource:    cfg.DataSource,
	//	SessionSecret: cfg.SessionSecret,
	//	SendMail:      config.SendEmailFunc(cfg),
	//	GothProviders: []goth.Provider{
	//		google.New(
	//			cfg.GoogleClientID,
	//			cfg.GoogleSecret,
	//			fmt.Sprintf("%s/auth/callback?provider=google", cfg.Domain),
	//			"email", "profile",
	//		),
	//	},
	//}
	//authnAPI := authn.New(ctx, authnConfig)
	//log.Println(authnAPI)

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(middleware.StripSlashes)

	name := "goliveview-starter"

	glvc := glv.Websocket(&name,
		glv.EnableHTMLFormatting(),
		glv.DisableTemplateCache(),
		glv.EnableDebugLog(),
		glv.EnableWatch(),
	)
	r.NotFound(glvc.NewView("./templates/404.html",
		glv.WithLayout("./templates/layouts/error.html")))
	r.Handle("/", glvc.NewView(
		"./templates/views/landing",
		glv.WithLayout("./templates/layouts/landing.html"),
		glv.WithViewHandler(&views.HandlerLandingView{})))
	r.Handle("/signup", glvc.NewView(
		"./templates/views/accounts/signup",
		glv.WithLayout("./templates/layouts/landing.html"),
		glv.WithViewHandler(&views.HandlerSignupView{})))
	r.Handle("/login", glvc.NewView(
		"./templates/views/accounts/login",
		glv.WithLayout("./templates/layouts/landing.html"),
		glv.WithViewHandler(&views.HandlerLoginView{})))

	workDir, _ := os.Getwd()
	public := http.Dir(filepath.Join(workDir, "./", "public", "assets"))
	staticHandler(r, "/static", public)

	fmt.Println("listening on http://localhost:4000")
	err := http.ListenAndServe(":4000", r)
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
