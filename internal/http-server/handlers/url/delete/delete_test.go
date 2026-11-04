package delete_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
)

type deleterStub struct{ delete func(string) (int64, error) }

func (s deleterStub) Delete(alias string) (int64, error) { return s.delete(alias) }

func TestDelete_Success(t *testing.T) {
	handler := delete.New(slogdiscard.NewDiscardLogger(), deleterStub{delete: func(alias string) (int64, error) {
		require.Equal(t, "example", alias)
		return 1, nil
	}})
	rr := serveDelete(handler, "example")
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDelete_NotFound(t *testing.T) {
	handler := delete.New(slogdiscard.NewDiscardLogger(), deleterStub{delete: func(string) (int64, error) {
		return 0, storage.ErrURLNotFound
	}})
	rr := serveDelete(handler, "missing")
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDelete_BadRequest(t *testing.T) {
	handler := delete.New(slogdiscard.NewDiscardLogger(), deleterStub{delete: func(string) (int64, error) {
		t.Fatal("service must not be called without an alias")
		return 0, nil
	}})
	rr := serveDelete(handler, "")
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func serveDelete(handler http.Handler, alias string) *httptest.ResponseRecorder {
	r := chi.NewRouter()
	r.Method(http.MethodDelete, "/{alias}", handler)
	req := httptest.NewRequest(http.MethodDelete, "/"+alias, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}
