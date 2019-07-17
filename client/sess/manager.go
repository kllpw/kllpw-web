package client

import (
	"log"
	"net/http"
	"github.com/gorilla/sessions"
)

const (
	authKey string = "auth"
	sessionKey string = "kllpw"
)

// Manager is for storing client sessions
type Manager struct {
	store *sessions.CookieStore
}

// NewManager returns a new manager with key from OS variable "SESSION_KEYS"
func NewManager() *Manager {
	m := Manager{}
	m.store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	return &m
}

// AuthenticateClient adds client to session store
func (m *Manager) AuthenticateClient(w http.ResponseWriter, r *http.Request) {
	session, _ := m.store.Get(r, sessionKey)
	session.Values[authKey] = true
	session.Save(r, w)
	
	log.Printf("Session Authenticated")
}

// DeauthenticateClient removes client from session store
func (m *Manager) DeauthenticateClient(w http.ResponseWriter, r *http.Request){
	session, _ := m.store.Get(r, sessionKey)
	session.Values[authKey] = false
	session.Save(r, w)
	
	log.Print("Session Deauthenticated")
}

// IsClientAuthed checks if request is in the session store and has the authKey value set to true
func (m *Manager) IsClientAuthed(w http.ResponseWriter, r *http.Request) bool {
	session, _ := m.store.Get(r, sessionKey)
	auth := session.Values[authKey]
	
	if auth != nil {
		log.Printf("Session Status: %t", auth)
		return auth.(bool)
	}
	return false
}