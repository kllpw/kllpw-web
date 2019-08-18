package cred

import (
	"golang.org/x/crypto/bcrypt"
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
		pwd, err := m.getSaltedHashedPassword(password)
		if err != nil {
			return false
		}
		m.clients[username] = credentials{username: username, password: pwd}
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
		err := bcrypt.CompareHashAndPassword(usrCrd.password, []byte(password))
		if err != nil {
			log.Println(err)
			return false
		}
		log.Printf("Login credentials valid.")
		return true
	}
	log.Printf("Login credentials invalid.")
	return false
}

func (m *Manager) getSaltedHashedPassword(password string) ( []byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return hash, nil

}
