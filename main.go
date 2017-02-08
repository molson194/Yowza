package main

import (
  "net/http"
  "html/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("template/login.html")
  t.Execute(w, nil)
}

func main() {
  http.HandleFunc("/", homeHandler)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/Users/Olson/Documents/Go/src/github.com/molson194/Yowza/static"))))
  http.ListenAndServe(":8080", nil)
}
