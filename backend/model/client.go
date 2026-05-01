package model

import "time"

// Client represents a delivery customer assigned to a zone.
type Client struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Address        string     `json:"address"`
	ZoneID         int        `json:"zone_id"`
	CreateDate     time.Time  `json:"create_date"`
	UnregisterDate *time.Time `json:"unregister_date"`
	PriorityOrder  int        `json:"priority_order"`
}
