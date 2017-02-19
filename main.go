package main

import (
  "net/http"
  "html/template"
  "database/sql"
  _ "github.com/lib/pq"
  "golang.org/x/crypto/bcrypt"
  "fmt"
)

/* TODO: ERROR HANDLING
 * Sign Up: insert user, unique email, db error (userId=0)
 * Sign In: email doesn't exist (emailId =0), password is wrong, db error
 * Template issues:
 * DB open:
 * DB create table:
 */

var db *sql.DB
var userId int

func homeHandler(w http.ResponseWriter, r *http.Request) {
  if userId == 0 {
    loginT, _ := template.ParseFiles("template/login.html")
    loginT.Execute(w, nil)
  } else {
    fmt.Println("User is logged in")
  }
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
  name := r.FormValue("name")
  email := r.FormValue("email")
  password := []byte(r.FormValue("password"))
  hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

  db.QueryRow("INSERT INTO users(name, email, password) VALUES($1, $2, $3) returning uid;", name, email, string(hash)).Scan(&userId)
  http.Redirect(w, r, "/", http.StatusFound)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
  email := r.FormValue("email")
  password := []byte(r.FormValue("password"))
  fmt.Println(password)

  var emailId int
  var dbHash string
  db.QueryRow("SELECT uid,password FROM users WHERE email=$1", email).Scan(&emailId, &dbHash)
  if bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password)) == nil {
    userId = emailId;
  }
  http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
  db,_ = sql.Open("postgres", "user=postgres password=postgres dbname=Olson sslmode=disable")
  // If any changes are made to structure of database, DROP TABLE table
  db.Exec("CREATE TABLE IF NOT EXISTS users (uid serial NOT NULL, name VARCHAR (50) NOT NULL, email VARCHAR (50) NOT NULL UNIQUE, password VARCHAR (100) NOT NULL)");

  http.HandleFunc("/", homeHandler)
  http.HandleFunc("/signup", signupHandler)
  http.HandleFunc("/signin", signinHandler)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/Users/Olson/Documents/Go/src/github.com/molson194/Yowza/static"))))
  http.ListenAndServe(":8080", nil)
}
