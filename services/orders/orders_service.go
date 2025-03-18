package orders

import (
	// Go Internal Packages
	"context"
	"fmt"

	// Local Packages
	errors "learn-go/errors"
	models "learn-go/models"
	utils "learn-go/utils"
)

type OrdersRepository interface {
	Insert(ctx context.Context, order models.Order) error
	GetOne(ctx context.Context, orderID string) (models.Order, error)
	Update(ctx context.Context, order models.Order) error
	Delete(ctx context.Context, orderID string) error
	Exists(ctx context.Context, orderID string) (bool, error)
}

type OrdersService struct {
	ordersRepository OrdersRepository
}

func NewService(ordersRepository OrdersRepository) *OrdersService {
	return &OrdersService{ordersRepository: ordersRepository}
}

func (s *OrdersService) Insert(ctx context.Context, order models.Order) (string, error) {
	order.ID = utils.GenerateRandomID()
	currTime := utils.GetCurrentTime()
	order.CreatedAt = currTime
	order.UpdatedAt = currTime
	err := s.ordersRepository.Insert(ctx, order)
	return order.ID, err
}

func (s *OrdersService) GetOne(ctx context.Context, orderID string) (models.Order, error) {
	return s.ordersRepository.GetOne(ctx, orderID)
}

func (s *OrdersService) Update(ctx context.Context, order models.Order) error {
	exists, err := s.ordersRepository.Exists(ctx, order.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.E(errors.NotFound, fmt.Sprintf("order not found with id %s", order.ID))
	}
	order.UpdatedAt = utils.GetCurrentTime()
	return s.ordersRepository.Update(ctx, order)
}

func (s *OrdersService) Delete(ctx context.Context, orderID string) error {
	return s.ordersRepository.Delete(ctx, orderID)
}
