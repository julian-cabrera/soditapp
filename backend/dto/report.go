package dto

// DeliveryHistoryEntryDTO represents one transaction row in a client's delivery history.
type DeliveryHistoryEntryDTO struct {
	ClientRefUUID string `json:"client_ref_uuid"`
	ClientName    string `json:"client_name"`
	ProductRefUUID string `json:"product_ref_uuid"`
	ProductName   string `json:"product_name"`
	ChargedUnits  int    `json:"charged_units"`
	CreditUnits   int    `json:"credit_units"`
	DeliveryDate  string `json:"delivery_date"`
}

// ZoneReportEntryDTO represents one transaction row in a zone delivery report.
type ZoneReportEntryDTO struct {
	ClientRefUUID string `json:"client_ref_uuid"`
	ClientName    string `json:"client_name"`
	ProductRefUUID string `json:"product_ref_uuid"`
	ProductName   string `json:"product_name"`
	ChargedUnits  int    `json:"charged_units"`
	CreditUnits   int    `json:"credit_units"`
	DeliveryDate  string `json:"delivery_date"`
}
