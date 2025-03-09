package models

import (
	// Local Packages
	"learn-go/errors"
)

type Order struct {
	ID          string     `json:"order_id"`
	UserID      string     `json:"user_id"`
	LineItems   []LineItem `json:"line_items"`
	OrderStatus string     `json:"order_status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	ShippedAt   string     `json:"shipped_at"`
	DeliveredAt string     `json:"delivered_at"`
}

type LineItem struct {
	ItemID   string  `json:"item_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (o *Order) ValidateCreation() error {
	ve := errors.ValidationErrs()

	if o.ID != "" {
		ve.Add("order_id", "must be empty during creation")
	}
	validateOrderFields(o, ve)

	if o.ShippedAt == "" {
		o.ShippedAt = "Will Be Shipped Soon"
	}
	if o.DeliveredAt == "" {
		o.DeliveredAt = "Will Be Delivered Soon"
	}

	return ve.Err()
}

func (o *Order) ValidateUpdate(orderID string) error {
	ve := errors.ValidationErrs()

	if o.ID != orderID {
		ve.Add("order_id", "does not match the existing order")
	}
	validateOrderFields(o, ve)

	return ve.Err()
}

func validateOrderFields(o *Order, ve *errors.ValidationErrorBuilder) {
	if o.UserID == "" {
		ve.Add("user_id", "cannot be empty")
	}
	if len(o.LineItems) == 0 {
		ve.Add("line_items", "cannot be empty")
	}
	if o.OrderStatus == "" {
		ve.Add("order_status", "cannot be empty")
	}
	for _, item := range o.LineItems {
		if item.ItemID == "" {
			ve.Add("item_id", "cannot be empty")
		}
		if item.Quantity <= 0 {
			ve.Add("quantity", "must be greater than zero")
		}
		if item.Price <= 0 {
			ve.Add("price", "must be greater than zero")
		}
	}
}
