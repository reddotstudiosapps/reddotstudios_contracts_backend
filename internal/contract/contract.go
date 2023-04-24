package contract

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
