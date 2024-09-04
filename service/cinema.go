package service

import (
	"context"
	"log/slog"

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
		return nil, err
	}

	return &cinemaService{
		i:                i,
		cinemaRepository: cinemaRepository,
	}, nil
}

func (c *cinemaService) Create(ctx context.Context, payload domain.CinemaPayload) (*domain.CinemaResponse, error) {
	log := slog.With(
		slog.String("service", "cinema"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing cinema creation process")

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	cinema := payload.ToCinema(session.UserID)
	if err := c.cinemaRepository.Create(ctx, *cinema); err != nil {
		log.Error("Failed to create cinema", slog.String("error", err.Error()))
		return nil, domain.ErrCreateCinema
	}

	log.Info("Cinema creation process executed succefully")
	return cinema.ToCinemaResponse(), nil
}

func (c *cinemaService) GetByID(ctx context.Context, cinemaID uuid.UUID) (*domain.CinemaResponse, error) {
	log := slog.With(
		slog.String("service", "cinema"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing get cinema by id process")

	cinema, err := c.cinemaRepository.GetByID(ctx, cinemaID)
	if err != nil {
		log.Error("Failed to get cinema by id", slog.String("error", err.Error()))
		return nil, domain.ErrCreateCinema
	}

	if cinema == nil {
		log.Warn("Cinema not found with this id", slog.String("cinenaID", cinemaID.String()))
		return nil, domain.ErrCinemaNotFound
	}

	log.Info("Get cinema by id process executed succefully")
	return cinema.ToCinemaResponse(), err
}

func (c *cinemaService) GetAll(ctx context.Context) ([]domain.CinemaResponse, error) {
	log := slog.With(
		slog.String("service", "cinema"),
		slog.String("func", "GetAll"),
	)

	log.Info("Initializing get all cinemas process")

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	cinemas, err := c.cinemaRepository.GetAll(ctx, session.UserID)
	if err != nil {
		log.Error("Failed to get all cinemas", slog.String("error", err.Error()))
		return nil, domain.ErrGetCinema
	}

	if cinemas == nil {
		log.Warn("No cinemas found")
		return nil, domain.ErrCinemaNotFound
	}

	var cinemasResponse []domain.CinemaResponse
	for _, cinema := range cinemas {
		cinemasResponse = append(cinemasResponse, *cinema.ToCinemaResponse())
	}

	log.Info("Get all cinemas process executed successfully")
	return cinemasResponse, nil
}

func (c *cinemaService) Delete(ctx context.Context, cinemaID uuid.UUID) error {
	log := slog.With(
		slog.String("service", "cinema"),
		slog.String("func", "Delete"),
	)

	log.Info("Initializing cinema deletion process")

	cinema, err := c.cinemaRepository.GetByID(ctx, cinemaID)
	if err != nil {
		log.Error("Failed to get cinema by ID", slog.String("error", err.Error()))
		return domain.ErrGetCinema
	}

	if cinema == nil {
		log.Warn("No cinema found with the provided ID")
		return domain.ErrCinemaNotFound
	}

	if err := c.cinemaRepository.Delete(ctx, cinemaID); err != nil {
		log.Error("Failed to delete cinema", slog.String("error", err.Error()))
		return domain.ErrDeleteCinema
	}

	log.Info("Cinema deletion process executed successfully")
	return nil
}
