package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/google/uuid"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type seatRepository struct {
	i  *do.Injector
	db *gorm.DB
}

func NewSeatRepository(i *do.Injector) (domain.SeatRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	return &seatRepository{
		i:  i,
		db: db,
	}, nil
}

func (s *seatRepository) Create(ctx context.Context, seat domain.Seat) error {
	log := slog.With(
		slog.String("repository", "seat"),
		slog.String("func", "Create"),
	)

	log.Info("Method Initiated")

	if err := s.db.WithContext(ctx).Create(&seat).Error; err != nil {
		log.Error("Failed to create seat", slog.String("error", err.Error()))
		return err
	}

	log.Info("Seat created succefully")

	return nil
}

func (s *seatRepository) GetByID(ctx context.Context, seatID uuid.UUID) (*domain.Seat, error) {
	log := slog.With(
		slog.String("repository", "seat"),
		slog.String("func", "GetByID"),
	)

	log.Info("Method Initiated")

	var seat domain.Seat
	if err := s.db.WithContext(ctx).Where("id = ?", seatID).First(&seat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No seats where found with the provided ID", slog.String("error", err.Error()))
			return nil, nil
		}

		log.Error("Failed to get seat by id", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Seat got by id succefully")

	return &seat, nil
}

func (s *seatRepository) GetAll(ctx context.Context) ([]domain.Seat, error) {
	log := slog.With(
		slog.String("repository", "seat"),
		slog.String("func", "GetAll"),
	)

	log.Info("Method Initiated")

	var seats []domain.Seat
	if err := s.db.WithContext(ctx).Find(&seats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No seats records where found", slog.String("error", err.Error()))
			return nil, nil
		}

		log.Error("Failed to get all seats", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("All seats got succefully")

	return seats, nil
}

func (s *seatRepository) Delete(ctx context.Context, seatID uuid.UUID) error {
	log := slog.With(
		slog.String("repository", "seat"),
		slog.String("func", "Delete"),
	)

	log.Info("Method initiated")

	if err := s.db.WithContext(ctx).Where("id = ?", seatID).Delete(&domain.Seat{}).Error; err != nil {
		log.Error("Failed to delete seat with provided ID", slog.String("error", err.Error()))
		return err
	}

	log.Info("Seat deleted sucefully")

	return nil
}
