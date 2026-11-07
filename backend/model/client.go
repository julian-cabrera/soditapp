package model

// Client represents a delivery customer assigned to a zone.
type Client struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Address        string  `json:"address"`
	ZoneID         string  `json:"zone_id"`
	CreateDate     string  `json:"create_date"`
	UnregisterDate *string `json:"unregister_date"`
	PriorityOrder  int     `json:"priority_order"`
	ReferenceUUID  string  `json:"reference_uuid"`
}
