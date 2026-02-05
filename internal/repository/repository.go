package repository

import (
	"app/internal/model"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository struct {
	db *pgx.Conn
	logger *slog.Logger
}

func NewSubscriptionRepository(db *pgx.Conn, logger *slog.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{db: db, logger: logger}
}

func (r *SubscriptionRepository) Craete(ctx context.Context, s *model.Subscriptions) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO subscriptions(id, service_name, price, user_id, start_date, end_date)
		 VALUES($1, $2, $3, $4, $5, $6)`,
		s.ID, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate)

	if err != nil{
		r.logger.Error("failed to create subscription", "error", err, "user_id", s.UserId)
		return err
	}

	r.logger.Info("subscription created", "id", s.ID)
	
	return nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*model.Subscriptions, error) {
	var s model.Subscriptions

	err := r.db.QueryRow(ctx,
		`SELECT id, sevice_name, price, user_id, start_date, end_date
	 From subscriptions
	 Where id=$1`, id).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price, &s.UserId,
		&s.StartDate,
		&s.EndDate)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]*model.Subscriptions, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subs []*model.Subscriptions
	for rows.Next() {
		var s model.Subscriptions
		if err := rows.Scan(&s.ID, s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}

	return subs, nil

}
