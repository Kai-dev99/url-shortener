package update

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type Request struct {
	Alias  string `json:"alias"`
	NewURL string `json:"url"`
}

type Response struct {
	resp.Response
	CountUpdated int64 `json:"countUpdated"`
}

type URLUpdater interface {
	Update(alias string, newURL string) (int64, error)
}

func New(log *slog.Logger, updater URLUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		countUpdated, err := updater.Update(req.Alias, req.NewURL)
		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", req.Alias)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to update url", "alias", req.Alias, sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update url"))
			return
		}

		log.Info("url updated", slog.String("alias", req.Alias), slog.Int64("count_updated", countUpdated))
		render.JSON(w, r, Response{
			Response:     resp.OK(),
			CountUpdated: countUpdated,
		})
	}
}
