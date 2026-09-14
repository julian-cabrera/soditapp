package model

import "time"

// DeliveryTransaction records a single delivery entry (per client, per product) for history and reporting.
type DeliveryTransaction struct {
	ID            int       `json:"id"`
	ClientID      string    `json:"client_id"`
	ProductID     string    `json:"product_id"`
	ChargedUnits  int       `json:"charged_units"`
	CreditUnits   int       `json:"credit_units"`
	DeliveryDate  time.Time `json:"delivery_date"`
}
