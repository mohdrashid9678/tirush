package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohdrashid9678/tirush/internal/models"
	"github.com/mohdrashid9678/tirush/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// RegisterUser creates a new user with hashed password
func (s *Service) RegisterUser(ctx context.Context, email, password, fullName string) (*models.User, error) {
	// Hash Password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPwd),
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ListSeats fetches seats for an event
func (s *Service) ListSeats(ctx context.Context, eventID uuid.UUID) ([]models.Seat, error) {
	return s.repo.GetSeatsByEvent(ctx, eventID)
}

// AttemptBooking tries to book a seat for a user, throwing error if unsuccessful
func (s *Service) AttemptBooking(ctx context.Context, req models.BookingRequest) error {
	return s.repo.BookSeat(ctx, req.SeatID, req.UserID)
}
