package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	"github.com/coopernurse/gorp"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Path[len("/view/"):]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s, err := loadSubmission(id)
	if err != nil || s == nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "view", s)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new", nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	author := r.FormValue("author")
	code := r.FormValue("code")
	s := newSubmission(title, author, code)
	if err := s.insert(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+strconv.FormatInt(s.Id, 10), http.StatusFound)
}

var templates = template.Must(template.ParseFiles("view.html", "new.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, s *Submission) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(new|save|view)/(|[a-zA-Z0-9]+)$")

func makeHandler(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		f(w, r)
	}
}

var dbmap *gorp.DbMap

func main() {
	dbmap = initDb()
	defer dbmap.Db.Close()

	fmt.Println("DB initialized")

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/new/", makeHandler(newHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	fmt.Println("Listening to port 1759")
	http.ListenAndServe(":1759", nil)
}
