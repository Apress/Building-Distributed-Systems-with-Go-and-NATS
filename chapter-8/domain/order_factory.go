package domain

import (
	"errors"
	"github.com/google/uuid"
)

type OrderFactory struct {
	menuService *MenuService
}

func NewOrderFactory(ms *MenuService) *OrderFactory {
	return &OrderFactory{menuService: ms}
}

func (f *OrderFactory) CreateOrder(itemName string, quantity int) (*Order, error) {
	menuItem, exists := f.menuService.GetMenuItem(itemName)
	if!exists {
		return nil, errors.New("item not found on menu")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	orderID := uuid.New().String() // Order Processor generates a unique ID
	return NewOrder(orderID, menuItem, quantity)
}