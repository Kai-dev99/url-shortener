package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		url          string
		respError    string
		mockError    error
		mockAliasOut string
		skipMock     bool
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			url:          "https://google.com",
			mockAliasOut: "test_alias",
		},
		{
			name:         "Empty alias",
			alias:        "",
			url:          "https://google.com",
			mockAliasOut: "abcdef",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
			mockError: validator.New().Struct(struct {
				URL string `validate:"required,url"`
			}{}),
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
			mockError: validator.New().Struct(struct {
				URL string `validate:"required,url"`
			}{URL: "some invalid URL"}),
		},
		{
			name:      "URL already exists",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "url already exists",
			mockError: storage.ErrURLExists,
		},
		{
			name:      "Save Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlSaverMock := mocks.NewURLSaver(t)

			if !tc.skipMock {
				outAlias := tc.mockAliasOut
				if outAlias == "" {
					outAlias = tc.alias
				}
				urlSaverMock.On("Save", tc.url, tc.alias).
					Return(outAlias, tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var resp save.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
