package config

import (
	"log"
	"text/template"
)

// AppConfig holds the application configuration settings
type AppConfig struct {
	UseCache bool

	TemplateCache map[string]*template.Template

	InfoLog *log.Logger
}
