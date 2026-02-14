package store

import "errors"

// --------------------------------------------------------------------
// Модель пользователя.
// --------------------------------------------------------------------
type User struct {
	Username string
	Password string
}

// --------------------------------------------------------------------
// Интерфейс хранилища — контракт для сервиса.
// --------------------------------------------------------------------
type UserStore interface {
	// Get возвращает пользователя по имени или ошибку.
	Get(username string) (User, error)
	// Create добавляет нового пользователя.
	Create(user User) error
}

// --------------------------------------------------------------------
// In-memory реализация.
// --------------------------------------------------------------------

// InMemoryStore — хранилище в оперативной памяти.
type InMemoryStore struct {
	users map[string]User
}

// NewInMemoryStore — конструктор, инициализирует пустую мапу или с тестовыми данными.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: map[string]User{
			"admin": {Username: "admin", Password: "password123"},
			"user1": {Username: "user1", Password: "secret"},
		},
	}
}

// Get возвращает пользователя по имени.
func (s *InMemoryStore) Get(username string) (User, error) {
	user, exists := s.users[username]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

// Create добавляет нового пользователя.
func (s *InMemoryStore) Create(user User) error {
	if _, exists := s.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	s.users[user.Username] = user
	return nil
}