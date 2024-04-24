package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DrEmbryo/snippetbox/cmd/pkg/models/db"
	"github.com/golangcollege/sessions"
	_ "modernc.org/sqlite"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	session *sessions.Session
	snippets *db.SnippetModel
	templateCache map[string]*template.Template
	users *db.UserModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	database, err := openDB()
	if err != nil {
	errorLog.Fatal(err)
	}
	
	defer database.Close()

	templateCache, err := newTemplateCache("./cmd/ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	secret := flag.String("secret","s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "secret")
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		session: session,
		snippets: &db.SnippetModel{DB: database},
		templateCache: templateCache,
		users: &db.UserModel{DB: database},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("listening on port %s", *addr)
	err = server.ListenAndServeTLS("./cmd/tls/cert.pem", "./cmd/tls/key.pem")
	errorLog.Fatal(err)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite","./cmd/db/data.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}