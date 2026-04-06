package dto

// DeliveryHistoryEntryDTO represents one transaction row in a client's delivery history.
type DeliveryHistoryEntryDTO struct {
	ClientID     int    `json:"client_id"`
	ClientName   string `json:"client_name"`
	ProductID    int    `json:"product_id"`
	ProductName  string `json:"product_name"`
	ChargedUnits int    `json:"charged_units"`
	CreditUnits  int    `json:"credit_units"`
	DeliveryDate string `json:"delivery_date"`
}

// ZoneReportEntryDTO represents one transaction row in a zone delivery report.
type ZoneReportEntryDTO struct {
	ClientID     int    `json:"client_id"`
	ClientName   string `json:"client_name"`
	ProductID    int    `json:"product_id"`
	ProductName  string `json:"product_name"`
	ChargedUnits int    `json:"charged_units"`
	CreditUnits  int    `json:"credit_units"`
	DeliveryDate string `json:"delivery_date"`
}
