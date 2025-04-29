package dto

type Dashboard struct {
	Users         DashboardUsers
	Services      DashboardServices
	WeeklyBooking WeeklyBookingData
	SearchTrend   Trend
	Application   Application
	Report        DashboardReports
}

type DashboardUsers struct {
	Total      int
	Reported   int
	Restricted int
	Vendor     int
	Recent     []User
}

type DashboardServices struct {
	Total  int
	Active int
}

type Trend struct {
	Searches []string
}

type WeeklyBookingData struct {
	Total      int                `json:"total"`
	Daily      []DailyBookingData `json:"daily"`
	Difference int                `json:"difference"`
}

type DailyBookingData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DashboardReports struct {
	Pending int
}
type Application struct {
	Pending int
}
