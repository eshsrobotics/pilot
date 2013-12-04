package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	"github.com/coopernurse/gorp"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var s []Submission
	_, err := dbmap.Select(&s, "select * from submissions order by Id desc")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "index", s)
}

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

func assetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[len("/"):])
}

var templates = template.Must(template.ParseFiles(
	"application.html",
	"index.html",
	"view.html",
	"new.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, i interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", i); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(|(new|save|view)/(|[a-zA-Z0-9]+))$")

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

	var port = flag.Int("port", 1759, "port to run server on")
	flag.Parse()

	http.HandleFunc("/", makeHandler(rootHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/new/", makeHandler(newHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/pilot.css", assetHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
