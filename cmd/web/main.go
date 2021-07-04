package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dhanekom/bookings/internal/config"
	"github.com/dhanekom/bookings/internal/handlers"
	"github.com/dhanekom/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("unable to create template cache", err)
	}

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	app.UseCache = false
	app.TemplateCache = tc
	render.NewTemplates(&app)
	handlers.NewRepo(&app)

	log.Printf("Starting server on port %s\n", portNumber)
	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}
