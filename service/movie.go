package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/GSVillas/movie-pass-api/client"
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/GSVillas/movie-pass-api/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type movieService struct {
	i                 *do.Injector
	movieRepository   domain.MovieRepository
	cloudFlareService client.CloudFlareService
}

func NewMovieService(i *do.Injector) (domain.MovieService, error) {
	movieRepository, err := do.Invoke[domain.MovieRepository](i)
	if err != nil {
		return nil, err
	}

	cloudFlareService, err := do.Invoke[client.CloudFlareService](i)
	if err != nil {
		return nil, err
	}

	return &movieService{
		i:                 i,
		movieRepository:   movieRepository,
		cloudFlareService: cloudFlareService,
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

func (m *movieService) Create(ctx context.Context, payload domain.MoviePayload) (*domain.MovieResponse, error) {
	log := slog.With(
		slog.String("service", "movie"),
		slog.String("func", "create"),
	)

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	movie := payload.ToMovie(session.UserID)

	if err := m.movieRepository.Create(ctx, *movie); err != nil {
		return nil, fmt.Errorf("error to create movie %w", err)
	}

	for _, image := range payload.Images {
		imageBytes, err := utils.ConvertImageToBytes(image)
		if err != nil {
			log.Error("error to convert image to bytes", slog.String("error", err.Error()))
			continue
		}

		task := domain.MovieImageUploadTask{
			MovieID: movie.ID,
			Image:   imageBytes,
			UserID:  session.UserID,
		}

		if err := m.movieRepository.AddUploadImageTaskToQueue(ctx, task); err != nil {
			log.Error(err.Error())
		}
	}

	return movie.ToMovieResponse(), nil
}

func (m *movieService) ProcessUploadImageQueue(ctx context.Context, task domain.MovieImageUploadTask) error {
	filename := fmt.Sprintf("movie_%s_image_%d.jpg", task.MovieID.String(), time.Now().Unix())
	imageURL, err := m.cloudFlareService.UploadImage(task.Image, filename)
	if err != nil {
		return fmt.Errorf("error to upload image to Cloudflare %w", err)
	}

	movieImage := domain.MovieImage{
		ID:       uuid.New(),
		MovieID:  task.MovieID,
		ImageURL: imageURL,
	}

	if err := m.movieRepository.CreateMovieImage(ctx, movieImage); err != nil {
		return fmt.Errorf("error to save movie image to the database error:%w", err)
	}

	return nil
}

func (m *movieService) GetAllByUserID(ctx context.Context) ([]*domain.MovieResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	movies, err := m.movieRepository.GetALlByUserID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("error to get all movies by user id. error: %w", err)
	}

	if movies == nil {
		return nil, domain.ErrMoviesNotFoundByUserID
	}

	var moviesResponse []*domain.MovieResponse
	for _, movie := range movies {
		moviesResponse = append(moviesResponse, movie.ToMovieResponse())
	}

	return moviesResponse, nil
}

func (m *movieService) Update(ctx context.Context, movieID uuid.UUID, payload domain.MovieUpdatePayload) (*domain.MovieResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	movie, err := m.movieRepository.GetByID(ctx, movieID, true)
	if err != nil {
		return nil, fmt.Errorf("error to get all movies by id. Error: %w", err)
	}

	if movie == nil {
		return nil, domain.ErrMoviesNotFound
	}

	if movie.UserID != session.UserID {
		return nil, domain.ErrMovieNotBelongUser
	}

	indicativeRatings, err := m.movieRepository.GetAllIndicativeRating(ctx)
	if err != nil {
		return nil, fmt.Errorf("err to get all indicative rating: Error: %w", err)
	}

	if indicativeRatings == nil {
		return nil, domain.ErrIndicativeRatingsNotFound
	}

	var indicativeRating domain.IndicativeRating
	if payload.IndicativeRatingID != nil {
		exists := false
		for _, rating := range indicativeRatings {
			if rating.ID == *payload.IndicativeRatingID {
				indicativeRating = *rating
				exists = true
				break
			}
		}
		if !exists {
			return nil, domain.ErrIndicativeRatingNotFound
		}
	}

	updates := map[string]any{}
	if payload.IndicativeRatingID != nil {
		movie.IndicativeRating = indicativeRating
		updates["indicativeRatingId"] = *payload.IndicativeRatingID
	}

	if payload.Title != nil {
		movie.Title = *payload.Title
		updates["title"] = *payload.Title
	}

	if payload.Duration != nil {
		movie.Duration = *payload.Duration
		updates["duration"] = *payload.Duration
	}

	if err := m.movieRepository.Update(ctx, movieID, updates); err != nil {
		return nil, fmt.Errorf("error to update movie. error: %w", err)
	}

	return movie.ToMovieResponse(), nil
}

func (m *movieService) Delete(ctx context.Context, movieID uuid.UUID) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return domain.ErrUserNotFoundInContext
	}

	movie, err := m.movieRepository.GetByID(ctx, session.UserID, false)
	if err != nil {
		return fmt.Errorf("error to get movies by user id error:%w", err)
	}

	if movie == nil {
		return domain.ErrMoviesNotFound
	}

	if movie.UserID != session.UserID {
		return domain.ErrMovieNotBelongUser
	}

	return nil
}

func (m *movieService) ProcessDeleteImageQueue(ctx context.Context, task domain.MovieImageDeleteTask) error {
	panic("unimplemented")
}
