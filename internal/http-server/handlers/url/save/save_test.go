package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
)

func TestSave_Success(t *testing.T) {
	saver := mocks.NewURLSaver(t)
	saver.On("Save", "https://example.com", "example").Return("example", nil).Once()

	rr := serve(t, save.New(slogdiscard.NewDiscardLogger(), saver), `{"url":"https://example.com","alias":"example"}`)
	require.Equal(t, http.StatusCreated, rr.Code)

	var response save.Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
	require.Equal(t, "example", response.Alias)
}

func TestSave_BadRequest(t *testing.T) {
	saver := mocks.NewURLSaver(t)

	// The JSON is deliberately malformed: the handler must reject it before calling service.
	rr := serve(t, save.New(slogdiscard.NewDiscardLogger(), saver), `{"url":`)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var response save.Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
	require.Equal(t, "failed to decode request", response.Error)
}

func TestSave_InvalidURL(t *testing.T) {
	saver := mocks.NewURLSaver(t)
	validationErr := validator.New().Struct(struct {
		URL string `validate:"required,url"`
	}{URL: "not-a-url"})
	saver.On("Save", "not-a-url", "example").Return("", validationErr).Once()

	rr := serve(t, save.New(slogdiscard.NewDiscardLogger(), saver), `{"url":"not-a-url","alias":"example"}`)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSave_Conflict(t *testing.T) {
	saver := mocks.NewURLSaver(t)
	saver.On("Save", "https://example.com", "taken").Return("", storage.ErrURLExists).Once()

	rr := serve(t, save.New(slogdiscard.NewDiscardLogger(), saver), `{"url":"https://example.com","alias":"taken"}`)
	require.Equal(t, http.StatusConflict, rr.Code)

	var response save.Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
	require.Equal(t, "url already exists", response.Error)
}

func TestSave_InternalError(t *testing.T) {
	saver := mocks.NewURLSaver(t)
	saver.On("Save", "https://example.com", "example").Return("", errors.New("storage unavailable")).Once()

	rr := serve(t, save.New(slogdiscard.NewDiscardLogger(), saver), `{"url":"https://example.com","alias":"example"}`)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func serve(t *testing.T, handler http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
