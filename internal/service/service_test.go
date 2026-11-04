package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

type stubStorage struct {
	saveFn   func(url, alias string) (int64, error)
	getFn    func(alias string) (string, error)
	updateFn func(alias, newURL string) (int64, error)
	deleteFn func(alias string) (int64, error)
}

func (s *stubStorage) SaveURL(url, alias string) (int64, error) {
	return s.saveFn(url, alias)
}

func (s *stubStorage) GetURL(alias string) (string, error) {
	return s.getFn(alias)
}

func (s *stubStorage) UpdateURL(alias, newURL string) (int64, error) {
	return s.updateFn(alias, newURL)
}

func (s *stubStorage) DeleteURL(alias string) (int64, error) {
	return s.deleteFn(alias)
}

func TestService_Save_ValidatesURL(t *testing.T) {
	svc := service.New(&stubStorage{})

	_, err := svc.Save("", "alias")
	require.Error(t, err)

	_, err = svc.Save("not-a-url", "alias")
	require.Error(t, err)
}

func TestService_Save_GeneratesAlias(t *testing.T) {
	var savedAlias string
	svc := service.New(&stubStorage{
		saveFn: func(url, alias string) (int64, error) {
			savedAlias = alias
			assert.Equal(t, "https://example.com", url)
			assert.Len(t, alias, 6)
			return 1, nil
		},
	})

	alias, err := svc.Save("https://example.com", "")
	require.NoError(t, err)
	assert.Equal(t, savedAlias, alias)
}

func TestService_Save_UserAliasExists(t *testing.T) {
	svc := service.New(&stubStorage{
		saveFn: func(url, alias string) (int64, error) {
			return 0, storage.ErrURLExists
		},
	})

	_, err := svc.Save("https://example.com", "taken")
	require.ErrorIs(t, err, storage.ErrURLExists)
}

func TestService_Get_PassesThrough(t *testing.T) {
	svc := service.New(&stubStorage{
		getFn: func(alias string) (string, error) {
			assert.Equal(t, "demo", alias)
			return "https://example.com", nil
		},
	})

	got, err := svc.Get("demo")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", got)
}

func TestService_Update_NotFound(t *testing.T) {
	svc := service.New(&stubStorage{
		updateFn: func(alias, newURL string) (int64, error) {
			return 0, nil
		},
	})

	_, err := svc.Update("missing", "https://example.com")
	require.ErrorIs(t, err, storage.ErrURLNotFound)
}

func TestService_Delete_NotFound(t *testing.T) {
	svc := service.New(&stubStorage{
		deleteFn: func(alias string) (int64, error) {
			return 0, nil
		},
	})

	_, err := svc.Delete("missing")
	require.ErrorIs(t, err, storage.ErrURLNotFound)
}

func TestService_Save_RetriesGeneratedAlias(t *testing.T) {
	attempts := 0
	svc := service.New(&stubStorage{
		saveFn: func(url, alias string) (int64, error) {
			attempts++
			if attempts == 1 {
				return 0, storage.ErrURLExists
			}
			return 1, nil
		},
	})

	alias, err := svc.Save("https://example.com", "")
	require.NoError(t, err)
	assert.NotEmpty(t, alias)
	assert.Equal(t, 2, attempts)
}

func TestService_Get_EmptyAlias(t *testing.T) {
	svc := service.New(&stubStorage{
		getFn: func(alias string) (string, error) {
			return "", errors.New("should not be called")
		},
	})

	_, err := svc.Get("")
	require.ErrorIs(t, err, storage.ErrURLNotFound)
}
