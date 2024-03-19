package Types

import (
	"database/sql"
	"github.com/volatiletech/null/v8"
	"time"
)

type NullString struct {
	sql.NullString
}
type NullInt64 struct {
	sql.NullInt64
}
type NullFloat64 struct {
	sql.NullFloat64
}
type NullBool struct {
	sql.NullBool
}
type NullTime struct {
	sql.NullTime
}

func NewNullString(s string) NullString {
	if len(s) == 0 {
		return NullString{}
	}
	return NullString{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

func NewNullInt64(n int64) NullInt64 {
	if n == 0 {
		return NullInt64{}
	}
	return NullInt64{
		sql.NullInt64{
			Int64: n,
			Valid: true,
		},
	}
}

func NewNullFloat64(n float64) NullFloat64 {
	if n == 0 {
		return NullFloat64{}
	}
	return NullFloat64{
		sql.NullFloat64{
			Float64: n,
			Valid:   true,
		},
	}
}

func NewNullBool(b bool) NullBool {
	if !b {
		return NullBool{}
	}
	return NullBool{
		sql.NullBool{
			Bool:  b,
			Valid: true,
		},
	}
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{
		sql.NullTime(null.Time{
			Time:  t,
			Valid: true,
		}),
	}
}
