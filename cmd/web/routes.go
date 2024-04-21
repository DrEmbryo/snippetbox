package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes () http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet) 
	
	fileServer := http.FileServer(http.Dir("./cmd/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}