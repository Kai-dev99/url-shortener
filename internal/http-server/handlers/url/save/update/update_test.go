package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/save/update"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
)

type updaterStub struct {
	update func(string, string) (int64, error)
}

func (s updaterStub) Update(alias, destination string) (int64, error) {
	return s.update(alias, destination)
}

func TestUpdate_Success(t *testing.T) {
	handler := update.New(slogdiscard.NewDiscardLogger(), updaterStub{update: func(alias, destination string) (int64, error) {
		require.Equal(t, "example", alias)
		require.Equal(t, "https://example.org", destination)
		return 1, nil
	}})
	rr := serveUpdate(handler, `{"alias":"example","url":"https://example.org"}`)
	require.Equal(t, http.StatusOK, rr.Code)

	var response update.Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
	require.EqualValues(t, 1, response.CountUpdated)
}

func TestUpdate_BadRequest(t *testing.T) {
	handler := update.New(slogdiscard.NewDiscardLogger(), updaterStub{update: func(string, string) (int64, error) {
		t.Fatal("service must not be called for malformed JSON")
		return 0, nil
	}})
	rr := serveUpdate(handler, `{"alias":`)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdate_NotFound(t *testing.T) {
	handler := update.New(slogdiscard.NewDiscardLogger(), updaterStub{update: func(string, string) (int64, error) {
		return 0, storage.ErrURLNotFound
	}})
	rr := serveUpdate(handler, `{"alias":"missing","url":"https://example.org"}`)
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdate_InternalError(t *testing.T) {
	handler := update.New(slogdiscard.NewDiscardLogger(), updaterStub{update: func(string, string) (int64, error) {
		return 0, errors.New("storage unavailable")
	}})
	rr := serveUpdate(handler, `{"alias":"example","url":"https://example.org"}`)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func serveUpdate(handler http.Handler, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/url", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
