package orders

import (
	// Go Internal Packages
	"context"
	"fmt"

	// Local Packages
	errors "learn-go/errors"
	omodels "learn-go/models/orders"
	helpers "learn-go/utils/helpers"
)

type OrdersRepository interface {
	Insert(ctx context.Context, order omodels.Order) error
	Get(ctx context.Context, orderID string) (omodels.Order, error)
	Update(ctx context.Context, order omodels.Order) error
	Delete(ctx context.Context, orderID string) error
	Exists(ctx context.Context, orderID string) (bool, error)
}

type OrdersService struct {
	ordersRepository OrdersRepository
}

func NewService(ordersRepository OrdersRepository) *OrdersService {
	return &OrdersService{ordersRepository: ordersRepository}
}

func (s *OrdersService) Insert(ctx context.Context, order omodels.Order) (omodels.Order, error) {
	order.OrderID = helpers.GenerateRandomID()
	currTime := helpers.GetCurrentTime()
	order.CreatedAt = currTime
	order.UpdatedAt = currTime
	err := s.ordersRepository.Insert(ctx, order)
	return order, err
}

func (s *OrdersService) Get(ctx context.Context, orderID string) (omodels.Order, error) {
	return s.ordersRepository.Get(ctx, orderID)
}

func (s *OrdersService) Update(ctx context.Context, order omodels.Order) error {
	exists, err := s.ordersRepository.Exists(ctx, order.OrderID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.E(errors.NotFound, fmt.Sprintf("order not found with id %s", order.OrderID))
	}
	order.UpdatedAt = helpers.GetCurrentTime()
	return s.ordersRepository.Update(ctx, order)
}

func (s *OrdersService) Delete(ctx context.Context, orderID string) error {
	return s.ordersRepository.Delete(ctx, orderID)
}
