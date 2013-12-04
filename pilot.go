package main

import (
  "fmt"
  "html/template"
  "net/http"
  "strconv"

  "github.com/coopernurse/gorp"
)

var dbmap *gorp.DbMap

func viewHandler(w http.ResponseWriter, r *http.Request) {
  idString := r.URL.Path[len("/view/"):]
  id, _ := strconv.ParseInt(idString, 10, 64)
  s, _ := loadSubmission(id)
  t, _ := template.ParseFiles("view.html")
  t.Execute(w, s)
}

func main() {
  dbmap = initDb()
  defer dbmap.Db.Close()

  fmt.Println("DB initialized")

  http.HandleFunc("/view/", viewHandler)
  fmt.Println("Listening to port 1759")
  http.ListenAndServe(":1759", nil)
}
