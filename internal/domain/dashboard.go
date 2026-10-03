package domain

type ActiveProgress struct {
	Doing      int     `json:"doing"`
	Todo       int     `json:"todo"`
	Total      int     `json:"total"`
	Ratio      float64 `json:"ratio"`
	Percentage float64 `json:"percentage"`
}

type DashboardSummary struct {
	TotalActiveProjects  int                `json:"totalActiveProjects"`
	TaskCountPerStatus   map[TaskStatus]int `json:"taskCountPerStatus"`
	ActiveProgress       ActiveProgress     `json:"activeProgress"`
	InReview             int                `json:"inReview"`
	DueSoon              int                `json:"dueSoon"`
	Overdue              int                `json:"overdue"`
	OverdueTaskCount     int                `json:"overdueTaskCount"`
	WorkloadPerMember    map[string]int     `json:"workloadPerMember"`
	RecentlyUpdatedTasks []Task             `json:"recentlyUpdatedTasks"`
	RecentTasks          []Task             `json:"recentTasks"`
	UpcomingDeadlines    []Task             `json:"upcomingDeadlines"`
	HighPriorityTasks    []Task             `json:"highPriorityTasks"`
	LatestUpdates        []ActivityLog      `json:"latestUpdates"`
}
