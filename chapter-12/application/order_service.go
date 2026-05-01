// application/order_service.go
package application

import (
	"errors"
	"log"

	"chapter-12/domain" // Import the domain package
)

// OrderServicePort is a "driving port" interface.
// It defines what the presentation layer (or any external client) needs from the application layer.
type OrderServicePort interface {
	HandleOrder(order domain.Order) (domain.Order, error)
}

// OrderProcessor is a "driven port" interface.
// It defines what the application layer needs from the infrastructure layer
// to process an order (e.g., interact with external systems, persist data).
type OrderProcessor interface {
	ProcessOrder(order domain.Order) (domain.Order, error)
}

// OrderApplicationService implements the application's use case for order processing.
// It depends on the OrderProcessor (driven) port.
// It also implicitly implements OrderServicePort because it has the HandleOrder method.
type OrderApplicationService struct {
	processor OrderProcessor // This is the 'driven port'
}

// NewOrderApplicationService creates an instance of the application service.
func NewOrderApplicationService(processor OrderProcessor) *OrderApplicationService {
	return &OrderApplicationService{processor: processor}
}

// HandleOrder is the application's use case function.
// It performs business validation and delegates the actual processing to the infrastructure layer.
func (s *OrderApplicationService) HandleOrder(order domain.Order) (domain.Order, error) {
	// --- Business Validation ---
	if order.ID == "" {
		return domain.Order{}, errors.New("order ID is required")
	}
	if order.ProductID == "" {
		return domain.Order{}, errors.New("product ID is required")
	}
	if order.Quantity <= 0 {
		return domain.Order{}, errors.New("quantity must be positive")
	}

	log.Printf("Application Service: Handling order %s for Product %s (Qty: %d)\n", order.ID, order.ProductID, order.Quantity)

	// --- Delegate to the infrastructure layer (via the port) for actual processing ---
	processedOrder, err := s.processor.ProcessOrder(order)
	if err!= nil {
		log.Printf("Application Service: Error processing order %s: %v\n", order.ID, err)
		return processedOrder, err
	}

	log.Printf("Application Service: Order %s processed successfully.\n", order.ID)
	return processedOrder, nil
}