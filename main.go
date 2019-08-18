package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kllpw/kllpw-web/ascii"
	"github.com/kllpw/kllpw-web/client"
	"log"
	"net/http"
	"os"
)

var clientManager = client.NewManager(os.Getenv("SESSION_KEYS"))

func protection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if clientManager.IsValidClient(w, r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
	return http.HandlerFunc(fn)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	header := ascii.RenderStringHTML(" kll.pw")
	fmt.Fprint(w,
		"<html>"+
			header+
			`
        <a href="/user">register</a>
        <a href="/user">login</a>
        </html>`)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if clientManager.LoginClient(w, r) {
		fmt.Fprint(w, "Successful")
	} else {
		clientManager.LogoutClient(w, r)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func loginFormHandler(w http.ResponseWriter, r *http.Request) {
	header := ascii.RenderStringHTML(" kll.pw")
	fmt.Fprint(w,
		`<html><header>
        <script>
        function authenticateUser(user, password)
        {
            var token = user + ":" + password;
            var hash = btoa(token); 

            return "Basic " + hash;
        }
        function requestAuthentication() {
            var username=document.getElementById("username").value;
            var password=document.getElementById("password").value;
            // New XMLHTTPRequest
            var request = new XMLHttpRequest();
            request.open("POST", "/user/login", false);
            request.setRequestHeader("Authorization", authenticateUser(username, password));  
			request.send();
			
			// view request status
			document.getElementById("response").innerHTML = request.responseText;
			if (request.status == "200"){
				window.location.href = '/user/home';
			}
            
        }
        function register() {
            var username=document.getElementById("username").value;
            var password=document.getElementById("password").value;
            
            // New XMLHTTPRequest
            var request = new XMLHttpRequest();
            request.open("POST", "/user/register", false);
            request.setRequestHeader("Authorization", authenticateUser(username, password));  
            request.send();
			// view request status
            document.getElementById("response").innerHTML = request.responseText;
			
        }
        </script>
        </header>`+
			header+
			`<div>
		<a href="/">home </a>
		<br>
        <label>Username:</label><input id="username" name="username"></input>
        <label>Password:</label><input type="password" id="password" name="password"></input>
        <input type="button" value="register"onclick="register();"></input>
		<input type="button" value="login" onclick="requestAuthentication();"></input>
		<br>
		<label>Response:</label><label id="response" name="response">Waiting....</label>
        </div>
        </html>`)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if clientManager.RegisterClient(w, r) {
		fmt.Fprint(w, "Client Registered")
	} else {
		fmt.Fprint(w, "Client Registration failed")
	}

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clientManager.LogoutClient(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func userHomeHandeler(w http.ResponseWriter, r *http.Request) {
	client := clientManager.GetClient(w, r)
	name := ascii.RenderString(client.Name)
	fmt.Fprint(w,
		"<html><pre>"+
		name +
		`</pre>
        <a href="/user/logout">logout</a>
        <a href="/user/home">home/dashboard</a>
        </html>`)
}

func main() {
	fmt.Print("\n" + ascii.RenderString(" kll.pw"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)

	userRoute := r.PathPrefix("/user").Subrouter()
	userRoute.HandleFunc("", loginFormHandler)
	userRoute.HandleFunc("/login", loginHandler)
	userRoute.HandleFunc("/logout", logoutHandler)
	userRoute.HandleFunc("/register", registerHandler)

	protected := r.PathPrefix("/user").Subrouter()
	protected.Use(protection)
	protected.HandleFunc("/home", userHomeHandeler)

	log.Println("Ready...")
	log.Fatal(http.ListenAndServe("localhost:8000", r))

}
