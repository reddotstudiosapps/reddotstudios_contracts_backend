package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ClientDetails struct {
	ClientName  string `json:"clientName"`
	ClientEmail string `json:"clientEmail"`
}

type Deliverable struct {
	Description  string `json:"description"`
	Quantity     string `json:"quantity"`
	Mode         string `json:"mode"`
	DeliveryDate string `json:"deliveryDate"`
}

type EventDetails struct {
	EventName         string `json:"eventName"`
	EventDate         string `json:"eventDate"`
	EventCoverageTime string `json:"eventCoverageTime"`
	EventVenue        string `json:"eventVenue"`
}

type PaymentDetails struct {
	TotalAmount        int64  `json:"totalAmount"`
	AdvancePaid        int64  `json:"advancePaid,omitempty"`
	AdvancePaymentMode string `json:"advancePaymentMode,omitempty"`
	PerHourExtra       int64  `json:"perHourExtra"`
}

type Contract struct {
	ClientDetails      ClientDetails  `json:"clientDetails"`
	EventDetails       EventDetails   `json:"eventDetails"`
	PaymentDetails     PaymentDetails `json:"paymentDetails"`
	DeliverableDetails []Deliverable  `json:"deliverableDetails"`
}

func NewContractHandler(c *fiber.Ctx) error {
	var contract Contract
	if err := c.BodyParser(&contract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse JSON",
		})
	}

	if err := ValidateContract(contract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Contract created successfully",
	})
}

func main() {
	// Create a new Fiber instance
	app := fiber.New()

	// Register the logger middleware
	app.Use(logger.New())

	config := cors.Config{
		AllowOrigins:     "http://localhost:3000,https://rds-contracts-ui.vercel.app",
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           3600,
	}

	// Add the CORS middleware
	app.Use(cors.New(config))

	// Define a basic GET route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	// Handle a new contract
	app.Post("/newcontract", NewContractHandler)

	port := 8080
	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error starting server on port %d: %v", port, err)
	}
}

func ValidateContract(contract Contract) error {
	if contract.ClientDetails.ClientName == "" {
		return errors.New("client name is required")
	}

	if contract.ClientDetails.ClientEmail == "" {
		return errors.New("client email is required")
	} else if !strings.Contains(contract.ClientDetails.ClientEmail, "@") {
		return errors.New("client email is not valid")
	}

	if contract.EventDetails.EventName == "" {
		return errors.New("event name is required")
	}

	if contract.EventDetails.EventDate == "" {
		return errors.New("event date is required")
	}

	if contract.EventDetails.EventCoverageTime == "" {
		return errors.New("event coverage time is required")
	}

	if contract.EventDetails.EventVenue == "" {
		return errors.New("event venue is required")
	}

	if contract.PaymentDetails.TotalAmount <= 0 {
		return errors.New("total amount should be greater than zero")
	}

	if contract.PaymentDetails.AdvancePaid < 0 {
		return errors.New("advance paid cannot be negative")
	}

	if contract.PaymentDetails.PerHourExtra < 0 {
		return errors.New("per hour extra cannot be negative")
	}

	if len(contract.DeliverableDetails) == 0 {
		return errors.New("at least one deliverable is required")
	}

	for _, deliverable := range contract.DeliverableDetails {
		if deliverable.Description == "" {
			return errors.New("deliverable description is required")
		}

		if deliverable.Quantity == "" {
			return errors.New("deliverable quantity is required")
		}

		if deliverable.Mode == "" {
			return errors.New("deliverable mode is required")
		}

		if deliverable.DeliveryDate == "" {
			return errors.New("deliverable delivery date is required")
		}
	}

	return nil
}
