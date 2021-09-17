package main

import (
	"encoding/gob"
	"flag"
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
	"github.com/dhanekom/bookings/internal/repository/dbrepo"
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

	defer close(app.MailChan)
	listenForMail()

	msg := models.MailData{
		To:      "john@do.ca",
		From:    "me@here.com",
		Subject: "Some subject",
		Content: "Hallo <strong>world</strong>",
	}

	app.MailChan <- msg

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

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	userCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database post")
	dbSSL := flag.String("dbssl", "disable", "Database sslsettings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required parameters")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.TemplatePath = "./templates"

	tc, err := render.CreateTemplateCache(app.TemplatePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create template cache - %s", err)
	}

	app.InProduction = *inProduction
	app.UseCache = *userCache

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
	// host=localhost port=5432 dbname=bookings user=pos password=pos
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	fmt.Println(connectionString)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cannot connet to database! Dying...")
	}

	log.Println("connected to database")

	app.TemplateCache = tc
	myDBRepo := dbrepo.NewPostgresRepo(db.SQL, &app)
	render.NewRendered(&app)
	handlers.NewRepo(&app, myDBRepo)
	helpers.NewHelpers(&app)

	return db, nil
}
