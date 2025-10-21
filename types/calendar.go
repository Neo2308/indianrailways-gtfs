package types

type Calendar struct {
	ServiceId string `json:"service_id" csv:"service_id"`
	Monday    int    `json:"monday" csv:"monday"`
	Tuesday   int    `json:"tuesday" csv:"tuesday"`
	Wednesday int    `json:"wednesday" csv:"wednesday"`
	Thursday  int    `json:"thursday" csv:"thursday"`
	Friday    int    `json:"friday" csv:"friday"`
	Saturday  int    `json:"saturday" csv:"saturday"`
	Sunday    int    `json:"sunday" csv:"sunday"`
	StartDate Date   `json:"start_date" csv:"start_date"`
	EndDate   Date   `json:"end_date" csv:"end_date"`
}
