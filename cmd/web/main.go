package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/chelobotix/booking-go/internal/config"
	"github.com/chelobotix/booking-go/internal/driver"
	"github.com/chelobotix/booking-go/internal/handlers"
	"github.com/chelobotix/booking-go/internal/helpers"
	"github.com/chelobotix/booking-go/internal/models"
	"github.com/chelobotix/booking-go/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(appConfig.MailChan)
	listenForMail()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	_ = srv.ListenAndServe()
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	appConfig.MailChan = mailChan

	appConfig.Production = false
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.Production
	appConfig.Session = session

	// connect to DB
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=booking user=x5 password=")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cant create templateCache")
		return nil, err
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = true

	repo := handlers.NewRepo(&appConfig, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&appConfig)

	helpers.NewHelpers(&appConfig)

	return db, nil
}
