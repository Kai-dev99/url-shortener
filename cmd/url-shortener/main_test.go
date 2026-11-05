package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

type serviceStub struct{}

func (serviceStub) Save(string, string) (string, error)  { return "", nil }
func (serviceStub) Get(string) (string, error)           { return "", nil }
func (serviceStub) Update(string, string) (int64, error) { return 0, nil }
func (serviceStub) Delete(string) (int64, error)         { return 0, nil }

func TestProtectedRoutes_MissingBasicAuth(t *testing.T) {
	router := newRouter(slogdiscard.NewDiscardLogger(), serviceStub{}, "user", "password")
	req := httptest.NewRequest(http.MethodPost, "/url/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestProtectedRoutes_InvalidBasicAuth(t *testing.T) {
	router := newRouter(slogdiscard.NewDiscardLogger(), serviceStub{}, "user", "password")
	req := httptest.NewRequest(http.MethodPost, "/url/", nil)
	req.SetBasicAuth("user", "wrong-password")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
