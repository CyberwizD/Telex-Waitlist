package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/CyberwizD/Telex-Waitlist/internal/services"
	"github.com/CyberwizD/Telex-Waitlist/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type WaitlistHandler struct {
	service     services.WaitlistService
	adminAPIKey string
}

func NewWaitlistHandler(service services.WaitlistService, adminAPIKey string) *WaitlistHandler {
	return &WaitlistHandler{service: service, adminAPIKey: adminAPIKey}
}

type waitlistRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// Submit handles public waitlist submissions.
func (h *WaitlistHandler) Submit(c *gin.Context) {
	var req waitlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.Submit(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			response.JSONError(c, http.StatusConflict, "email already on waitlist")
			return
		}
		if errors.Is(err, services.ErrValidation) {
			msg := sanitizeValidationMessage(err)
			response.JSONError(c, http.StatusBadRequest, msg)
			return
		}
		response.JSONError(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSONData(c, http.StatusCreated, entry)
}

// List returns paginated entries and is protected by an admin token.
func (h *WaitlistHandler) List(c *gin.Context) {
	if h.adminAPIKey == "" {
		response.JSONError(c, http.StatusForbidden, "admin listing disabled")
		return
	}
	if c.GetHeader("X-Admin-Token") != h.adminAPIKey {
		response.JSONError(c, http.StatusUnauthorized, "invalid admin token")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	entries, total, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSONPage(c, http.StatusOK, entries, total, limit, offset)
}

func sanitizeValidationMessage(err error) string {
	msg := err.Error()
	prefix := services.ErrValidation.Error() + ": "
	msg = strings.TrimPrefix(msg, prefix)
	if msg == "" || msg == services.ErrValidation.Error() {
		return "invalid input"
	}
	return msg
}
