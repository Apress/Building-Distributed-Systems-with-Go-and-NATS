package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"os"
)

// Order structure to hold the order ID and item
type Order struct {
	OrderID string `json:"order_id"`
	Item    string `json:"item"`
}

// Order Status structure to simulate status updates
type OrderStatus struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	fmt.Println("Order processor starting")
	// Define the menu with items and prices
	menu := map[string]int{
		"pizza":       10,
		"sushi":       15,
		"burger":      8,
		"pasta":       12,
		"ice cream":   5,
		"fries":       4,
		"milkshake":   6,
		"hot dog":     7,
	}

	nu := os.Getenv("NATS_URL")
	if nu == "" {
		nu = nats.DefaultURL
	}

	// Connect to NATS
	nc, err := nats.Connect(nu)
	if err != nil {
		fmt.Println("Error connecting to NATS:", err)
		return
	}
	defer nc.Close()

	// Subscribe to incoming orders
	nc.Subscribe("restaurant.orders", func(msg *nats.Msg) {
		var order Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			fmt.Println("Error unmarshalling order:", err)
			return
		}

		// Check if the ordered item is in the menu
		if price, exists := menu[order.Item]; exists {
			// Valid order: Process it
			fmt.Printf("Processing order %s for item: %s (Price: $%d)\n", order.OrderID, order.Item, price)

			// Simulate order processing and status updates
			statuses := []string{
				"Preparing",
				"Ready for Delivery",
			}
			for _, status := range statuses {
				time.Sleep(5 * time.Second) // Simulate time between status updates

				// Publish status updates to the "order.status" topic
				statusUpdate := OrderStatus{
					OrderID: order.OrderID,
					Status:  status,
					}
					statusData, _ := json.Marshal(statusUpdate)
					nc.Publish("order.status", statusData)
				}
			} else {
				// Invalid order: Item not on the menu
				fmt.Printf("Order %s for item '%s' is invalid. Item not on the menu.\n", order.OrderID, order.Item)

				// Publish an invalid order status
				statusUpdate := OrderStatus{
					OrderID: order.OrderID,
					Status:  "Invalid order: Item not on the menu",
				}
				statusData, _ := json.Marshal(statusUpdate)
				nc.Publish("order.status", statusData)
			}
		})

		// Keep the connection alive
		select {}
	}
