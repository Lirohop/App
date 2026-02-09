package service

import (
	"github.com/Lirohop/App/internal/model"
	"github.com/Lirohop/App/internal/repository"
	"github.com/Lirohop/App/internal/utils"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo   *repository.SubscriptionRepository
	logger *slog.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, logger *slog.Logger) *SubscriptionService {
	return &SubscriptionService{repo: repo, logger: logger}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub *model.Subscription) error {

	if sub.ServiceName == "" {
		s.logger.Warn("invalid subscription data", "reason", "service name is empty")
		return errors.New("service name is required")
	}

	if sub.Price <= 0 {
		s.logger.Warn("invalid subscription data", "reason", "price is not valid")
		return errors.New("price must be greater than zero")
	}

	if sub.StartDate.IsZero() {
		s.logger.Warn("invalid subscription data", "reason", "start date is zero")
		return errors.New("start date is required")
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		s.logger.Warn("invalid subscription data", "reason", "end date is not valid")
		return errors.New("end date must be after start date")
	}

	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}

	err := s.repo.Create(ctx, sub)
	if err != nil {
		s.logger.Error("failed to create subscription", "error", err)
		return err
	}

	s.logger.Info("subscription successfully created", "id", sub.ID)
	return nil
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, sub *model.Subscription) error {

	if sub.ID == uuid.Nil {
		s.logger.Warn("invalid subscription data", "reason", "id is nil")
		return errors.New("id is nil")
	}

	if sub.ServiceName == "" {
		s.logger.Warn("invalid subscription data", "reason", "service name is empty")
		return errors.New("service name is required")
	}

	if sub.Price <= 0 {
		s.logger.Warn("invalid subscription data", "reason", "price is not valid")
		return errors.New("price must be greater than zero")
	}

	if sub.StartDate.IsZero() {
		s.logger.Warn("invalid subscription data", "reason", "start date is zero")
		return errors.New("start date is required")
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		s.logger.Warn("invalid subscription data", "reason", "end date is not valid")
		return errors.New("end date must be after start date")
	}

	err := s.repo.Update(ctx, sub)

	if err != nil {
		s.logger.Error("failed to update subscription", "error", err, "id", sub.ID)
		return err
	}

	s.logger.Info("subscription updated", "id", sub.ID)

	return nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {

	if id == uuid.Nil {
		s.logger.Error("id is nil")
		return errors.New("id is valid")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete subscription", "error", err)
		return err
	}

	s.logger.Info("subscription successfully deleted", "id", id)
	return nil
}

func (s *SubscriptionService) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {

	if id == uuid.Nil {
		s.logger.Error("id is nil")
		return nil, errors.New("id is valid")
	}

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get subscription by id", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("subscription fetched", "id", id)

	return sub, err
}

func (s *SubscriptionService) List(ctx context.Context) ([]*model.Subscription, error) {

	subs, err := s.repo.List(ctx)

	if err != nil {
		s.logger.Error("failed to list subscriptions", "error", err)
		return nil, err
	}

	s.logger.Info("subscriptions fetched", "count", len(subs))

	return subs, nil
}

func (s *SubscriptionService) CalculateSubscriptionsTotalCost(
	ctx context.Context,
	userId uuid.UUID,
	serviceName string,
	dateStart time.Time,
	dateEnd time.Time,
) (int, error) {

	sub, err := s.repo.GetByUserAndService(ctx, userId, serviceName)
	if err != nil {
		s.logger.Error("failed to get subscription", "error", err)
		return 0, err
	}

	periodStart := utils.MaxTime(dateStart, sub.StartDate)

	var periodEnd time.Time
	if sub.EndDate != nil {
		periodEnd = utils.MinTime(dateEnd, *sub.EndDate)
	} else {
		periodEnd = dateEnd
	}

	if periodStart.After(periodEnd) {
		s.logger.Info("no overlapping period for subscription",
			"user_id", userId,
			"service", serviceName,
		)
		return 0, nil
	}

	months := (periodEnd.Year()-periodStart.Year())*12 +
		int(periodEnd.Month()-periodStart.Month()) + 1

	totalPrice := months * sub.Price

	s.logger.Debug(
		"calculated subscription cost",
		"user_id", userId,
		"service", serviceName,
		"months", months,
		"total_price", totalPrice,
	)

	return totalPrice, nil
}
