package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/chelobotix/booking-go/internal/models"
	"html/template"
	"log"
)

// AppConfig hold the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	Production    bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	MailChan      chan models.MailData
}
