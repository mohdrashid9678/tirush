package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohdrashid9678/tirush/internal/models"
)

var ErrSeatUnavailable = errors.New("seat is already booked or unavailable")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a new user
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, email, password_hash, full_name, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.FullName, user.CreatedAt)
	return err
}

// GetSeatsByEvent returns all seats for an event
func (r *Repository) GetSeatsByEvent(ctx context.Context, eventID uuid.UUID) ([]models.Seat, error) {
	query := `SELECT id, event_id, section, row_number, seat_number, status, version FROM seats WHERE event_id = $1 ORDER BY section, row_number, seat_number`
	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var s models.Seat
		if err := rows.Scan(&s.ID, &s.EventID, &s.Section, &s.RowNumber, &s.SeatNumber, &s.Status, &s.Version); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

// BookSeat performs the Transactional Booking with Optimistic Locking
func (r *Repository) BookSeat(ctx context.Context, seatID, userID uuid.UUID) error {
	// 1. Start a Transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Always Rollback. If Commit is called, Rollback is a no-op.
	defer tx.Rollback(ctx)

	// 2. Optimistic Locking Query
	// We only update IF the seat is AVAILABLE (Condition 1)
	// AND the version hasn't changed (Condition 2 - implied by status check here, but version adds safety)
	query := `
		UPDATE seats 
		SET status = 'BOOKED', user_id = $1, version = version + 1
		WHERE id = $2 AND status = 'AVAILABLE'
	`

	cmdTag, err := tx.Exec(ctx, query, userID, seatID)
	if err != nil {
		return err
	}

	// 3. Check if any row was actually updated
	if cmdTag.RowsAffected() == 0 {
		return ErrSeatUnavailable
	}

	// 4. Create the Booking Record
	bookingQuery := `INSERT INTO bookings (id, user_id, event_id, seat_id, status) 
	                 VALUES ($1, $2, (SELECT event_id FROM seats WHERE id = $3), $3, 'CONFIRMED')`

	_, err = tx.Exec(ctx, bookingQuery, uuid.New(), userID, seatID)
	if err != nil {
		return err
	}

	// 5. Commit
	return tx.Commit(ctx)
}
