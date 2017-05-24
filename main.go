package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

/* TODO: ERROR HANDLING
 * Sign Up: insert user, unique email, db error (userId=0)
 * Sign In: email doesn't exist (emailId =0), password is wrong, db error
 * Template issues:
 * DB open:
 * DB create table:
 */

type Profile struct {
	Name        string
	Email       string
	DisplayEdit string
	Statement   string
	Phone       string
	Location    string
	Summary     string
	Companies   string
	Skills      string
}

var db *sql.DB
var userId int
var prof *Profile

func loadProfile(uid int) *Profile {
	if userId == 0 {
		return nil
	} else {
		var name string
		var email string
		var statement string
		var phone string
		var location string
		var summary string
		var companies string
		// TODO: figure out array and how to range
		var skills string
		db.QueryRow("SELECT name,email,statement FROM users WHERE uid=$1", userId).Scan(&name, &email, &statement)
		db.QueryRow("SELECT phone,location,summary FROM users WHERE uid=$1", userId).Scan(&phone, &location, &summary)
		db.QueryRow("SELECT companies FROM users WHERE uid=$1", userId).Scan(&companies)
		fmt.Println(companies)
		db.QueryRow("SELECT skills FROM users WHERE uid=$1", userId).Scan(&skills)
		return &Profile{Name: name, Email: email, Statement: statement, Phone: phone, Location: location, Summary: summary, Companies: companies, Skills: skills}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if userId == 0 {
		loginT, _ := template.ParseFiles("template/login.html")
		loginT.Execute(w, nil)
	} else {
		prof = loadProfile(userId)
		prof.DisplayEdit = "" // NOTE: for other profiles "display:none"
		profT, _ := template.ParseFiles("template/profile.html")
		profT.Execute(w, prof)
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

	var emailId int
	var dbHash string
	db.QueryRow("SELECT uid,password FROM users WHERE email=$1", email).Scan(&emailId, &dbHash)
	if bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password)) == nil {
		userId = emailId
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if userId == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		editT, _ := template.ParseFiles("template/edit.html")
		editT.Execute(w, prof)
	}
}

func saveeditHandler(w http.ResponseWriter, r *http.Request) {
	if prof != nil && userId != 0 {
		r.ParseForm()
		name := r.FormValue("name")
		if prof.Name != name {
			db.Exec("UPDATE users SET name=$1 where uid=$2", name, userId)
		}
		email := r.FormValue("email")
		if prof.Email != email {
			db.Exec("UPDATE users SET email=$1 where uid=$2", email, userId)
		}
		// TODO: upload image to s3
		statement := r.FormValue("statement")
		if prof.Statement != statement {
			db.Exec("UPDATE users SET statement=$1 where uid=$2", statement, userId)
		}
		phone := r.FormValue("phone")
		if prof.Phone != phone {
			db.Exec("UPDATE users SET phone=$1 where uid=$2", phone, userId)
		}
		location := r.FormValue("location")
		if prof.Location != location {
			db.Exec("UPDATE users SET location=$1 where uid=$2", location, userId)
		}
		summary := r.FormValue("summary")
		if prof.Summary != summary {
			db.Exec("UPDATE users SET summary=$1 where uid=$2", summary, userId)
		}
		var companies []string
		for k, v := range r.Form {
			if k == "company" {
				companies = v
			}
		}
		// TODO: only if change
		db.Exec("UPDATE users SET companies=$1 where uid=$2", pq.Array(companies), userId)

		skills := r.FormValue("skills")
		if prof.Skills != skills {
			db.Exec("UPDATE users SET skills=$1 where uid=$2", skills, userId)
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	db, _ = sql.Open("postgres", "user=postgres password=postgres dbname=Olson sslmode=disable")
	// If any changes are made to structure of database, DROP TABLE table
	db.Exec("CREATE TABLE IF NOT EXISTS users (uid serial NOT NULL, name VARCHAR (50) NOT NULL, email VARCHAR (50) NOT NULL UNIQUE, password VARCHAR (100) NOT NULL, statement TEXT, phone VARCHAR(10), location TEXT, summary TEXT, companies text[], skills TEXT)")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/saveedit", saveeditHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/Users/Olson/Documents/Go/src/github.com/molson194/Yowza/static"))))
	http.ListenAndServe(":8080", nil)
}
