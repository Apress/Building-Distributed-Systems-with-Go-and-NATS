package domain

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending          OrderStatus = "Pending"
	OrderStatusPreparing        OrderStatus = "Preparing"
	OrderStatusReadyForDelivery OrderStatus = "Ready for Delivery"
	OrderStatusOutForDelivery   OrderStatus = "Out For Delivery"
	OrderStatusDelivered        OrderStatus = "Delivered"
	OrderStatusInvalid          OrderStatus = "Invalid"
)

// Order is an Entity because it has an identity (OrderID) and a lifecycle.
type Order struct {
	OrderID   string
	MenuItem  MenuItem // Value Object
	Quantity  int
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(orderID string, item MenuItem, quantity int) (*Order, error) {
	if orderID == "" || quantity <= 0 {
		return nil, errors.New("order ID and quantity must be valid")
	}
	return &Order{
		OrderID:   orderID,
		MenuItem:  item,
		Quantity:  quantity,
		Status:    OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (o *Order) MarkAsPreparing() error {
	if o.Status!= OrderStatusPending {
		return errors.New("order must be pending to be marked as preparing")
	}
	o.Status = OrderStatusPreparing
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) MarkAsReadyForDelivery() error {
	if o.Status!= OrderStatusPreparing {
		return errors.New("order must be preparing to be marked as ready for delivery")
	}
	o.Status = OrderStatusReadyForDelivery
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) MarkAsOutForDelivery() error {
	if o.Status!= OrderStatusReadyForDelivery {
		return errors.New("order must be ready for delivery to be marked as out for delivery")
	}
	o.Status = OrderStatusOutForDelivery
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) MarkAsDelivered() error {
	if o.Status!= OrderStatusOutForDelivery {
		return errors.New("order must be out for delivery to be marked as delivered")
	}
	o.Status = OrderStatusDelivered
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) MarkAsInvalid(reason string) {
	o.Status = OrderStatusInvalid
	o.UpdatedAt = time.Now()
	// In a real system, the reason might be stored or an event published.
}