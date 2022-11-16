package utils

import "database/sql"

func NullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
