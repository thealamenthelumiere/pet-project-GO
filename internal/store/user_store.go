package store

import "errors"

type User struct {
	Username string
	Password string
}

type UserStore interface {
	
	Get(username string) (User, error)
	
	Create(user User) error
}

type InMemoryStore struct {
	users map[string]User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: map[string]User{
			"admin": {Username: "admin", Password: "password123"},
			"user1": {Username: "user1", Password: "secret"},
		},
	}
}

func (s *InMemoryStore) Get(username string) (User, error) {
	user, exists := s.users[username]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *InMemoryStore) Create(user User) error {
	if _, exists := s.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	s.users[user.Username] = user
	return nil
}