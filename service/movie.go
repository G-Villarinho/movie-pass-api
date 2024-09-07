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

	log.Info("Initializing create movie process")

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	movie := payload.ToMovie(session.UserID)

	if err := m.movieRepository.Create(ctx, *movie); err != nil {
		log.Error("Failed to create movie", slog.String("error", err.Error()))
		return nil, domain.ErrCreateMovie
	}

	for _, image := range payload.Images {
		imageBytes, err := utils.ConvertImageToBytes(image)
		if err != nil {
			log.Error("Failed to convert image to bytes", slog.String("error", err.Error()))
			continue
		}

		task := domain.MovieImageUploadTask{
			MovieID: movie.ID,
			Image:   imageBytes,
			UserID:  session.UserID,
		}
		if err := m.movieRepository.AddUploadImageTaskToQueue(ctx, task); err != nil {
			log.Error("Failed to enqueue image upload task", slog.String("error", err.Error()))
		}
	}

	log.Info("Create movie process completed successfully")
	return movie.ToMovieResponse(), nil
}

func (m *movieService) ProcessUploadImageQueue(ctx context.Context) error {
	log := slog.With(
		slog.String("service", "movie"),
		slog.String("func", "processUploadImageQueue"),
	)

	log.Info("Initializing image upload queue processor")

	minInterval := 5 * time.Second
	maxInterval := 1 * time.Minute
	checkInterval := minInterval

	for {
		task, err := m.movieRepository.GetNextUploadImageTaskFromQueue(ctx)
		if err != nil {
			log.Error("Failed to get task from queue", slog.String("error", err.Error()))
			time.Sleep(checkInterval)
			continue
		}

		if task == nil {
			log.Info("No tasks found in the queue, waiting before retrying")

			checkInterval *= 2
			if checkInterval > maxInterval {
				checkInterval = maxInterval
			}

			time.Sleep(checkInterval)
			continue
		}

		checkInterval = minInterval

		log.Info("Processing task", slog.String("movieID", task.MovieID.String()))

		filename := fmt.Sprintf("movie_%s_image_%d.jpg", task.MovieID.String(), time.Now().Unix())
		imageURL, err := m.cloudFlareService.UploadImage(task.Image, filename)
		if err != nil {
			log.Error("Failed to upload image to Cloudflare", slog.String("error", err.Error()))
			continue
		}

		movieImage := domain.MovieImage{
			ID:       uuid.New(),
			MovieID:  task.MovieID,
			ImageURL: imageURL,
		}

		if err := m.movieRepository.CreateMovieImage(ctx, movieImage); err != nil {
			log.Error("Failed to save movie image to the database", slog.String("error", err.Error()))
			continue
		}

		log.Info("Successfully uploaded image to Cloudflare", slog.String("movieID", task.MovieID.String()))
	}
}

func (m *movieService) GetAllByUserID(ctx context.Context) ([]*domain.MovieResponse, error) {
	log := slog.With(
		slog.String("service", "movie"),
		slog.String("func", "GetAllByUserID"),
	)

	log.Info("Initializing get all movie by user id process")

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	movies, err := m.movieRepository.GetALlByUserID(ctx, session.UserID)
	if err != nil {
		log.Error("Failed to get all movies by user id", slog.String("error", err.Error()))
		return nil, domain.ErrGetMoviesByUserID
	}

	if movies == nil {
		log.Warn("movies not found for this user id", slog.String("userID", session.UserID.String()))
		return nil, domain.ErrMoviesNotFoundByUserID
	}

	var moviesResponse []*domain.MovieResponse
	for _, movie := range movies {
		moviesResponse = append(moviesResponse, movie.ToMovieResponse())
	}

	log.Info("Get all movies by user id process executed succefully")
	return moviesResponse, nil
}
