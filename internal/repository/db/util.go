package db

import "database/sql"

// ToNullString маленькая утилита-хелпер для чистой конвертации
func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
