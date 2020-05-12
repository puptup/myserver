package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "form.html")
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("name")
		address := r.FormValue("address")
		token := name + ":" + address
		cookie := http.Cookie{
			Name:  "token",
			Value: token,
		}
		http.SetCookie(w, &cookie)
		http.ServeFile(w, r, "form.html")
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "login.html")
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		login := r.FormValue("login")
		password := r.FormValue("password")

		if login+":"+password == generateToken() {
			authorization = true
			fmt.Fprintf(w, "Successful authorization")
		} else {
			fmt.Fprintf(w, "Login or password entered incorrectly")
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

var login string = "Kirill"
var password string = "qwerty"
var authorization bool = false

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if authorization {
			endpoint(w, r)
		} else {
			fmt.Fprintf(w, "Please, enter your log:pass in /login/ ")
		}
	})
}

func generateToken() string {
	return login + ":" + password
}

func main() {
	fmt.Println("Server started")
	http.HandleFunc("/login", logHandler)
	http.Handle("/", isAuthorized(formHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
