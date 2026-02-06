package model

import (
	"time"

	. "github.com/google/uuid"
)

type Subscription struct {
	ID          UUID       `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserId      UUID       `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}
