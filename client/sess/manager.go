package client

import (
	"encoding/gob"
	"log"
	"net/http"
	"github.com/gorilla/sessions"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	authKey string = "auth"
	sessionKey string = "kllpw"
)

// Manager is for storing client sessions
type Manager struct {
	store *sessions.CookieStore
	activeSessions []*uuid.UUID
}


// NewManager returns a new manager with key from OS variable "SESSION_KEYS"
func NewManager(sesskey string) *Manager {
	m := Manager{}
	m.store = sessions.NewCookieStore([]byte(sessionKey))
	return &m
}

func init(){
	gob.Register(uuid.UUID{})
}

// AuthenticateClient adds client to session store
func (m *Manager) AuthenticateClient(w http.ResponseWriter, r *http.Request) {
	session, _ := m.store.Get(r, sessionKey)
	cUUID, _ := uuid.NewV4()
	m.activeSessions = append(m.activeSessions, cUUID)
	session.Values[authKey] = cUUID
	session.Save(r, w)
	log.Print("Session Authenticated ", cUUID)
}

// DeauthenticateClient removes client from session store
func (m *Manager) DeauthenticateClient(w http.ResponseWriter, r *http.Request){
	session, _ := m.store.Get(r, sessionKey)
	currentUUID := session.Values[authKey]
	if currentUUID != nil {
		u := currentUUID.(*uuid.UUID)
		uuidPos := m.findUUID(u)
		if uuidPos > -1 {
			m.removeUUID(uuidPos)
			log.Print("Session Deauthenticated ", currentUUID)
		}
	}
	session.Save(r, w)
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


// IsClientAuthed checks if request is in the session store and has the authKey value set to true
func (m *Manager) IsClientAuthed(w http.ResponseWriter, r *http.Request) bool {
	session, _ := m.store.Get(r, sessionKey)
	currentUUID := session.Values[authKey]
	if currentUUID != nil {
		u := currentUUID.(*uuid.UUID)
		uuidPos := m.findUUID(u)
		if uuidPos > -1 {
			log.Print(currentUUID, " is an Active session")
			return true
		}
	}
	log.Print(currentUUID, " is an not an Active session")
	return false
}