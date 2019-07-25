package client

import (
	"log"
	"net/http"
	uuid "github.com/nu7hatch/gouuid"
	sess "github.com/kllpw/kllpw-web/client/sess"
	cred "github.com/kllpw/kllpw-web/client/cred"
)
// Client is uuid of current client
type Client struct {
	uuid *uuid.UUID
}

var sessManager sess.Manager
var credManager cred.Manager

// Manager handles client with session manager and credetials managers
type Manager struct {
	sessManager *sess.Manager
	credManager *cred.Manager
}

// NewManager returns a client manger with session and credential checks
func NewManager(sesskey string) *Manager {
	sm := sess.NewManager(sesskey)
	cm := cred.NewManager()
	m := Manager{sessManager: sm, credManager: cm}
	return &m
}

// IsValidClient checks for stored session key if non found checks basic auth credentials
func (m *Manager) IsValidClient(w http.ResponseWriter, r *http.Request) bool {
	return m.sessManager.IsClientAuthed(w, r)
}

// RegisterClient stores provided username and password for checks later
func (m *Manager) RegisterClient(w http.ResponseWriter, r *http.Request) bool {
	u, p, _ := r.BasicAuth()
	return m.credManager.RegisterClient(u,p)
}

// LoginClient checks provided credentials and stores session credentials and
// returns wether or not login was successful
func (m *Manager) LoginClient(w http.ResponseWriter, r *http.Request) bool {
	u, p, _ := r.BasicAuth()
	if m.credManager.CheckClientCredentials(u, p) {
		m.sessManager.AuthenticateClient(w, r)
		log.Printf("Login successful")
		return true
	}
	log.Printf("Login failed")
	return false
}

// LogoutClient removes session stored auth cookie
func (m *Manager) LogoutClient(w http.ResponseWriter, r *http.Request) {
	m.sessManager.DeauthenticateClient(w, r)
}

