package hdb

// Static carpark information from data.gov.sg
type CarparkInfoResponse struct {
	CarParkNo           string  `json:"car_park_no"`
	Address             string  `json:"address"`
	XCoord              string  `json:"x_coord"`
	YCoord              string  `json:"y_coord"`
	CarParkType         string  `json:"car_park_type"`
	TypeOfParkingSystem string  `json:"type_of_parking_system"`
	ShortTermParking    string  `json:"short_term_parking"`
	FreeParking         string  `json:"free_parking"`
	NightParking        string  `json:"night_parking"`
	CarParkDecks        string `json:"car_park_decks"`
	GantryHeight        string `json:"gantry_height"`
	CarParkBasement     string  `json:"car_park_basement"`
}

// Wrapper for data.gov.sg carpark info API response
type CarparkInfoAPIResponse struct {
	Success bool              `json:"success"`
	Result  CarparkInfoResult `json:"result"`
}

type CarparkInfoResult struct {
	Records []CarparkInfoResponse `json:"records"`
}

// Real-time availability from data.gov.sg API
type CarparkAvailabilityResponse struct {
	ApiInfo    map[string]interface{} `json:"api_info"`
	Items      []AvailabilityItem     `json:"items"`
}

type AvailabilityItem struct {
	Timestamp   string         `json:"timestamp"`
	CarparkData []CarparkData  `json:"carpark_data"`
}

type CarparkData struct {
	CarparkNumber  string     `json:"carpark_number"`
	UpdateDateTime string     `json:"update_datetime"`
	CarparkInfo    []LotInfo  `json:"carpark_info"`
}

type LotInfo struct {
	TotalLots     string `json:"total_lots"`
	LotType       string `json:"lot_type"`
	LotsAvailable string `json:"lots_available"`
}