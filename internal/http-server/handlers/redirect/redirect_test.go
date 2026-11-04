package redirect_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("Get", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}

func TestRedirect_NotFound(t *testing.T) {
	urlGetterMock := mocks.NewURLGetter(t)
	urlGetterMock.On("Get", "missing").Return("", storage.ErrURLNotFound).Once()

	r := chi.NewRouter()
	r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRedirect_InternalError(t *testing.T) {
	urlGetterMock := mocks.NewURLGetter(t)
	urlGetterMock.On("Get", "broken").Return("", errors.New("database unavailable")).Once()

	r := chi.NewRouter()
	r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

	req := httptest.NewRequest(http.MethodGet, "/broken", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
