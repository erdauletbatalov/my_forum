package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"forum/pkg/models/sqlite"
	"forum/web"

	_ "github.com/mattn/go-sqlite3"
)

const Post = 0
const Comment = 0

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

	templateCache, err := web.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &web.Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Forum:         &sqlite.ForumModel{DB: db},
		TemplateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
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
	err = setup(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setup(db *sql.DB) error {
	fmt.Println("setup processing")
	query, err := ioutil.ReadFile("./pkg/models/sqlite/setup.sql")
	fmt.Println("setup processing")
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	fmt.Println("setup processing")
	if _, err := db.Exec(string(query)); err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	fmt.Println("setup success")
	return nil
}
