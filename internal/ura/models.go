package ura

// Generic URAResponse
type URAResponse[T any] struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	Result  T      `json:"Result"`
}

type CarparkDetailsResponse struct {
	PpCode        string  `json:"ppCode"`
	PpName        string  `json:"ppName"`
	ParkingSystem string  `json:"parkingSystem"`
	VehCat        string  `json:"vehCat"`
	WeekdayRate   string  `json:"weekdayRate"`
	WeekdayMin    string  `json:"weekdayMin"`
	SatdayRate    string  `json:"satdayRate"`
	SatdayMin     string  `json:"satdayMin"`
	SunPHRate     string  `json:"sunPHRate"`
	SunPHMin      string  `json:"sunPHMin"`
	StartTime     string  `json:"startTime"`
	EndTime       string  `json:"endTime"`
	ParkCapacity  int     `json:"parkCapacity"`
	Geometries    []struct {
		Coordinates string `json:"coordinates"`
	} `json:"geometries"`
}

type CarparkAvailabilityResponse struct {
	CarparkNo 		string `json:"carparkNo"`
	LotType	 			string `json:"lotType"`
	LotsAvailable string `json:"lotsAvailable"`
	Geometries    []struct {
		Coordinates string `json:"coordinates"`
	} `json:"geometries"`
}

type CarparkSeasonDetailsResponse struct {
	PpCode        string  `json:"ppCode"`
	PpName        string  `json:"ppName"`
	VehCat        string  `json:"vehCat"`
	ParkingHrs    string  `json:"parkingHrs"`
	TicketType		string 	`json:"ticketType"`
	MonthlyRate 	string 	`json:"monthlyRate"`
	Geometries    []struct {
		Coordinates string `json:"coordinates"`
	} `json:"geometries"`
}