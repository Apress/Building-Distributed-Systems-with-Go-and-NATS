package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"chapter_8/chapter-8/application"
	"chapter_8/chapter-8/domain"
	"chapter_8/chapter-8/infrastructure/messaging"
	"chapter_8/chapter-8/infrastructure/persistence"
	"chapter_8/chapter-8/presentation"

	"github.com/nats-io/nats.go"
)

// Main application entry point
func main() {
	fmt.Println("Starting Restaurant Order Processing Application...")

	// 1. Initialize NATS Connection
	// 1. Initialize NATS Connection
	// Check for NATS_URL environment variable, otherwise use nats.DefaultURL
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL // Default to "nats://127.0.0.1:4222" if env var is not set
		fmt.Printf("NATS_URL environment variable not set. Using default NATS URL: %s\n", natsURL)
	} else {
		fmt.Printf("Using NATS URL from environment variable: %s\n", natsURL)
	}

	nc, err := nats.Connect(natsURL, nats.Name("OrderProcessorService"))
	if err != nil {
		fmt.Printf("Error connecting to NATS: %v\n", err)
		os.Exit(1)
	}
	defer nc.Close()
	fmt.Printf("Connected to NATS server: %s\n", natsURL)

	// 2. Setup Domain Layer Components
	// Define initial menu items
	pizza, _ := domain.NewMenuItem("Pizza", 12.50)
	burger, _ := domain.NewMenuItem("Burger", 8.00)
	salad, _ := domain.NewMenuItem("Salad", 7.00)
	drinks, _ := domain.NewMenuItem("Coke", 2.00)

	initialMenuItems := []domain.MenuItem{pizza, burger, salad, drinks}

	// MenuService manages the menu items within the domain
	menuService := domain.NewMenuService(initialMenuItems)
	// OrderFactory is responsible for creating new Order entities
	orderFactory := domain.NewOrderFactory(menuService)

	// 3. Setup Infrastructure Layer Components
	// In-memory repository for orders (simulating a database)
	orderRepo := persistence.NewInMemoryOrderRepository()
	// NATS publisher for order status updates
	statusPublisher := messaging.NewNATSOrderStatusPublisher(nc)

	// 4. Setup Application Layer
	// OrderProcessingService orchestrates the core business logic/use cases
	orderProcessingSvc := application.NewOrderProcessingService(
		orderRepo,
		orderFactory,
		menuService,
		statusPublisher,
	)

	// 5. Setup Presentation Layer (NATS Receivers)
	// NATSOrderReceiver listens for incoming new order requests and status updates
	natsReceiver := presentation.NewNATSOrderReceiver(nc, orderProcessingSvc)

	// Start listening for new order requests
	if err := natsReceiver.Start(); err != nil {
		fmt.Printf("Error starting NATS order receiver: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Subscribed to 'restaurant.orders' for new order requests.")

	// Start listening for order status updates (e.g., from a simulated delivery service)
	if err := natsReceiver.StartOrderStatusUpdater(); err != nil {
		fmt.Printf("Error starting NATS order status updater: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Subscribed to 'order.status' for internal status updates.")

	fmt.Println("\nApplication is running. Send NATS messages to 'restaurant.orders' to create orders.")
	fmt.Println("Or send NATS messages to 'order.status' to update existing order statuses.")
	fmt.Println("Type 'exit' or 'quit' to stop the application.")
	fmt.Println("---")

	// Start a simple NATS publisher client for testing
	go startTestPublisher(nc)

	// Keep the main goroutine alive
	select {}
}

// startTestPublisher provides a CLI to send test NATS messages
func startTestPublisher(nc *nats.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type 'new <item> <quantity>' to send a new order request (e.g., new Pizza 2)")
	fmt.Println("Type 'update <orderID> <status>' to send a status update (e.g., update 12345 Preparing)")
	fmt.Println("Available statuses: Pending, Preparing, Ready for Delivery, Out For Delivery, Delivered, Invalid")
	fmt.Println("---")

	for {
		fmt.Print("Enter command: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting test publisher.")
			return
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "new":
			if len(parts) == 3 {
				item := parts[1]
				quantityStr := parts[2]
				var quantity int
				_, err := fmt.Sscanf(quantityStr, "%d", &quantity)
				if err != nil {
					fmt.Printf("Invalid quantity: %v\n", err)
					continue
				}

				dto := presentation.OrderRequestDTO{Item: item, Quantity: quantity}
				data, err := json.Marshal(dto)
				if err != nil {
					fmt.Printf("Error marshalling new order DTO: %v\n", err)
					continue
				}

				err = nc.Publish("restaurant.orders", data)
				if err != nil {
					fmt.Printf("Error publishing new order: %v\n", err)
				} else {
					fmt.Printf("Published new order request: %s\n", string(data))
				}
			} else {
				fmt.Println("Usage: new <item> <quantity>")
			}
		case "update":
			if len(parts) >= 3 {
				orderID := parts[1]
				statusStr := strings.Join(parts[2:], " ")

				// Basic validation for status string
				validStatus := false
				for _, s := range []domain.OrderStatus{
					domain.OrderStatusPending, domain.OrderStatusPreparing,
					domain.OrderStatusReadyForDelivery, domain.OrderStatusOutForDelivery,
					domain.OrderStatusDelivered, domain.OrderStatusInvalid,
				} {
					if strings.EqualFold(string(s), statusStr) {
						validStatus = true
						break
					}
				}
				if !validStatus {
					fmt.Printf("Invalid status: %s. Please use one of: Pending, Preparing, Ready for Delivery, Out For Delivery, Delivered, Invalid\n", statusStr)
					continue
				}

				dto := presentation.OrderStatusUpdateDTO{
					OrderID: orderID,
					Status:  statusStr,
					// Item field is optional for status updates from external systems, but often good for context.
					// For this simple client, we'll leave it empty.
					Item: "",
				}
				data, err := json.Marshal(dto)
				if err != nil {
					fmt.Printf("Error marshalling status update DTO: %v\n", err)
					continue
				}

				err = nc.Publish("order.status", data)
				if err != nil {
					fmt.Printf("Error publishing status update: %v\n", err)
				} else {
					fmt.Printf("Published status update for order %s: %s\n", orderID, string(data))
				}
			} else {
				fmt.Println("Usage: update <orderID> <status>")
			}
		default:
			fmt.Println("Unknown command. Use 'new' or 'update'.")
		}
		time.Sleep(100 * time.Millisecond) // Small delay to avoid busy loop
	}
}