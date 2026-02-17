package store

import (
	"context"
	"errors"
)

type User struct {
	Username string
	Password string
}

type UserStore interface {
    Get(ctx context.Context, username string) (User, error)
    Create(ctx context.Context, user User) error
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

func (s *InMemoryStore) Get(ctx context.Context, username string) (User, error) {
select {
    case <-ctx.Done():
        return User{}, ctx.Err()
    default:
    }

    user, exists := s.users[username]
    if !exists {
        return User{}, errors.New("user not found")
    }
    return user, nil
}

func (s *InMemoryStore) Create(ctx context.Context, user User) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
	if _, exists := s.users[user.Username]; exists {
        return errors.New("user already exists")
    }
    s.users[user.Username] = user
    return nil
}