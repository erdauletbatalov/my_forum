package web

import (
	"html/template"
	"log"

	"github.com/erdauletbatalov/forum/pkg/models/sqlite"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Forum         *sqlite.ForumModel
	TemplateCache map[string]*template.Template
}
