package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
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

	// Variable to keep track of user orders
	userOrders := make(map[string]string) // Key: Order ID, Value: Item Name

	fmt.Println("Welcome to the Restaurant CLI!")
	fmt.Println("Type 'menu' to view the menu, 'order <item>' to order food or 'exit' to quit. Have a good evening")


	// Subscribe to restaurant updates
	go func() {
		sub, _ := nc.SubscribeSync("restaurant.updates")
		for {
			msg, err := sub.NextMsg(0) // Blocks until a message is received
			if err == nil {
				fmt.Println("[Update] " + string(msg.Data))
			}
		}
	}()

	// Subscribe to order status updates (for all orders)
	go func() {
		sub, _ := nc.SubscribeSync("order.status")
		for {
			msg, err := sub.NextMsg(0) // Blocks until a message is received
			if err == nil {
				var status OrderStatus
				if err := json.Unmarshal(msg.Data, &status); err != nil {
					fmt.Println("Error unmarshalling order status:", err)
					continue
				}

				// Show update for all orders
				if item, exists := userOrders[status.OrderID]; exists {
					if status.Status == "Delivered" {
						fmt.Printf("[Your Order %s] %s: Your order for %s has been delivered!\n", status.OrderID, status.Status, item)
					} else {
						fmt.Printf("[Your Order %s] %s: %s\n", status.OrderID, status.Status, item)
					}
				} else {
					// Show updates for other orders
					fmt.Printf("[Order %s] %s\n", status.OrderID, status.Status)
				}
			}
		}
	}()

	// Command-line input loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		command := scanner.Text()

		// Handle "exit"
		if command == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Handle "menu"
		if command == "menu" {
			msg, err := nc.Request("restaurant.menu", nil, nats.DefaultTimeout)
			if err != nil {
				fmt.Println("Error fetching menu:", err)
				continue
			}
			var menu map[string]int
			if err := json.Unmarshal(msg.Data, &menu); err != nil {
				fmt.Println("Error parsing menu:", err)
				continue
			}
			fmt.Println("Menu:")
			for item, price := range menu {
				fmt.Printf("  %s: $%d\n", item, price)
			}
			continue
		}

		// Handle "order <item>"
		if strings.HasPrefix(command, "order ") {
			item := strings.TrimSpace(strings.TrimPrefix(command, "order "))
			if item == "" {
				fmt.Println("Please specify an item to order.")
				continue
			}

			// Generate a new Order ID
			orderID := uuid.New().String()
			order := Order{
				OrderID: orderID,
				Item:    item,
			}

			// Publish the order to the order processor for processing
			orderData, _ := json.Marshal(order)
			nc.Publish("restaurant.orders", orderData)
			userOrders[orderID] = item // Keep track of the order ID and item
			fmt.Printf("Order placed for: %s. Your Order ID is %s\n", strings.Title(item), orderID)

			// No need to simulate sending order status updates here; this will be done by the Order Processor.

			continue
		}

		// Unknown command
		fmt.Println("Unknown command. Try 'menu', 'order <item>', 'track <order_id>', or 'exit'.")
	}
}

