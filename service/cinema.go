package service

import (
	"context"
	"fmt"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type cinemaService struct {
	i                *do.Injector
	cinemaRepository domain.CinemaRepository
}

func NewCinemaSevice(i *do.Injector) (domain.CinemaService, error) {
	cinemaRepository, err := do.Invoke[domain.CinemaRepository](i)
	if err != nil {
		return nil, fmt.Errorf("error to initialize CinemaRepository: %w", err)
	}

	return &cinemaService{
		i:                i,
		cinemaRepository: cinemaRepository,
	}, nil
}

func (c *cinemaService) Create(ctx context.Context, payload domain.CinemaPayload) (*domain.CinemaResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	cinema := payload.ToCinema(session.UserID)
	if err := c.cinemaRepository.Create(ctx, *cinema); err != nil {
		return nil, fmt.Errorf("error to create cinema for user ID %s: %w", session.UserID, err)
	}

	return cinema.ToCinemaResponse(), nil
}

func (c *cinemaService) GetByID(ctx context.Context, cinemaID uuid.UUID) (*domain.CinemaResponse, error) {
	cinema, err := c.cinemaRepository.GetByID(ctx, cinemaID)
	if err != nil {
		return nil, fmt.Errorf("error fetching cinema by ID %s: %w", cinemaID.String(), err)
	}

	if cinema == nil {
		return nil, domain.ErrCinemaNotFound
	}

	return cinema.ToCinemaResponse(), nil
}

func (c *cinemaService) GetAll(ctx context.Context, pagination *domain.Pagination) (*domain.Pagination, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	cinemasPagination, err := c.cinemaRepository.GetAll(ctx, session.UserID, pagination)
	if err != nil {
		return nil, fmt.Errorf("error to fetch cinemas for user ID %s: %w", session.UserID, err)
	}

	if cinemasPagination.Rows == nil {
		return nil, domain.ErrCinemaNotFound
	}

	var cinemasResponse []domain.CinemaResponse
	for _, cinema := range cinemasPagination.Rows.([]domain.Cinema) {
		cinemasResponse = append(cinemasResponse, *cinema.ToCinemaResponse())
	}

	cinemasPagination.Rows = cinemasResponse

	return cinemasPagination, nil
}

func (c *cinemaService) Delete(ctx context.Context, cinemaID uuid.UUID) error {
	cinema, err := c.cinemaRepository.GetByID(ctx, cinemaID)
	if err != nil {
		return fmt.Errorf("error to retrieve cinema by ID %s: %w", cinemaID.String(), err)
	}

	if cinema == nil {
		return domain.ErrCinemaNotFound
	}

	if err := c.cinemaRepository.Delete(ctx, cinemaID); err != nil {
		return fmt.Errorf("error to delete cinema with ID %s: %w", cinemaID.String(), err)
	}

	return nil
}
