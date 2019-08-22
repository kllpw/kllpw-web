package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kllpw/kllpw-web/ascii"
	"github.com/kllpw/kllpw-web/render"
	"github.com/kllpw/kllpw-web/user"
	"log"
	"net/http"
	"os"
)

var userManager = user.NewManager(os.Getenv("SESSION_KEYS"))

func protection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if userManager.IsUserAuthenticated(w, r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
	return http.HandlerFunc(fn)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	render.WritePageToTemplate(w, render.Index, render.GetPageTemplate(render.Index))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if userManager.LoginUser(w, r) {
		render.WritePageToTemplate(w, render.Login, render.GetPageTemplate(render.Login))
	} else {
		userManager.LogoutUser(w, r)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func loginFormHandler(w http.ResponseWriter, r *http.Request) {
	render.WritePageToTemplate(w, render.Login, render.GetPageTemplate(render.Login))
}

func registerFormHandler(w http.ResponseWriter, r *http.Request) {
	render.WritePageToTemplate(w, render.Register, render.GetPageTemplate(render.Register))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if userManager.RegisterUser(w, r) {
		fmt.Fprint(w, "User Registered")
	} else {
		fmt.Fprint(w, "User Registration failed")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	userManager.LogoutUser(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func userHomeHandler(w http.ResponseWriter, r *http.Request) {
	user := userManager.GetUser(w, r)
	popUserHome := render.UserHome
	popUserHome.Header = ascii.RenderString(user.Name)
	popUserHome.ContentTitle = "details:"
	popUserHome.Content = map[string]interface{}{
		"User" : user,
	}
	render.WritePageToTemplate(w, popUserHome, render.GetPageTemplate(popUserHome))
}

func main() {
	fmt.Print("\n" + ascii.RenderString(" kll.pw"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./render/templates/static/"))))

	userRoute := r.PathPrefix("/user").Subrouter()
	userRoute.HandleFunc("/login", loginFormHandler)
	userRoute.HandleFunc("/logout", logoutHandler)
	userRoute.HandleFunc("/register", registerFormHandler)

	userRequest := userRoute.PathPrefix("/req").Subrouter()
	userRequest.HandleFunc("/login", loginHandler)
	userRequest.HandleFunc("/logout", logoutHandler)
	userRequest.HandleFunc("/register", registerHandler)

	protected := r.PathPrefix("/user").Subrouter()
	protected.Use(protection)
	protected.HandleFunc("/home", userHomeHandler)

	log.Println("Ready...")
	log.Fatal(http.ListenAndServeTLS("", os.Getenv("SSLCERT"), os.Getenv("SSLKEY"), r))
}
