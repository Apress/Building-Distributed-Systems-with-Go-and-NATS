package application

import (
	"chapter_8/chapter-8/domain"
	"errors"
	"fmt"
)

// Interfaces (Ports) for external dependencies, implemented by Infrastructure.
// These define WHAT the application needs, not HOW it gets it.
type OrderRepository interface {
	Save(order *domain.Order) error
	FindByID(orderID string) (*domain.Order, error)
	// Potentially other methods like FindByStatus, UpdateStatus, etc.
}

type OrderStatusPublisher interface {
	PublishOrderStatus(orderID string, status domain.OrderStatus, item string) error
}

// OrderProcessingService is an Application Service that orchestrates use cases.
type OrderProcessingService struct {
	orderRepo       OrderRepository
	orderFactory    *domain.OrderFactory
	menuService     *domain.MenuService
	statusPublisher OrderStatusPublisher
}

func NewOrderProcessingService(
	repo OrderRepository,
	factory *domain.OrderFactory,
	menuSvc *domain.MenuService,
	publisher OrderStatusPublisher,
) *OrderProcessingService {
	return &OrderProcessingService{
		orderRepo:       repo,
		orderFactory:    factory,
		menuService:     menuSvc,
		statusPublisher: publisher,
	}
}

// ProcessNewOrder is a specific use case method.
func (s *OrderProcessingService) ProcessNewOrder(itemName string, quantity int) (string, error) {
	// Application logic: First, check if the item is valid using the Domain's MenuService.
	if!s.menuService.IsItemValid(itemName) {
		// Publish an invalid status update immediately, even before order creation.
		s.statusPublisher.PublishOrderStatus("N/A", domain.OrderStatusInvalid, fmt.Sprintf("Item '%s' not on menu", itemName))
		return "", errors.New("item not on menu")
	}
	// Application logic: Use the Domain's OrderFactory to create a new Order entity.
	order, err := s.orderFactory.CreateOrder(itemName, quantity)
	if err!= nil {
		s.statusPublisher.PublishOrderStatus("N/A", domain.OrderStatusInvalid, fmt.Sprintf("Failed to create order for '%s': %v", itemName, err))
		return "", fmt.Errorf("failed to create order: %w", err)
	}
	// Application logic: Persist the new order using the OrderRepository (Infrastructure via port).
	if err := s.orderRepo.Save(order); err!= nil {
		s.statusPublisher.PublishOrderStatus(order.OrderID, domain.OrderStatusInvalid, fmt.Sprintf("Failed to save order '%s': %v", order.OrderID, err))
		return "", fmt.Errorf("failed to save order: %w", err)
	}
	// Application logic: Publish the initial status update (Infrastructure via port).
	s.statusPublisher.PublishOrderStatus(order.OrderID, order.Status, order.MenuItem.Name)
	return order.OrderID, nil
}

// UpdateOrderStatus is another use case method, potentially triggered by an internal event (e.g., from Delivery Tracker).
func (s *OrderProcessingService) UpdateOrderStatus(orderID string, newStatus domain.OrderStatus) error {
	// Application logic: Retrieve the order from persistence.
	order, err := s.orderRepo.FindByID(orderID)
	if err!= nil {
		return fmt.Errorf("order %s not found: %w", orderID, err)
	}

	// Check if the order is already in the requested status.
	// If so, no state transition is needed, and we can return early.
	if order.Status == newStatus {
		fmt.Printf("Application Layer: Order %s is already in status '%s'. No transition needed.\n", orderID, newStatus)
		return nil
	}

	// Delegate status change to the domain entity (business logic).
	var statusErr error
	switch newStatus {
	case domain.OrderStatusPreparing:
		statusErr = order.MarkAsPreparing()
	case domain.OrderStatusReadyForDelivery:
		statusErr = order.MarkAsReadyForDelivery()
	case domain.OrderStatusOutForDelivery:
		statusErr = order.MarkAsOutForDelivery()
	case domain.OrderStatusDelivered:
		statusErr = order.MarkAsDelivered()
	case domain.OrderStatusInvalid:
		order.MarkAsInvalid("External invalidation") // Direct call if no specific transition rules
	default:
		return errors.New("invalid status transition")
	}
	if statusErr!= nil {
		return fmt.Errorf("failed to update order status: %w", statusErr)
	}
	// Application logic: Save the updated order.
	if err := s.orderRepo.Save(order); err!= nil {
		return fmt.Errorf("failed to save updated order: %w", err)
	}
	// Application logic: Publish the updated status.
	s.statusPublisher.PublishOrderStatus(order.OrderID, order.Status, order.MenuItem.Name)
	return nil
}