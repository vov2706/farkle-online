package database

import (
	"app/http/responses"
	"strings"

	"gorm.io/gorm"
)

func Paginate(page, perPage uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if perPage > responses.MaxPerPage {
			perPage = responses.DefaultPerPage
		}

		offset := (page - 1) * perPage

		return db.Offset(int(offset)).Limit(int(perPage))
	}
}

func OrderBy(column, direction string) func(db *gorm.DB) *gorm.DB {
	dir := strings.ToLower(direction)

	if dir != "asc" && dir != "desc" {
		direction = "asc"
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Order(column + " " + direction)
	}
}

func OrderByDesc(column string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(column + " desc")
	}
}
