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
	fmt.Println("Delivery tracker starting")

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

	// Subscribe to "Ready for Delivery" orders
	nc.Subscribe("order.status", func(msg *nats.Msg) {
		var statusUpdate OrderStatus
		if err := json.Unmarshal(msg.Data, &statusUpdate); err != nil {
			fmt.Println("Error unmarshalling status update:", err)
			return
		}

		// Check if the status is "Ready for Delivery"
		if statusUpdate.Status == "Ready for Delivery" {
			fmt.Printf("Order %s is ready for delivery. Starting delivery process...\n", statusUpdate.OrderID)

			// Simulate "Out for Delivery" status
			time.Sleep(3 * time.Second) // Simulate delay before assigning to delivery
			outForDelivery := OrderStatus{
				OrderID: statusUpdate.OrderID,
				Status:  "Out for delivery",
			}
			outForDeliveryData, _ := json.Marshal(outForDelivery)
			nc.Publish("order.status", outForDeliveryData)
			fmt.Printf("Order %s is now out for delivery.\n", statusUpdate.OrderID)

			// Simulate delivery process
			time.Sleep(5 * time.Second) // Simulate delivery time
			delivered := OrderStatus{
				OrderID: statusUpdate.OrderID,
				Status:  "Delivered",
			}
			deliveredData, _ := json.Marshal(delivered)
			nc.Publish("order.status", deliveredData)
			fmt.Printf("Order %s has been delivered.\n", statusUpdate.OrderID)
		}
	})

	// Keep the connection alive
	select {}
}
