package user

import (
	"github.com/kllpw/kllpw-web/user/cred"
	sess "github.com/kllpw/kllpw-web/user/session"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"net/http"
)

// User is uuid of current user
type User struct {
	UUID *uuid.UUID
	Name string
}

var sessManager sess.Manager
var credManager cred.Manager

// Manager handles user with session manager and credentials managers
type Manager struct {
	sessManager      *sess.Manager
	credManager      *cred.Manager
	userSessionStore map[*uuid.UUID]*User
}

// NewManager returns a user manger with session and credential checks
func NewManager(sesskey string) *Manager {
	sm := sess.NewManager(sesskey)
	cm := cred.NewManager()
	css := make(map[*uuid.UUID]*User)
	m := Manager{sessManager: sm, credManager: cm, userSessionStore: css}
	return &m
}

// GetUser returns a populated user from session side checks returns nil if not found
func (m *Manager) GetUser(w http.ResponseWriter, r *http.Request) *User {
	if m.sessManager.IsUserAuthenticated(w, r) {
		uuid := m.sessManager.GetUserUUID(w, r)
		return m.userSessionStore[uuid]
	}
	return nil
}

// IsUserAuthenticated checks for stored session key if non found checks basic auth credentials
func (m *Manager) IsUserAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	return m.sessManager.IsUserAuthenticated(w, r)
}

// RegisterUser stores provided username and password for checks later
func (m *Manager) RegisterUser(w http.ResponseWriter, r *http.Request) bool {
	u, p, _ := r.BasicAuth()
	return m.credManager.RegisterUser(u, p)
}

// LoginUser checks provided credentials and stores session credentials and
// returns whether or not login was successful
func (m *Manager) LoginUser(w http.ResponseWriter, r *http.Request) bool {
	u, p, _ := r.BasicAuth()
	if m.credManager.IsUserCredentialsValid(u, p) {
		cUUID := m.sessManager.AuthenticateUser(w, r)
		user := User{UUID: cUUID, Name: u}
		m.userSessionStore[cUUID] = &user
		log.Printf("Login successful")
		return true
	}
	log.Printf("Login failed")
	return false
}

// LogoutUser removes session stored auth cookie
func (m *Manager) LogoutUser(w http.ResponseWriter, r *http.Request) {
	uuid := m.sessManager.GetUserUUID(w, r)
	delete(m.userSessionStore, uuid)
	m.sessManager.DeauthenticateUser(w, r)
}
