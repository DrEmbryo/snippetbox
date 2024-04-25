package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DrEmbryo/snippetbox/cmd/pkg/forms"
	"github.com/DrEmbryo/snippetbox/cmd/pkg/models"
)

func (app *application) home (w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
}

func (app *application) showSnippet (w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get(":id")) 
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{Snippet: s})
}

func (app *application) createSnippet (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
		
	title := form.Get("title")
	content := form.Get("content")
	expires := form.Get("expires")

	expiresInt, err := strconv.Atoi(expires)
	if err != nil {
		app.serverError(w, err)
		return 
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, expiresInt)

	id, err := app.snippets.Insert(title, content, expiresAt.String());
	if err != nil {
		app.serverError(w, err)
		return 
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm (w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return 
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
	}

	name := form.Get("name")
	email := form.Get("email")
	password := form.Get("password")

	_, err = app.users.Insert(name, email, password)
	if err != nil {
		app.serverError(w, err)
		return	
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
    })		
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
	app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	email := form.Get("email")
	password := form.Get("password")

	id, err := app.users.Authenticate(email, password)
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You've been logged out seccessfully!")
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
