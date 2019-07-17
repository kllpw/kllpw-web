package client

import (
	"log"
	"net/http"
	sess "./sess" 
	cred "./cred"
)

var sessManager =  sess.NewManager()
var credManager =  cred.NewManager()

// IsValidClient checks for stored session key if non found checks basic auth credentials
func IsValidClient(w http.ResponseWriter, r *http.Request) bool {
	return sessManager.IsClientAuthed(w, r)
}

// RegisterClient stores provided username and password for checks later
func RegisterClient(w http.ResponseWriter, r *http.Request) {
	log.Printf("Registering Client")
	credManager.RegisterClient(w, r)
}

// LoginClient checks provided credentials and stores session credentials and
// returns wether or not login was successful
func LoginClient(w http.ResponseWriter, r *http.Request) bool {
	if credManager.CheckClientCredentials(w, r) {
		sessManager.AuthenticateClient(w, r)
		log.Printf("Login successful")
		return true
	}
	log.Printf("Login failed")
	return false
}

// LogoutClient removes session stored auth cookie
func LogoutClient(w http.ResponseWriter, r *http.Request) {
	sessManager.DeauthenticateClient(w, r)
}

