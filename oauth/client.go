package oauth

import (
	"errors"
	"strings"
	"time"

	"go-echo-api/entity"
	"go-echo-api/utils"
	"go-echo-api/utils/password"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	// ErrClientNotFound ...
	ErrClientNotFound = errors.New("Client not found")
	// ErrInvalidClientSecret ...
	ErrInvalidClientSecret = errors.New("Invalid client secret")
	// ErrClientIDTaken ...
	ErrClientIDTaken = errors.New("Client ID taken")
)

// ClientExists returns true if client exists
func (s *Service) ClientExists(clientID string) bool {
	_, err := s.FindClientByClientID(clientID)
	return err == nil
}

// FindClientByClientID looks up a client by client ID
func (s *Service) FindClientByClientID(clientID string) (*entity.OauthClient, error) {
	// Client IDs are case insensitive
	client := new(entity.OauthClient)
	notFound := s.db.Where("key = LOWER(?)", clientID).
		First(client).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrClientNotFound
	}

	return client, nil
}

// CreateClient saves a new client to database
func (s *Service) CreateClient(clientID, secret, redirectURI string) (*entity.OauthClient, error) {
	return s.createClientCommon(s.db, clientID, secret, redirectURI)
}

// CreateClientTx saves a new client to database using injected db object
func (s *Service) CreateClientTx(tx *gorm.DB, clientID, secret, redirectURI string) (*entity.OauthClient, error) {
	return s.createClientCommon(tx, clientID, secret, redirectURI)
}

// AuthClient authenticates client
func (s *Service) AuthClient(clientID, secret string) (*entity.OauthClient, error) {
	// Fetch the client
	client, err := s.FindClientByClientID(clientID)
	if err != nil {
		return nil, ErrClientNotFound
	}

	// Verify the secret
	if password.VerifyPassword(client.Secret, secret) != nil {
		return nil, ErrInvalidClientSecret
	}

	return client, nil
}

func (s *Service) createClientCommon(db *gorm.DB, clientID, secret, redirectURI string) (*entity.OauthClient, error) {
	// Check client ID
	if s.ClientExists(clientID) {
		return nil, ErrClientIDTaken
	}

	// Hash password
	secretHash, err := password.HashPassword(secret)
	if err != nil {
		return nil, err
	}

	client := &entity.OauthClient{
		MyGormModel: entity.MyGormModel{
			ID:        uuid.New().String(),
			CreatedAt: time.Now().UTC(),
		},
		Key:         strings.ToLower(clientID),
		Secret:      string(secretHash),
		RedirectURI: utils.StringOrNull(redirectURI),
	}
	if err := db.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}
