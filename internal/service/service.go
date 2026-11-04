package service

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

// TODO: move to config
const aliasLength = 6

const maxAliasAttempts = 5

// Storage is the persistence port used by the service layer.
type Storage interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	UpdateURL(alias string, newURL string) (int64, error)
	DeleteURL(alias string) (int64, error)
}

// Service contains business logic for URL shortening.
type Service struct {
	storage  Storage
	validate *validator.Validate
}

func New(storage Storage) *Service {
	return &Service{
		storage:  storage,
		validate: validator.New(),
	}
}

type saveInput struct {
	URL string `validate:"required,url"`
}

type updateInput struct {
	Alias  string `validate:"required"`
	NewURL string `validate:"required,url"`
}

// Save validates the URL, generates an alias when needed, and persists the mapping.
func (s *Service) Save(rawURL, alias string) (string, error) {
	if err := s.validate.Struct(saveInput{URL: rawURL}); err != nil {
		return "", err
	}

	userProvidedAlias := alias != ""
	if !userProvidedAlias {
		alias = random.NewRandomString(aliasLength)
	}

	for attempt := 0; attempt < maxAliasAttempts; attempt++ {
		_, err := s.storage.SaveURL(rawURL, alias)
		if err == nil {
			return alias, nil
		}

		if errors.Is(err, storage.ErrURLExists) {
			if userProvidedAlias {
				return "", storage.ErrURLExists
			}
			alias = random.NewRandomString(aliasLength)
			continue
		}

		return "", err
	}

	return "", fmt.Errorf("failed to generate unique alias")
}

// Get returns the original URL for alias.
func (s *Service) Get(alias string) (string, error) {
	if alias == "" {
		return "", storage.ErrURLNotFound
	}

	return s.storage.GetURL(alias)
}

// Update changes the destination URL for an existing alias.
func (s *Service) Update(alias, newURL string) (int64, error) {
	if err := s.validate.Struct(updateInput{Alias: alias, NewURL: newURL}); err != nil {
		return 0, err
	}

	n, err := s.storage.UpdateURL(alias, newURL)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, storage.ErrURLNotFound
	}

	return n, nil
}

// Delete removes a URL mapping by alias.
func (s *Service) Delete(alias string) (int64, error) {
	if alias == "" {
		return 0, storage.ErrURLNotFound
	}

	n, err := s.storage.DeleteURL(alias)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, storage.ErrURLNotFound
	}

	return n, nil
}
