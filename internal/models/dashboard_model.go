package models

type DashboardModel struct {
	UserData        UserData
	BugReportData   WeeklyBugReportData
	TransactionData WeeklyTransactionData
}

type UserData struct {
	Total    int `json:"total"`
	Verified int `json:"verified"`
	Expert   int `json:"expert"`
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

type WeeklyTransactionData struct {
	Total      int                    `json:"total"`
	Daily      []DailyTransactionData `json:"daily"`
	Difference int                    `json:"difference"` // Difference compared to last week
}

type DailyTransactionData struct {
	Date  string `json:"date" db:"date"`
	Count int    `json:"count" db:"count"`
}
