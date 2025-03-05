package models

type DashboardModel struct {
	UserData        UserData
	ReportData      ReportData
	TransactionData WeeklyTransactionData
}

// User data
type UserData struct {
	Total      int `json:"total"`
	Verified   int `json:"verified"`
	Expert     int `json:"expert"`
	Restricted int `json:"restricted"`
}

// Report Data
type ReportData struct {
	WeeklyBugReport    WeeklyBugReportData    `json:"weeklyBugReport"`
	WeeklyVendorReport WeeklyVendorReportData `json:"weeklyVendorReport"`
}

type WeeklyBugReportData struct {
	Total      int                  `json:"total"`
	Daily      []DailyBugReportData `json:"daily"`
	Difference int                  `json:"difference"` // Difference compared to last week
}

type DailyBugReportData struct {
	Date  string `json:"date" db:"date"`
	Count int    `json:"count" db:"count"`
}

type WeeklyVendorReportData struct {
	Total      int                     `json:"total"`
	Daily      []DailyVendorReportData `json:"daily"`
	Difference int                     `json:"difference"` // Difference compared to last week
}

type DailyVendorReportData struct {
	Date  string `json:"date" db:"date"`
	Count int    `json:"count" db:"count"`
}

// Transaction Data
type WeeklyTransactionData struct {
	Total      int                    `json:"total"`
	Daily      []DailyTransactionData `json:"daily"`
	Difference int                    `json:"difference"` // Difference compared to last week
}

type DailyTransactionData struct {
	Date  string `json:"date" db:"date"`
	Count int    `json:"count" db:"count"`
}
