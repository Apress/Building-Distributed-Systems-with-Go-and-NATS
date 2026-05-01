// presentation/nats_handler.go
package presentation

import (
	"encoding/json"
	"fmt"
	"log"

	"chapter-12/application"     // Import application for OrderServicePort
	"chapter-12/domain"         // Import domain for Order type
	"chapter-12/infrastructure" // Import infrastructure for NATSOrderRequest and NATSOrderAdapter
	"github.com/nats-io/nats.go"
)

// NATSOrderHandler is a "driving adapter" that handles incoming NATS messages.
// It orchestrates the flow from NATS message reception to application logic and back to NATS reply.
type NATSOrderHandler struct {
	appService  application.OrderServicePort    // Now correctly typed as the interface from the application layer
	natsAdapter *infrastructure.NATSOrderAdapter // Adapter to send the reply back via NATS
}

// NATSOrderHandler is a "driving adapter" that handles incoming NATS messages.
// It orchestrates the flow from NATS message reception to application logic and back to NATS reply.

func NewNATSOrderHandler(appSvc application.OrderServicePort, natsAdpt *infrastructure.NATSOrderAdapter) *NATSOrderHandler {
	return &NATSOrderHandler{appService: appSvc, natsAdapter: natsAdpt}
}

// HandleNATSOrder is the NATS callback function that processes incoming order messages.
// This function will be passed to nc.QueueSubscribe in main.
func (h *NATSOrderHandler) HandleNATSOrder(m *nats.Msg) {
	log.Printf("Presentation Layer: Received NATS message on subject '%s' from NATS Queue Group.\n", m.Subject)

	var req infrastructure.NATSOrderRequest
	if err := json.Unmarshal(m.Data, &req); err!= nil {
		log.Printf("Presentation Layer: Invalid NATS request payload: %v\n", err)
		// Attempt to send a generic error reply if reply subject is available
		if req.ReplySubject!= "" {
			h.natsAdapter.SendOrderReply(req.ReplySubject, domain.Order{ID: "unknown"}, fmt.Errorf("invalid request format"))
		}
		return
	}

	// Map NATS request DTO to domain entity
	order := domain.Order{
		ID:        req.OrderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    "Pending", // Initial status before processing
	}

	// Call the application service (through the driving port) to handle the order
	processedOrder, appErr := h.appService.HandleOrder(order)

	// Send the reply back to the API Gateway using the infrastructure adapter
	// The replySubject is carried from the incoming NATS request.
	h.natsAdapter.SendOrderReply(req.ReplySubject, processedOrder, appErr)
}