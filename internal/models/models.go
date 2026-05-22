package models

import "time"

type Carpark struct {
	CarparkCode   string
	CarparkName   string
	DataSource    string  // "ura" | "hdb"
	CarparkType   *string // HDB only; nil for URA
	ParkingSystem *string // "electronic" | "coupon"; nil if unknown
	Lat           *float64
	Lon           *float64
	TotalLots     *int // parkCapacity (URA) or total_lots from availability (HDB)
}

type ShortTermRate struct {
	CarparkCode   string
	DataSource    string
	VehicleType   string  // "C" | "M" | "H"
	DayType       string  // "weekday" | "saturday" | "sunday_ph" | "all"
	StartTime     string  // "HH:MM"; empty string → NULL
	EndTime       string  // "HH:MM"; empty string → NULL
	RatePer30Min  float64
	MinDuration   string
}

type SeasonRate struct {
	CarparkCode string
	DataSource  string
	VehicleType string // "C" | "M" | "H"
	TicketType  string // "Commercial" | "Residential"
	ParkingHrs  string
	MonthlyRate float64
}

type Availability struct {
	CarparkCode   string
	DataSource    string
	VehicleType   string // "C" | "M" | "H"
	LotsAvailable int
	TotalLots     *int // nil for URA (no total in availability endpoint)
	SnapshotTime  time.Time
}

type Features struct {
	CarparkCode        string
	DataSource         string
	ShortTermParking   string
	FreeParking        string
	NightParking       bool
	CarParkDecks       int
	GantryHeight       float64
	CarParkBasement    bool
	IsCentralArea      bool
	IsPeakHourCarpark  bool
}
