package utils

import (
	"time"

	"github.com/google/uuid"
)

func IsValidUUIDv4(id uuid.UUID) bool {
	return id.Version() == 4 && id != uuid.Nil
}

func ParseUUIDFromString(str string)( uuid.UUID, error){
	id, err := uuid.Parse(str)
	if err != nil{
		return uuid.Nil, err
	}
	return id, nil
}

func MaxTime(A time.Time, B time.Time) time.Time{
	if A.After(B){
		return A
	}
	return B
}
func MinTime(A time.Time, B time.Time) time.Time{
	if A.Before(B){
		return A
	}
	return B	
}

func ParseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		t.Year(),
		t.Month(),
		1,
		0, 0, 0, 0,
		time.UTC,
	), nil
}
