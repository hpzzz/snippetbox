package main

import (
	"fmt"
	"net/http"
	"strconv"

	"karolharasim.com/snippetbox/pkg/forms"
	"karolharasim.com/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
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
	//retrieve the value and then delete the key and value from the session data; acts like one time fetch
	// flash := app.session.PopString(r, "flash")

	//pass the flash message to the template
	app.render(w, r, "show.page.tmpl", &templateData{
		// Flash: flash,
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
        // Pass a new empty forms.Form object to the template.
        Form: forms.New(nil),
    })


}


func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

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
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)


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
	//Validate the form contents using the form helper
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	//If there are any errors, redisplay the signup form
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return 
	}
	//Try ro create a new user record in the database. If the email already exists add an error message to the form and redisplay it
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return 
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	//Otherwise add a confirmation flash message to the session confirming that their signup worked and asking them to log in
	app.session.Put(r, "flash", "your signup was successful. Please log in")
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

	//Check whether the credentials are valid, of they're not, add a generic error message to the form failures map and re-display the login page
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form,})
		return 
	} else if err != nil {
		app.serverError(w, err)
		return 
	}

	//Add the ID of the current user to the session, so that they are now 'logged in'
	app.session.Put(r, "userID", id)

	//Redirect the user to the create snippet page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) { 
	// Remove the userID from the session data so that the user is 'logged out'.
    app.session.Remove(r, "userID")
    // Add a flash message to the session to confirm to the user that they've been logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!") //doesnt show up
	
    http.Redirect(w, r, "/", 303)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) handleAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API")
}