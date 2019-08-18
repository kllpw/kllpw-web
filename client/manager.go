package client

import (
	"log"
	"net/http"
	uuid "github.com/nu7hatch/gouuid"
	sess "github.com/kllpw/kllpw-web/client/session"
	cred "github.com/kllpw/kllpw-web/client/cred"
)
// Client is uuid of current client
type Client struct {
	UUID *uuid.UUID
	Name string
}

var sessManager sess.Manager
var credManager cred.Manager

// Manager handles client with session manager and credetials managers
type Manager struct {
	sessManager *sess.Manager
	credManager *cred.Manager
	clientSessionStore map[*uuid.UUID]*Client
}

// NewManager returns a client manger with session and credential checks
func NewManager(sesskey string) *Manager {
	sm := sess.NewManager(sesskey)
	cm := cred.NewManager()
	css := make(map[*uuid.UUID]*Client)
	m := Manager{sessManager: sm, credManager: cm, clientSessionStore: css}
	return &m
}
// GetClient returns a populated client from session side checks returns nil if not found
func (m *Manager) GetClient(w http.ResponseWriter, r *http.Request) *Client {
	if m.sessManager.IsClientAuthed(w, r) {
		uuid, err := m.sessManager.GetClientUUID(w, r)
		if err != nil {
			log.Printf("Client not found %s", uuid)
			return nil
		}
		return m.clientSessionStore[uuid]
	}
	return nil
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
		uuid, _ := m.sessManager.GetClientUUID(w, r)
		if uuid != nil {
			log.Printf("Login successful reissued token")
			return true
		}
		cUUID := m.sessManager.AuthenticateClient(w, r)
		client := Client{UUID: cUUID, Name: u}
		m.clientSessionStore[cUUID] = &client
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

