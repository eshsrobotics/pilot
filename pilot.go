package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/coopernurse/gorp"
)

var dbmap *gorp.DbMap

func renderTemplate(w http.ResponseWriter, tmpl string, s *Submission) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, s); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Path[len("/view/"):]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s, err := loadSubmission(id)
	if err != nil || s == nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "view", s)
}

func main() {
	dbmap = initDb()
	defer dbmap.Db.Close()

	fmt.Println("DB initialized")

	http.HandleFunc("/view/", viewHandler)
	fmt.Println("Listening to port 1759")
	http.ListenAndServe(":1759", nil)
}
