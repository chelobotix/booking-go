package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/chelobotix/booking-go/pkg/config"
	"github.com/chelobotix/booking-go/pkg/handlers"
	"github.com/chelobotix/booking-go/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {
	appConfig.Production = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.Production
	appConfig.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cant create templateCache")
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = true

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	render.NewTemplates(&appConfig)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	_ = srv.ListenAndServe()
}
