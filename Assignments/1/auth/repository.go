package auth

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gitlab.com/garyteofanus/law-assignment/domain"
	"time"
)

type repository struct {
	db *redis.Client
}

func NewRepo(db *redis.Client) Repository {
	return &repository{
		db: db,
	}
}

type Repository interface {
	GetUser() (*domain.User, error)
	GetSession() (*domain.Session, error)
	SetSession(session domain.Session) error
}

func (r repository) GetUser() (*domain.User, error) {
	username, err := r.db.Get(context.TODO(), "username").Result()
	if err != nil {
		return nil, err
	}

	password, err := r.db.Get(context.TODO(), "password").Result()
	if err != nil {
		return nil, err
	}

	clientID, err := r.db.Get(context.TODO(), "client_id").Result()
	if err != nil {
		return nil, err
	}

	clientSecret, err := r.db.Get(context.TODO(), "client_secret").Result()
	if err != nil {
		return nil, err
	}

	fullName, err := r.db.Get(context.TODO(), "full_name").Result()
	if err != nil {
		return nil, err
	}

	npm, err := r.db.Get(context.TODO(), "npm").Result()
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       clientID,
		Secret:   clientSecret,
		Username: username,
		Password: password,
		FullName: fullName,
		NPM:      npm,
	}, nil
}

func (r repository) GetSession() (*domain.Session, error) {
	accessToken, err := r.db.Get(context.TODO(), "access_token").Result()
	if err != nil {
		return nil, err
	}

	refreshToken, err := r.db.Get(context.TODO(), "refresh_token").Result()
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r repository) SetSession(session domain.Session) error {
	// Access token duration is 5 minutes
	err := r.db.Set(context.TODO(), "access_token", session.AccessToken, time.Minute*5).Err()
	if err != nil {
		return err
	}

	// Refresh token duration is 200 days
	err = r.db.Set(context.TODO(), "refresh_token", session.RefreshToken, time.Hour*24*200).Err()
	if err != nil {
		return err
	}

	return nil
}
