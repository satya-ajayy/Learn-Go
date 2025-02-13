package redis

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"fmt"

	// Local Packages
	errors "learn-go/errors"
	omodels "learn-go/models/orders"

	// External Packages
	"github.com/redis/go-redis/v9"
)

type OrdersRepository struct {
	client *redis.Client
}

func NewOrdersRepository(client *redis.Client) *OrdersRepository {
	return &OrdersRepository{client: client}
}

func orderIDKey(id string) string {
	return fmt.Sprintf("order:%s", id)
}

func (r *OrdersRepository) Insert(ctx context.Context, order omodels.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	tx := r.client.TxPipeline()
	key := orderIDKey(order.OrderID)
	res := tx.SetNX(ctx, key, data, 0)
	if err := res.Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err := tx.SAdd(ctx, "orders", key).Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to add order to set: %w", err)
	}

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	return nil
}

func (r *OrdersRepository) Get(ctx context.Context, orderID string) (omodels.Order, error) {
	key := orderIDKey(orderID)
	value, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return omodels.Order{}, errors.E(errors.NotFound, "order not found")
	}
	if err != nil {
		return omodels.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	var order omodels.Order
	if err := json.Unmarshal([]byte(value), &order); err != nil {
		return omodels.Order{}, fmt.Errorf("failed to decode order: %w", err)
	}
	return order, nil
}

func (r *OrdersRepository) Delete(ctx context.Context, orderID string) error {
	key := orderIDKey(orderID)
	tx := r.client.TxPipeline()

	err := tx.Del(ctx, key).Err()

	if errors.Is(err, redis.Nil) {
		tx.Discard()
		return errors.E(errors.NotFound, "order not found")
	}
	if err != nil {
		tx.Discard()
		return fmt.Errorf("failed to get order: %w", err)
	}
	if err := tx.SRem(ctx, "orders", key).Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to remove order from set: %w", err)
	}

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

func (r *OrdersRepository) Update(ctx context.Context, order omodels.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)
	err = r.client.SetXX(ctx, key, string(data), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

func (r *OrdersRepository) Exists(ctx context.Context, orderID string) (bool, error) {
	res, err := r.client.Exists(ctx, orderIDKey(orderID)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check order: %w", err)
	}
	return res == 1, nil
}
