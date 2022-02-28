package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	sqlite "github.com/erdauletbatalov/forum.git/pkg/models/sqlite"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	forum         *sqlite.ForumModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8080", "Сетевой адрес веб-сервера")
	dsn := flag.String("dsn", "./forum.db", "Название SQLite3 источника данных")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// // Инициализируем новый кэш шаблона...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// И добавляем его в зависимостях нашего
	// веб-приложения.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		forum:         &sqlite.ForumModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на http://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	setup(db)
	return db, nil
}

func setup(db *sql.DB) error {
	query, err := ioutil.ReadFile("./pkg/models/sqlite/setup.sql")
	if err != nil {
		return fmt.Errorf("setup: %s", err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		return fmt.Errorf("setup: %s", err)
	}
	return nil
}
