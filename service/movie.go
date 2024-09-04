package service

import (
	"context"
	"log/slog"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/samber/do"
)

type movieService struct {
	i               *do.Injector
	movieRepository domain.MovieRepository
}

func NewMovieService(i *do.Injector) (domain.MovieService, error) {
	movieRepository, err := do.Invoke[domain.MovieRepository](i)
	if err != nil {
		return nil, err
	}

	return &movieService{
		i:               i,
		movieRepository: movieRepository,
	}, nil
}

func (m *movieService) GetAllIndicativeRating(ctx context.Context) ([]domain.IndicativeRatingResponse, error) {
	log := slog.With(
		slog.String("service", "movie"),
		slog.String("func", "GetAllIndicativeRating"),
	)

	log.Info("Initializing get all indicative rating process")

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	indicativeRating, err := m.movieRepository.GetAllIndicativeRating(ctx, session.UserID)
	if err != nil {
		log.Error("Failed to get all indicative rating", slog.String("error", err.Error()))
		return nil, domain.ErrGetAllIndicativeRating
	}

	if indicativeRating == nil {
		log.Warn("indicative ratings not found")
		return nil, domain.ErrIndicativeRatingsNotFound
	}
}
