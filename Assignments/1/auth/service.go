package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"gitlab.com/garyteofanus/law-assignment/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

type Service interface {
	AuthorizeUser(user domain.User) (string, string, error)
	AccessResource(token string) (*domain.User, *domain.Session, error)
}

func (s service) AuthorizeUser(user domain.User) (string, string, error) {
	u, err := s.repo.GetUser()
	if err != nil {
		log.Println("[failed to get user]", err)
		return "", "", errors.New("failed to fetch user from database")
	}

	if user.Username != u.Username {
		return "", "", errors.New("invalid username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		log.Println("[failed to compare password]", err)
		return "", "", errors.New("invalid password")
	}

	if user.ID != u.ID {
		log.Println("[failed to compare id]", err)
		return "", "", errors.New("invalid client id")
	}

	if user.Secret != u.Secret {
		log.Println("[failed to compare secret]", err)
		return "", "", errors.New("invalid client secret")
	}

	accessToken := generateSecureToken(40)
	refreshToken := generateSecureToken(40)

	if err := s.repo.SetSession(domain.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s service) AccessResource(token string) (*domain.User, *domain.Session, error) {
	session, err := s.repo.GetSession()
	if err != nil {
		return nil, nil, err
	}

	if session.AccessToken != token {
		return nil, nil, errors.New("invalid token")
	}

	user, err := s.repo.GetUser()
	if err != nil {
		log.Println("[failed to get user]", err)
		return nil, nil, errors.New("failed to fetch user from database")
	}

	return user, session, nil
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
