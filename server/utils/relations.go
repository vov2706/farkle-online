package utils

import "gorm.io/gorm"

func ApplyRelations[T any](query gorm.ChainInterface[T], relations []string) gorm.ChainInterface[T] {
	for _, rel := range relations {
		if rel == "" {
			continue
		}
		query = query.Preload(rel, nil)
	}

	return query
}
