package repository

import (
	"app/internal/model"
	"context"
	"log/slog"

	. "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewSubscriptionRepository(db *pgxpool.Pool, logger *slog.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{db: db, logger: logger}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *model.Subscription) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO subscriptions(id, service_name, price, user_id, start_date, end_date)
		 VALUES($1, $2, $3, $4, $5, $6)`,
		s.ID, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate)

	if err != nil {
		r.logger.Error("failed to create subscription", "error", err, "user_id", s.UserId)
		return err
	}

	r.logger.Info("subscription created", "id", s.ID)

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		r.logger.Error("failed to delete subscription", "error", err, "id", id)
		return err
	}

	r.logger.Info("subscription was deleted", "id", id)

	return nil
}

func (r *SubscriptionRepository) Update(
	ctx context.Context,
	s *model.Subscription,
) error {
	_, err := r.db.Exec(ctx,
		`UPDATE subscriptions
         SET service_name = $1,
             price = $2,
             user_id = $3,
             start_date = $4,
             end_date = $5
         WHERE id = $6`,
		s.ServiceName,
		s.Price,
		s.UserId,
		s.StartDate,
		s.EndDate,
		s.ID,
	)

	if err != nil {
		r.logger.Error("failed to update subscription", "error", err, "id", s.ID)
		return err
	}

	r.logger.Info("subscription was updated", "id", s.ID)

	return nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*model.Subscription, error) {
	var s model.Subscription

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

	r.logger.Info("subscription finded by id", "id", s.ID)

	return &s, nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]*model.Subscription, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subs []*model.Subscription
	for rows.Next() {
		var s model.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserId,
			&s.StartDate,
			&s.EndDate,
		); err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}

	r.logger.Info("subscriptions were found", slog.Int("count", len(subs)))

	return subs, nil

}
