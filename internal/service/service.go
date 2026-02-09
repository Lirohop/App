package service

import (
	"app/internal/model"
	"app/internal/repository"
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
		s.logger.Error("service name is empty")
		return errors.New("service name is required")
	}

	if sub.Price <= 0 {
		s.logger.Error("price is not valid")
		return errors.New("price must be greater than zero")
	}

	if sub.StartDate.IsZero() {
		s.logger.Error("start date is zero")
		return errors.New("start date is required")
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		s.logger.Error("end date is not valid")
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
	err := s.repo.Update(ctx, sub)
	
	if err != nil{
		s.logger.Error("???")
		return err
	}

	s.logger.Info("???")
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

func (s *SubscriptionService) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*model.Subscription ,error) {

	if id == uuid.Nil {
		s.logger.Error("id is nil")
		return nil, errors.New("id is valid")
	}

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil{
		s.logger.Error("???")
		return nil, err
	}

	s.logger.Info("???")

	return sub, err
}

func (s *SubscriptionService) List(ctx context.Context) ([]*model.Subscription, error) {

	subs, err := s.repo.List(ctx)

	if err != nil{
		s.logger.Error("???")
		return nil, err
	}

	s.logger.Info("???")

	return subs, nil
}

func (s *SubscriptionService) CalculateSubscriptionsTotalCost (ctx context.Context, userId uuid.UUID, serviceName string, dateStart time.Time, dateEnd time.Time)(int, error){


	return 0, nil
}
