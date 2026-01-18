package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never return password in JSON
	FullName     string    `json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type Event struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	TotalSeats     int       `json:"total_seats"`
	AvailableSeats int       `json:"available_seats"` // Derived or Cached
	Date           time.Time `json:"date"`
}

type Seat struct {
	ID         uuid.UUID  `json:"id"`
	EventID    uuid.UUID  `json:"event_id"`
	Section    string     `json:"section"`
	RowNumber  string     `json:"row_number"`
	SeatNumber string     `json:"seat_number"`
	Status     string     `json:"status"` // AVAILABLE, BOOKED
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Version    int        `json:"version"` // Optimistic Lock
}

type BookingRequest struct {
	EventID uuid.UUID `json:"event_id" binding:"required"`
	SeatID  uuid.UUID `json:"seat_id" binding:"required"`
	UserID  uuid.UUID `json:"user_id"`
}
