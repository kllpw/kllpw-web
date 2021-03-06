package client

import (
	"encoding/gob"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	authKey string = "auth"
)

// Manager is for storing http user sessions
type Manager struct {
	store          *sessions.CookieStore
	activeSessions []*uuid.UUID
	sessionKey     string
}

// NewManager returns a new manager with key provided
func NewManager(sessionkey string) *Manager {
	m := Manager{}
	m.sessionKey = sessionkey
	m.store = sessions.NewCookieStore([]byte(m.sessionKey))
	m.activeSessions = make([]*uuid.UUID, 0)
	return &m
}

func init() {
	gob.Register(uuid.UUID{})
}

// AuthenticateUser adds user to session store
func (m *Manager) AuthenticateUser(w http.ResponseWriter, r *http.Request) *uuid.UUID {
	session, _ := m.store.Get(r, m.sessionKey)
	cUUID, _ := uuid.NewV4()
	m.activeSessions = append(m.activeSessions, cUUID)
	session.Values[authKey] = cUUID
	session.Save(r, w)
	log.Print("Session Authenticated ", cUUID)
	return cUUID
}

// DeauthenticateUser removes user from session store
func (m *Manager) DeauthenticateUser(w http.ResponseWriter, r *http.Request) {
	session, _ := m.store.Get(r, m.sessionKey)
	_, uuidPos, _ := m.getUserUUIDAndPosition(w, r)
	if uuidPos > -1 {
		m.removeUUID(uuidPos)
		session.Save(r, w)
	}
}

func (m *Manager) getUserUUIDAndPosition(w http.ResponseWriter, r *http.Request) (*uuid.UUID, int, error) {
	session, _ := m.store.Get(r, m.sessionKey)
	currentUUID := session.Values[authKey]
	if currentUUID != nil {
		uuidPos := m.getUUIDPosition(currentUUID)
		if uuidPos > -1 {
			return m.activeSessions[uuidPos], uuidPos, nil
		}
	}
	return nil, -1, errors.New("no UUID found")
}

func (m *Manager) getUUIDPosition(currentUUID interface{}) int {
	uuidPos := -1
	switch u := currentUUID.(type) {
	case uuid.UUID:
		uuidPos = m.findUUID(&u)
	case *uuid.UUID:
		uuidPos = m.findUUID(u)
	}
	return uuidPos
}

// GetUserUUID gets uuid for current request session
func (m *Manager) GetUserUUID(w http.ResponseWriter, r *http.Request) *uuid.UUID {
	cUUID, _, _ := m.getUserUUIDAndPosition(w, r)
	return cUUID
}

func (m *Manager) findUUID(currentUUID *uuid.UUID) (pos int) {
	for pos, v := range m.activeSessions {
		if currentUUID.String() == v.String() {
			return pos
		}
	}
	return -1
}

func (m *Manager) removeUUID(index int) {
	m.activeSessions = append(m.activeSessions[:index], m.activeSessions[index+1:]...)
}

// IsUserAuthenticated checks if request is in the session store and has the authKey value set to true
func (m *Manager) IsUserAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	_, _, err := m.getUserUUIDAndPosition(w, r)
	if err != nil {
		return false
	}
	return true
}
