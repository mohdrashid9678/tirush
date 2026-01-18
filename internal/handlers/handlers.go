package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/mohdrashid9678/tirush/internal/models"
	"github.com/mohdrashid9678/tirush/internal/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes defines API endpoints
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/users", h.CreateUser)
	r.Get("/events/{eventID}/seats", h.GetSeats)
	r.Post("/book", h.BookTicket)
}

// CreateUser handles user registration
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	user, err := h.svc.RegisterUser(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// GetSeats retrieves seats for a given event
func (h *Handler) GetSeats(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventID")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid Event ID"})
		return
	}

	seats, err := h.svc.ListSeats(r.Context(), eventID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, seats)
}

// BookTicket handles seat booking requests
func (h *Handler) BookTicket(w http.ResponseWriter, r *http.Request) {
	var req models.BookingRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid payload"})
		return
	}

	// For now userID is passed in the request directly, TODO: Extract from auth once implemented
	if err := h.svc.AttemptBooking(r.Context(), req); err != nil {
		if err.Error() == "seat is already booked or unavailable" {
			render.Status(r, http.StatusConflict) // 409 CONFLICT
			render.JSON(w, r, map[string]string{"error": "Too slow! Seat taken."})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "booked", "seat_id": req.SeatID.String()})
}
