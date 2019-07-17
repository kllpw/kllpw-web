package cred

import (
	"net/http"
	"log"
)

// Manager is for storing client Creds
type Manager struct {
	clients map[string]credentials
}

// NewManager returns a new manager
func NewManager() *Manager {
	m := Manager{}
	m.clients = make(map[string]credentials, 0)
	return &m
}

// RegisterClient adds client to  store
func (m *Manager) RegisterClient(w http.ResponseWriter, r *http.Request){
	usr, pwd, _ := r.BasicAuth()
	log.Printf("Registering: %s", usr)
	m.clients[usr] = credentials{username: usr, password: []byte(pwd)}
}

// CheckClientCredentials checks if the client is registered and details match those stored server side
func (m *Manager) CheckClientCredentials(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	username, password, _ := r.BasicAuth()
	log.Printf("Validating login for: %s", username)
	usrCrd, found := m.clients[username]
	if found {
		if password == string(usrCrd.password) {
			log.Printf("Login credentials valid.")
			return true
		} 
	}
	log.Printf("Login credentials invalid.")
	return false
} 