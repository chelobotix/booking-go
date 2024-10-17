package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/chelobotix/booking-go/internal/config"
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
	err := run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	_ = srv.ListenAndServe()
}

func run() error {
	gob.Register(models.Reservation{})

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

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cant create templateCache")
		return err
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = true

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	render.NewTemplates(&appConfig)

	helpers.NewHelpers(&appConfig)

	return nil
}
