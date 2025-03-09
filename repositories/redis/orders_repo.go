package redis

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"fmt"

	// Local Packages
	errors "learn-go/errors"
	models "learn-go/models"
	helpers "learn-go/utils/helpers"

	// External Packages
	"github.com/redis/go-redis/v9"
)

type OrdersRepository struct {
	client *redis.Client
}

func NewOrdersRepository(client *redis.Client) *OrdersRepository {
	return &OrdersRepository{client: client}
}

func (r *OrdersRepository) GetOne(ctx context.Context, orderID string) (models.Order, error) {
	key := helpers.GetOrderID(orderID)
	value, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return models.Order{}, errors.E(errors.NotFound, "order not found")
	}
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	var order models.Order
	if err := json.Unmarshal([]byte(value), &order); err != nil {
		return models.Order{}, fmt.Errorf("failed to decode order: %w", err)
	}
	return order, nil
}

func (r *OrdersRepository) Insert(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	tx := r.client.TxPipeline()
	key := helpers.GetOrderID(order.ID)

	res := tx.SetNX(ctx, key, data, 0)
	if err := res.Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err := tx.SAdd(ctx, "ORDERS", key).Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to add order to set: %w", err)
	}

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	return nil
}

func (r *OrdersRepository) Update(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := helpers.GetOrderID(order.ID)
	err = r.client.SetXX(ctx, key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

func (r *OrdersRepository) Delete(ctx context.Context, orderID string) error {
	key := helpers.GetOrderID(orderID)
	tx := r.client.TxPipeline()

	err := tx.Del(ctx, key).Err()
	if err != nil {
		tx.Discard()
		return fmt.Errorf("failed to get order: %w", err)
	}
	if err := tx.SRem(ctx, "ORDERS", key).Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to remove order from set: %w", err)
	}

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

func (r *OrdersRepository) Exists(ctx context.Context, orderID string) (bool, error) {
	key := helpers.GetOrderID(orderID)
	res, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check order: %w", err)
	}
	return res == 1, nil
}
