package config

import (
	"database/sql"
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration settings
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	Db            *sql.DB
	InProduction  bool
	Session       *scs.SessionManager
}
