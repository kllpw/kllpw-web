package cred

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Manager is for storing user Creds
type Manager struct {
	users map[string]credentials
}

// NewManager returns a new manager
func NewManager() *Manager {
	m := Manager{users: make(map[string]credentials, 0)}
	return &m
}

// RegisterUser adds user to  store
func (m *Manager) RegisterUser(username string, password string) bool {
	_, found := m.users[username]
	if !found {
		if m.users == nil {
			m.users = make(map[string]credentials, 0)
		}
		pwd, err := m.getSaltedHashedPassword(password)
		if err != nil {
			return false
		}
		m.users[username] = credentials{username: username, password: pwd}
		log.Printf("Registering: %s", username)
		return true
	}
	log.Printf("User already exists: %s", username)
	return false
}

// IsUserCredentialsValid checks if the user is registered and details match those stored server side
func (m *Manager) IsUserCredentialsValid(username string, password string) bool {
	log.Printf("Validating login for: %s", username)
	usrCrd, found := m.users[username]
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

func (m *Manager) getSaltedHashedPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return hash, nil

}
