package models

// Epic is an representation of Jira epic only with appropriate fields
type Epic struct {
	// Name is summary if issue which is epic
	Name string
	// JiraKey is issue's key in jira
	JiraKey string
	// JiraLink is full URL to this issue in jira.
	JiraLink string
	// TimeSpendSum spent time sum for all tasks, bugs, service-tasks inside this epic
	TimeSpendSum float64
}