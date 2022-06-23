package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"net/http"
	"projects/practic_5course_cesar/internal/cryptosvc"
)

type CryptoService interface {
	Encrypt(ctx context.Context, lang cryptosvc.Language, phrase cryptosvc.Data) (cryptosvc.Result, error)
	Decrypt(ctx context.Context, lang cryptosvc.Language, phrase cryptosvc.Data) (cryptosvc.Result, error)
}

type Handler struct {
	crypto  CryptoService
	logger  *zap.Logger
	decoder *schema.Decoder
}

func New(svc CryptoService, dc *schema.Decoder, log *zap.Logger) *Handler {
	return &Handler{
		crypto:  svc,
		decoder: dc,
		logger:  log,
	}
}

func (h *Handler) OK(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("check service status", zap.String("user_agent", r.UserAgent()))
	h.write(w, http.StatusOK, "crypto-service")
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.logger.Warn("invalid path", zap.String("user_agent", r.UserAgent()))
	h.write(w, http.StatusNotFound, "invalid path")
}

func (h *Handler) write(w http.ResponseWriter, statusCode int, v interface{}) { //nolint:varnamelen
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	h.logger.Info("output", zap.Any("code", statusCode), zap.Any("response", v))

	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Error("json encoder, write response", zap.Error(err), zap.Any("response", v))
	}
}

func (h *Handler) WriteAnswer(w http.ResponseWriter, code int, answer interface{}) {
	h.write(w, code, answer)
}

func (h *Handler) Encrypt(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	lang := params.Get("lang")
	data := params.Get("data")

	response, err := h.crypto.Encrypt(r.Context(), cryptosvc.Language(lang), cryptosvc.Data(data))
	if err != nil {
		h.logger.Error("encrypt fail", zap.Error(err), zap.String("language", lang), zap.String("data", data))
		h.WriteAnswer(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteAnswer(w, http.StatusOK, response)
}

func (h *Handler) Decrypt(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	lang := params.Get("lang")
	data := params.Get("data")

	response, err := h.crypto.Decrypt(r.Context(), cryptosvc.Language(lang), cryptosvc.Data(data))
	if err != nil {
		h.logger.Error("decrypt fail", zap.Error(err), zap.String("language", lang), zap.String("data", data))
		h.WriteAnswer(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteAnswer(w, http.StatusOK, response)
}
