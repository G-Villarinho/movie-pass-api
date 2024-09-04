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

func (m *movieService) GetAllIndicativeRating(ctx context.Context) ([]*domain.IndicativeRatingResponse, error) {
	log := slog.With(
		slog.String("service", "movie"),
		slog.String("func", "GetAllIndicativeRating"),
	)

	log.Info("Initializing get all indicative rating process")

	indicativeRatings, err := m.movieRepository.GetAllIndicativeRating(ctx)
	if err != nil {
		log.Error("Failed to get all indicative rating", slog.String("error", err.Error()))
		return nil, domain.ErrGetAllIndicativeRating
	}

	if indicativeRatings == nil {
		log.Warn("indicative ratings not found")
		return nil, domain.ErrIndicativeRatingsNotFound
	}

	var indicativeRatingsResponse []*domain.IndicativeRatingResponse
	for _, indicativeRattings := range indicativeRatings {
		indicativeRatingsResponse = append(indicativeRatingsResponse, indicativeRattings.ToIndicativeRatingResponse())
	}

	log.Info("Get all indicative rating process executed succefully")
	return indicativeRatingsResponse, nil
}
