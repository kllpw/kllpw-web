package cred

import (
	"log"
)

// Manager is for storing client Creds
type Manager struct {
	clients map[string]credentials
}

// NewManager returns a new manager
func NewManager() *Manager {
	m := Manager{clients: make(map[string]credentials, 0)}
	return &m
}

// RegisterClient adds client to  store
func (m *Manager) RegisterClient(username string, password string) bool {
	_, found := m.clients[username]
	if !found {
		if m.clients == nil {
			m.clients = make(map[string]credentials, 0)
		}
		m.clients[username] = credentials{username: username, password:[]byte(password)}
		log.Printf("Registering: %s", username)
		return true
	}
	log.Printf("Client already exists: %s", username)
	return false
}

// CheckClientCredentials checks if the client is registered and details match those stored server side
func (m *Manager) CheckClientCredentials(username string, password string) bool {
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
