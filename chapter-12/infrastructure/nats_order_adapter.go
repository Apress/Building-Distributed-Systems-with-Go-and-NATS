// infrastructure/nats_order_adapter.go
package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"chapter-12/domain" // Import the domain package
	"github.com/nats-io/nats.go"
)

// NATSOrderRequest represents the message format received from the API Gateway over NATS.
// It includes the OrderID as the correlation ID and the ReplySubject for the response.
type NATSOrderRequest struct {
	OrderID      string `json:"order_id"` // This is the correlation ID
	ProductID    string `json:"product_id"`
	Quantity     int    `json:"quantity"`
	ReplySubject string `json:"reply_subject"` // The subject the API Gateway is listening on for this specific request.
}

// NATSOrderResponse represents the message format sent back to the API Gateway over NATS.
type NATSOrderResponse struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"` // "Processed" or "Failed"
	Message   string `json:"message"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// NATSOrderAdapter is a "driven adapter" for NATS communication.
// It implements the application.OrderProcessor interface.
type NATSOrderAdapter struct {
	natsConn *nats.Conn
}

// NewNATSOrderAdapter creates a new NATS adapter for order processing.
func NewNATSOrderAdapter(nc *nats.Conn) *NATSOrderAdapter {
	return &NATSOrderAdapter{natsConn: nc}
}

// ProcessOrder simulates the actual processing of an order.
// This method is the concrete implementation of the application.OrderProcessor port.
func (a *NATSOrderAdapter) ProcessOrder(order domain.Order) (domain.Order, error) {
	log.Printf("Infrastructure Adapter: Simulating processing for order ID: %s, Product: %s, Qty: %d\n", order.ID, order.ProductID, order.Quantity)
	time.Sleep(2 * time.Second) // Simulate long-running work

	// --- Simulate Business Logic Outcome ---
	// Example: Orders with quantity > 100 are considered failed.
	if order.Quantity > 100 {
		order.Status = "Failed"
		log.Printf("Infrastructure Adapter: Order %s failed due to large quantity (%d).\n", order.ID, order.Quantity)
		return order, fmt.Errorf("order quantity too large: %d", order.Quantity)
	}

	order.Status = "Processed"
	log.Printf("Infrastructure Adapter: Order %s successfully processed.\n", order.ID)
	return order, nil
}

// SubscribeToOrderRequests sets up the NATS Queue Group subscription for incoming order requests.
// This is where the service listens for messages from the API Gateway.
// The handler function (from the Presentation layer) will be called when a message arrives.
func (a *NATSOrderAdapter) SubscribeToOrderRequests(queueGroup string, handler func(msg *nats.Msg)) error {
	subject := "orders.process" // The subject for incoming order requests
	_, err := a.natsConn.QueueSubscribe(subject, queueGroup, handler) // [1]
	if err!= nil {
		return fmt.Errorf("failed to subscribe to %s in queue group %s: %w", subject, queueGroup, err)
	}
	log.Printf("Infrastructure Adapter: Subscribed to '%s' in queue group '%s'\n", subject, queueGroup)
	return nil
}

// SendOrderReply sends the processing result back to the API Gateway.
// This method is called by the Presentation layer after the application logic is complete.
func (a *NATSOrderAdapter) SendOrderReply(replySubject string, order domain.Order, appErr error) {
	var natsResp NATSOrderResponse
	if appErr!= nil {
		natsResp = NATSOrderResponse{
			OrderID:   order.ID,
			Status:    "Failed",
			Message:   fmt.Sprintf("Processing failed: %v", appErr),
			ProductID: order.ProductID,
			Quantity:  order.Quantity,
		}
	} else {
		natsResp = NATSOrderResponse{
			OrderID:   order.ID,
			Status:    "Processed",
			Message:   "Order processed successfully.",
			ProductID: order.ProductID,
			Quantity:  order.Quantity,
		}
	}

	respBytes, marshalErr := json.Marshal(natsResp)
	if marshalErr!= nil {
		log.Printf("Infrastructure Adapter: Error marshalling response for order %s: %v\n", order.ID, marshalErr)
		return
	}

	// Publish the response to the specific reply subject provided by the API Gateway
	if publishErr := a.natsConn.Publish(replySubject, respBytes); publishErr!= nil { // [1]
		log.Printf("Infrastructure Adapter: Error publishing reply for order %s to '%s': %v\n", order.ID, replySubject, publishErr)
	} else {
		log.Printf("Infrastructure Adapter: Published reply for order %s to '%s'\n", order.ID, replySubject)
	}
}