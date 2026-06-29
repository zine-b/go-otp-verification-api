package http

import (
	"encoding/json"
	"errors"
	"net/http"

	appin "prepareGo/internal/application/port/in"
	"prepareGo/internal/domain"
)
type OTPHandler struct {
	useCase appin.OTPUseCase
}

func NewOTPHandler(useCase appin.OTPUseCase) *OTPHandler {
	return &OTPHandler{
		useCase: useCase,
	}
}

type sendOTPRequest struct {
	Phone          string `json:"phone"`
	Channel        string `json:"channel"`
	IdempotencyKey string `json:"idempotency_key"`
}

type verifyOTPRequest struct {
	VerificationID string `json:"verification_id"`
	Code           string `json:"code"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// déclarer les routes HTTP
// Le ServeMux --> le routeur HTTP natif de Go.
// quand il reçois telle URL + telle méthode HTTP,
// il exécute telle fonction

func (h *OTPHandler) RegisterRoutes(mux *http.ServeMux) {
	//Quand le client appelle GET /health --> exécute la méthode h.Health
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /otp/send", h.SendOTP)
	mux.HandleFunc("POST /otp/verify", h.VerifyOTP)
}

func (h *OTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}


func (h *OTPHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req sendOTPRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	// pk useCase.SendOTP ???
	result, err := h.useCase.SendOTP(r.Context(), appin.SendOTPCommand{
		Phone:          req.Phone,
		Channel:        req.Channel,
		IdempotencyKey: req.IdempotencyKey,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPhoneRequired):
			writeError(w, http.StatusBadRequest, err.Error())
			return
	
		case errors.Is(err, domain.ErrRateLimitExceeded):
			writeError(w, http.StatusTooManyRequests, err.Error())
			return
	
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *OTPHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req verifyOTPRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	err := h.useCase.VerifyOTP(r.Context(), appin.VerifyOTPCommand{
		VerificationID: req.VerificationID,
		Code:           req.Code,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrVerificationNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return

		case errors.Is(err, domain.ErrInvalidCode):
			writeError(w, http.StatusUnauthorized, err.Error())
			return

		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "verified",
		"message": "otp verified successfully",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{
		Error: message,
	})
}