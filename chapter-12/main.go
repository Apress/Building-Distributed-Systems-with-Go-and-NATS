// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"chapter-12/application"
	"chapter-12/infrastructure"
	"chapter-12/presentation"

	"github.com/nats-io/nats.go"
)

func main() {
	// --- 1. Connect to NATS ---
	nc, err := nats.Connect(nats.DefaultURL)
	if err!= nil {
		log.Fatalf("Order Service: Failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	defer nc.Drain() // Ensure graceful shutdown of NATS connection

	log.Println("Order Service connected to NATS server.")

	// --- 2. Wire up DDD Layers ---

	// Infrastructure Layer Adapter (Driven Adapter): Implements the OrderProcessor port
	// This adapter handles actual NATS communication for processing and replying.
	natsOrderAdapter := infrastructure.NewNATSOrderAdapter(nc)

	// Application Layer Service: Uses the OrderProcessor port (implemented by natsOrderAdapter)
	// This service contains the core business logic and use cases.
	orderAppService := application.NewOrderApplicationService(natsOrderAdapter)

	// Presentation Layer Handler (Driving Adapter): Handles incoming NATS messages.
	// It uses the application service and the NATS adapter to send replies.
	natsOrderHandler := presentation.NewNATSOrderHandler(orderAppService, natsOrderAdapter)

	// --- 3. Configure and Start NATS Queue Group Listener ---

	// Define the queue group for this service. Multiple instances will join this group.
	queueGroup := "order_processors_group"
	instanceID := os.Getenv("INSTANCE_ID") // Allow setting instance ID for clearer logs
	if instanceID == "" {
		instanceID = "OrderProcessor-Default"
	}
	log.Printf("Order Service Instance %s joining NATS queue group '%s'\n", instanceID, queueGroup)

	// Subscribe to incoming order requests using a Queue Group.
	// NATS will distribute messages among instances in this group. [1]
	err = natsOrderAdapter.SubscribeToOrderRequests(queueGroup, natsOrderHandler.HandleNATSOrder)
	if err!= nil {
		log.Fatalf("Order Service: Failed to subscribe to order requests: %v", err)
	}

	log.Println("Order Service running and listening for messages...")

	// --- Keep the service running until interrupted ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Block until a signal is received
	log.Println("Order Service shutting down.")
}