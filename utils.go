package main

import (
	"chatter/datastore"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// errRedirect redirects users to a error page.
func errRedirect(w http.ResponseWriter, r *http.Request, msg string) {
	q := url.QueryEscape("msg=" + msg)
	http.Redirect(w, r, "/err?"+q, http.StatusNotFound)
}

// session verifies cookies validation against all private html pages.
func session(w http.ResponseWriter, r *http.Request) (s *datastore.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		return
	}

	s = &datastore.Session{UUID: cookie.Value}
	if ok, err := s.Check(); err != nil || !ok {
		err = errors.New("Invalid session")
	}
	return
}

// renderHTML parses necessary template files and renders a layout page to satisfy http handlers.
func renderHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, f := range filenames {
		files = append(files, fmt.Sprintf("templ/%s.html", f))
	}

	t := template.Must(template.ParseFiles(files...))
	err := t.ExecuteTemplate(w, "Layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// templates creates a layout template with other necessary templats.
// It returns a prepared layout template.
func templates(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, f := range filenames {
		files = append(files, fmt.Sprintf("templ/%s.html", f))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}
