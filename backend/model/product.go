package model

// Product represents a deliverable product.
type Product struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Active        bool   `json:"active"`
	ReferenceUUID string `json:"reference_uuid"`
}
