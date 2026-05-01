package persistence

import (
	"errors"
	"chapter_8/chapter-8/domain" // Corrected import path based on your other files
	"sync"
)

// InMemoryOrderRepository implements the application.OrderRepository interface.
type InMemoryOrderRepository struct {
	orders map[string]*domain.Order
	mu     sync.RWMutex // Mutex for concurrent access
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	} // Corrected brace position
}

// Save adds or updates an order in the repository.
func (r *InMemoryOrderRepository) Save(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure the order has an ID to be used as the map key
	// Using .ID as per domain/order.go
	if order.OrderID == "" {
		return errors.New("order ID cannot be empty")
	}

	r.orders[order.OrderID] = order // Correctly assign the order to the map using its ID as the key
	return nil
}

// FindByID retrieves an order by its ID.
func (r *InMemoryOrderRepository) FindByID(orderID string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[orderID] // Correctly access the map using the orderID key
	if !exists {
		return nil, errors.New("order not found")
	}

	// Return a deep copy to prevent external modification of the stored object,
	// reinforcing domain purity (though a simple struct copy suffices for this example).
	// This creates a *new* struct with copied values.
	copyOrder := *order
	return &copyOrder, nil
}