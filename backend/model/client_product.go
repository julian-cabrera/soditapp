package model

// ClientProduct tracks the running credit balance for a client-product pair.
type ClientProduct struct {
	ID        int `json:"id"`
	ClientID  int `json:"client_id"`
	ProductID int `json:"product_id"`
	Balance   int `json:"balance"`
}
