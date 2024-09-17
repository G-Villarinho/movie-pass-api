package repository

import (
	"math"

	"github.com/GSVillas/movie-pass-api/domain"
	"gorm.io/gorm"
)

func paginate(value any, pagination *domain.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64

	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows

	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
