package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo      *repository.SubscriptionRepository
	logger *slog.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, logger *slog.Logger) *SubscriptionService {
	return &SubscriptionService{repo: repo, logger: logger}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	
	if sub.ServiceName == "" {
		return errors.New("service name is required")
	}

	if sub.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	if sub.StartDate.IsZero() {
		return errors.New("start date is required")
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return errors.New("end date must be after start date")
	}

	if sub.ID == uuid.Nil{
		sub.ID = uuid.New()
	}

	err := s.repo.Create(ctx, sub)
	if  err != nil{
		s.logger.Error("failed to create subscription", "error", err)
		return err
	}

	s.logger.Info("subscription successfully created", "id", sub.ID)
	return nil
}
