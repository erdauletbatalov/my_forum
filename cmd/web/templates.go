package main

import (
	"fmt"
	"html/template" // новый импорт
	"path/filepath" // новый импорт

	"github.com/erdauletbatalov/forum.git/pkg/models"
)

type templateData struct {
	User *models.User
	// Users []*models.User
	Error     bool
	IsSession bool
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	fmt.Println(pages)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
