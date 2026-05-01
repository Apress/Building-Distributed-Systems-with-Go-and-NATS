package presentation

import (
	"chapter_8/chapter-8/application" // Depends on application layer
	"chapter_8/chapter-8/domain"      // Depends on domain layer for OrderStatus type
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
)

// OrderRequestDTO represents the data structure for an incoming new order request.
// It mirrors the expected payload from external clients (e.g., the customer application).
type OrderRequestDTO struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	// OrderID could optionally be included here if the client provides a correlation ID,
	// but the Order Processor will generate its authoritative OrderID.
	// OrderID string `json:"order_id,omitempty"`
}

// NATSOrderReceiver is an Inbound Adapter for receiving new order requests via NATS.
type NATSOrderReceiver struct {
	nc                  *nats.Conn
	orderProcessingSvc *application.OrderProcessingService
}

func NewNATSOrderReceiver(nc *nats.Conn, ops *application.OrderProcessingService) *NATSOrderReceiver {
	return &NATSOrderReceiver{
		nc:                  nc,
		orderProcessingSvc: ops,
	}
}

// Start subscribes to the NATS subject for new orders (Pub/Sub style).
// It receives order requests, unmarshals them, and delegates to the application layer.
func (r *NATSOrderReceiver) Start() error {
	_, err := r.nc.Subscribe("restaurant.orders", func(msg *nats.Msg) {
		var orderReq OrderRequestDTO
		if err := json.Unmarshal(msg.Data, &orderReq); err!= nil {
			fmt.Printf("Presentation Layer: Error unmarshalling order request: %v\n", err)
			// For Pub/Sub, typically just log and return. No direct response to the sender.
			return
		}

		fmt.Printf("Presentation Layer: Received new order request for item: '%s', quantity: %d\n", orderReq.Item, orderReq.Quantity)

		// Delegate to the Application Layer for processing.
		// The application service handles domain logic and publishing status updates.
		orderID, err := r.orderProcessingSvc.ProcessNewOrder(orderReq.Item, orderReq.Quantity)
		if err!= nil {
			fmt.Printf("Presentation Layer: Failed to process order for '%s': %v\n", orderReq.Item, err)
			// The application service already publishes an 'Invalid' status if creation/save fails.
		} else {
			fmt.Printf("Presentation Layer: Order for '%s' (ID: %s) successfully initiated.\n", orderReq.Item, orderID)
		}
	})
	if err!= nil {
		return fmt.Errorf("failed to subscribe to 'restaurant.orders' subject: %w", err)
	}
	return nil
}

// OrderStatusUpdateDTO represents the data structure for an incoming order status update.
// This is typically received from other internal services (e.g., Delivery Tracker).
type OrderStatusUpdateDTO struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Item    string `json:"item"` // For context, though not strictly needed for status update logic
}

// StartOrderStatusUpdater subscribes to NATS for status updates (e.g., from Delivery Tracker)
// and delegates to the application layer to update the order status.
// This acts as an inbound adapter for internal service communication.
func (r *NATSOrderReceiver) StartOrderStatusUpdater() error {
	_, err := r.nc.Subscribe("order.status", func(msg *nats.Msg) {
		var statusUpdate OrderStatusUpdateDTO
		if err := json.Unmarshal(msg.Data, &statusUpdate); err!= nil {
			fmt.Printf("Presentation Layer: Error unmarshalling status update: %v\n", err)
			return
		}

		fmt.Printf("Presentation Layer: Received status update for Order ID '%s': %s\n", statusUpdate.OrderID, statusUpdate.Status)

		// Convert DTO status string to domain.OrderStatus type
		domainStatus := domain.OrderStatus(statusUpdate.Status)

		// Delegate to the Application Layer to update the order status.
		err := r.orderProcessingSvc.UpdateOrderStatus(statusUpdate.OrderID, domainStatus)
		if err!= nil {
			fmt.Printf("Presentation Layer: Failed to update status for Order ID '%s': %v\n", statusUpdate.OrderID, err)
		} else {
			fmt.Printf("Presentation Layer: Order ID '%s' status updated to '%s' successfully.\n", statusUpdate.OrderID, domainStatus)
		}
	})
	if err!= nil {
		return fmt.Errorf("failed to subscribe to 'order.status' subject: %w", err)
	}
	return nil
}