package handlers

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// Local Packages
	errors "learn-go/errors"
	omodels "learn-go/models/orders"

	// External Packages
	"github.com/go-chi/chi/v5"
)

type OrdersService interface {
	Insert(ctx context.Context, order omodels.Order) (omodels.Order, error)
	Get(ctx context.Context, orderID string) (omodels.Order, error)
	Update(ctx context.Context, order omodels.Order) error
	Delete(ctx context.Context, orderID string) error
}

type OrdersHandler struct {
	svc OrdersService
}

func NewOrdersHandler(orderService OrdersService) *OrdersHandler {
	return &OrdersHandler{svc: orderService}
}

func (a *OrdersHandler) GetOne(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		return nil, http.StatusBadRequest, errors.EmptyParamErr("orderId")
	}

	order, err := a.svc.Get(r.Context(), orderID)
	if err == nil {
		return order, http.StatusOK, nil
	}
	return
}

func (a *OrdersHandler) Insert(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var order omodels.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return nil, http.StatusBadRequest, errors.InvalidBodyErr(err)
	}
	if err := order.ValidateCreation(); err != nil {
		return nil, http.StatusBadRequest, errors.ValidationFailedErr(err)
	}

	order, err = a.svc.Insert(r.Context(), order)
	if err == nil {
		return order, http.StatusCreated, nil
	}
	return
}

func (a *OrdersHandler) Update(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		return nil, http.StatusBadRequest, errors.EmptyParamErr("orderId")
	}

	var updatedOrder omodels.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		return nil, http.StatusBadRequest, errors.InvalidBodyErr(err)
	}

	if err := updatedOrder.ValidateUpdate(orderID); err != nil {
		return nil, http.StatusBadRequest, errors.ValidationFailedErr(err)
	}

	err = a.svc.Update(r.Context(), updatedOrder)
	if err == nil {
		return updatedOrder, http.StatusOK, nil
	}
	return
}

func (a *OrdersHandler) Delete(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		return nil, http.StatusBadRequest, errors.EmptyParamErr("orderId")
	}

	err = a.svc.Delete(r.Context(), orderID)
	if err == nil {
		return map[string]string{"message": fmt.Sprintf("%s is deleted", orderID)}, http.StatusOK, nil
	}
	return
}
