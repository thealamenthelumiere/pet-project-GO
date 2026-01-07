package service

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// User представляет пользователя
type User struct {
    Username string
    Password string
}

// Claims для JWT
type Claims struct {
    Username string `json:"sub"`
    jwt.RegisteredClaims
}

// TokenError тип ошибки для токенов
type TokenError struct {
    message string
}

func NewTokenError(message string) *TokenError {
    return &TokenError{message: message}
}

func (e *TokenError) Error() string {
    return e.message
}

// UserStore интерфейс для хранилища
type UserStore interface {
    Get(username string) (User, error)
}

// UserService интерфейс бизнес-логики
type UserService interface {
    ValidateCredentials(username, password string) (bool, error)
    GenerateToken(username string) (string, error)
    RefreshToken(token string) (string, error)
}

// userService реализация (неэкспортируемая)
type userService struct {
    store  UserStore
    secret string
}

// NewUserService конструктор (ЭКСПОРТИРУЕМАЯ функция)
func NewUserService(store UserStore, secret string) UserService {
    return &userService{
        store:  store,
        secret: secret,
    }
}

// GenerateToken создает JWT токен
func (s *userService) GenerateToken(username string) (string, error) {
    claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   username,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secret))
}

// ValidateCredentials проверяет учетные данные
func (s *userService) ValidateCredentials(username, password string) (bool, error) {
    user, err := s.store.Get(username)
    if err != nil {
        return false, err
    }

    return user.Password == password, nil
}

// RefreshToken обновляет токен
func (s *userService) RefreshToken(tokenString string) (string, error) {
    // Парсим и валидируем токен
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.secret), nil
    })

    if err != nil {
        return "", NewTokenError("invalid token")
    }

    if !token.Valid {
        return "", NewTokenError("invalid token")
    }

    // Проверяем что токен не истек (имеет смысл обновлять)
    if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
        return "", NewTokenError("token expired")
    }

    // Генерируем новый токен
    return s.GenerateToken(claims.Username)
}