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
	"github.com/dhanekom/bookings/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	log.Printf("Starting server on port %s\n", portNumber)
	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	app.TemplatePath = "./templates"

	tc, err := render.CreateTemplateCache(app.TemplatePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create template cache - %s", err)
	}

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=pos password=pos")
	if err != nil {
		log.Fatal("cannot connet to database! Dying...")
	}

	log.Println("connected to database")

	app.UseCache = false
	app.TemplateCache = tc
	render.NewRendered(&app)
	handlers.NewRepo(&app, db)
	helpers.NewHelpers(&app)

	return db, nil
}
