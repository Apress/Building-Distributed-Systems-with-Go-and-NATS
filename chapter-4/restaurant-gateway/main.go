package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// Function to generate dynamic menu updates
func generateUpdate() string {
	updates := []string{
		"New item added: Burger for $8!",
		"Half price on Pizza today only!",
		"Buy one Sushi, get one free!",
		"Special offer: Pasta for $7!",
		"Limited time: Ice Cream for $3!",
	}
	return updates[rand.Intn(len(updates))]
}

func main() {
	// Check chapter 4 Minikube example for this
	nu := os.Getenv("NATS_URL")
	if nu == "" {
		nu = nats.DefaultURL
	}
	fmt.Println("Restaurant gateway starting")
	// Connect to NATS - check Chapter 4 Minikube example for any discrepancy
	nc, err := nats.Connect(nu)
	if err != nil {
		fmt.Println("Error connecting to NATS:", err)
		return
	}
	defer nc.Close()

	// Menu with more items
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

	// Request/Reply: Handle customer queries
	nc.Subscribe("restaurant.menu", func(msg *nats.Msg) {
		// Construct a JSON representation of the menu
		menuResponse := `{`
		for item, price := range menu {
			menuResponse += fmt.Sprintf(`"%s": %d,`, item, price)
		}
		menuResponse = menuResponse[:len(menuResponse)-1] + `}` // Remove last comma and close JSON
		fmt.Println(menuResponse)
		msg.Respond([]byte(menuResponse))
		fmt.Println("Sent menu to customer.")
	})

	// Publish/Subscribe: Publish menu updates periodically
	go func() {
		for {
			time.Sleep(10 * time.Second) // Publish updates every 10 seconds
			update := generateUpdate()
			nc.Publish("restaurant.updates", []byte(update))
			fmt.Println("Published menu update:", update)
		}
	}()

	// Keep the connection alive
	select {}
}

