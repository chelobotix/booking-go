package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
)

// AppConfig hold the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	Production    bool
	Session       *scs.SessionManager
}
