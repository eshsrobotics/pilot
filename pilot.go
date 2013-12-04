package main

import (
  "fmt"
  "net/http"
  "strconv"

  "github.com/coopernurse/gorp"
)

var dbmap *gorp.DbMap

func submissionHandler(w http.ResponseWriter, r *http.Request) {
  idString := r.URL.Path[len("/submission/"):]
  id, _ := strconv.ParseInt(idString, 10, 64)
  s, _ := loadSubmission(id)
  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div><div>%s</div>", s.Title, s.Author, s.Code)
}

func main() {
  dbmap = initDb()
  defer dbmap.Db.Close()

  fmt.Println("DB initialized")

  http.HandleFunc("/submission/", submissionHandler)
  http.ListenAndServe(":1759", nil)

  fmt.Println("Listening to port 1759")
}
