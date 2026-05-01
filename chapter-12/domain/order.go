package domain

// Order represents the core domain entity for an order.
// It contains attributes directly relevant to the business concept of an order.
type Order struct {
	ID        string // This will serve as the correlation ID (UUID)
	ProductID string
	Quantity  int
	Status    string // e.g., "Pending", "Processed", "Failed"
}