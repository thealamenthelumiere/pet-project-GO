package service

import (
    "context"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"

    "github.com/thealamenthelumiere/pet-project-GO/internal/store"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound     = errors.New("user not found")
)
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type UserService interface {
    Login(ctx context.Context, username, password string) (*TokenPair, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
}

type userService struct {
	store  store.UserStore
	secret string
	accessTokenTTL time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(store store.UserStore, secret string, accessTTL, refreshTTL time.Duration) UserService {
    return &userService{
        store:          store,
        secret:         secret,
        accessTokenTTL: accessTTL,
        refreshTokenTTL: refreshTTL,
    }
}

func (s *userService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
    user, err := s.store.Get(ctx, username) 
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    if user.Password != password {
        return nil, ErrInvalidCredentials
    }

    accessToken, err := s.generateToken(username, s.accessTokenTTL)
    if err != nil {
        return nil, err
    }
    refreshToken, err := s.generateToken(username, s.refreshTokenTTL)
    if err != nil {
        return nil, err
    }

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
        return []byte(s.secret), nil
    })
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrTokenExpired
        }
        return nil, ErrInvalidToken
    }
    if !token.Valid {
        return nil, ErrInvalidToken
    }

    username := claims.Subject
    if username == "" {
        return nil, ErrInvalidToken
    }

    accessToken, err := s.generateToken(username, s.accessTokenTTL)
    if err != nil {
        return nil, err
    }
    newRefreshToken, err := s.generateToken(username, s.refreshTokenTTL)
    if err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
    }, nil
}
func (s *userService) generateToken(username string, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}