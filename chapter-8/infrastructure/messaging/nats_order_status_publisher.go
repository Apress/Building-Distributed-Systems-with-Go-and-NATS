package messaging

import (
	"encoding/json"
	"fmt"
	"chapter_8/chapter-8/domain" // Depends on domain layer
	"github.com/nats-io/nats.go"
)

// NATSOrderStatusPublisher implements the application.OrderStatusPublisher interface.
type NATSOrderStatusPublisher struct {
	nc *nats.Conn
}

func NewNATSOrderStatusPublisher(nc *nats.Conn) *NATSOrderStatusPublisher {
	return &NATSOrderStatusPublisher{nc: nc}
}

func (p *NATSOrderStatusPublisher) PublishOrderStatus(orderID string, status domain.OrderStatus, itemName string) error {
	// Define a DTO for external communication, ensuring we only expose what's necessary.
	// This DTO mirrors the structure from Chapter 4's OrderStatus for consistency.
	statusUpdateDTO := struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
		Item    string `json:"item"` // Added for context in tracking by other services
	}{
		OrderID: orderID,
		Status:  string(status),
		Item:    itemName,
	}
	statusData, err := json.Marshal(statusUpdateDTO)
	if err!= nil {
		return fmt.Errorf("failed to marshal order status update DTO: %w", err)
	}
	// Create a NATS message. Use PublishMsg to allow adding headers.
	msg := &nats.Msg{
		Subject: "order.status", // The subject for order status updates
		Data:    statusData,
	}
	// Example of adding a custom header for error context.
	if status == domain.OrderStatusInvalid {
		msg.Header = nats.Header{} // Initialize header if not already
		msg.Header.Add("ErrorType", "InvalidOrder")
		msg.Header.Add("Reason", "Item Not On Menu") // More specific error detail
	}
	return p.nc.PublishMsg(msg)
}