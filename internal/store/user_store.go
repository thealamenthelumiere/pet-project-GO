package store

import (
    "errors"
    "github.com/thealamenthelumiere/pet-project-GO/internal/service"
)


// service.User
type InMemoryStore struct {
    users map[string]service.User  //  используем service.User
}

func NewInMemoryStore() *InMemoryStore {
    return &InMemoryStore{
        users: map[string]service.User{
            "admin": {Username: "admin", Password: "password123"},
            "user1": {Username: "user1", Password: "secret"},
        },
    }
}

func (s *InMemoryStore) Get(username string) (service.User, error) {  //  возвращаем service.User
    user, exists := s.users[username]
    if !exists {
        return service.User{}, errors.New("user not found")
    }
    return user, nil
}
