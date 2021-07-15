package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dhanekom/bookings/internal/config"
	"github.com/dhanekom/bookings/internal/handlers"
	"github.com/dhanekom/bookings/internal/helpers"
	"github.com/dhanekom/bookings/internal/models"
	"github.com/dhanekom/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on port %s\n", portNumber)
	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}

func run() error {
	gob.Register(models.Reservation{})

	app.TemplatePath = "./templates"

	tc, err := render.CreateTemplateCache(app.TemplatePath)
	if err != nil {
		return fmt.Errorf("unable to create template cache - %s", err)
	}

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\n", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\n", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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
	helpers.NewHelpers(&app)

	return nil
}
