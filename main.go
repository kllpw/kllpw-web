package main

import (
    "fmt"
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "./client"
)

func protection(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        if client.IsValidClient(w, r) {
            next.ServeHTTP(w, r)
        } else {
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    }
	return http.HandlerFunc(fn)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
   client.LoginClient(w, r)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    client.RegisterClient(w, r)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    client.LogoutClient(w, r)
}

func handlerDash(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello authed client")
}

func main() {
      r := mux.NewRouter()
    r.HandleFunc("/", indexHandler)
	
	userRoute := r.PathPrefix("/user").Subrouter()
	userRoute.HandleFunc("/login", loginHandler)
    userRoute.HandleFunc("/logout", logoutHandler)
    userRoute.HandleFunc("/register", registerHandler)

	protected := r.PathPrefix("/admin").Subrouter()
	protected.Use(protection)
    protected.HandleFunc("/dashboard", handlerDash)


    log.Println("Ready...")
    log.Fatal(http.ListenAndServe("localhost:8000", r))
    
}