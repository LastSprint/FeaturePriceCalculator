package models

// PreSaleFeature describes relations between pre sale feature in real jira epics
// which were used to did this feature
type PreSaleFeature struct {
	// Name is pre sale feature name
	Name string
	// Estimate is pre sale feature estimate
	Estimate float64
	// RealEpicsForCurrentPreSaleFeature is array if jira epics
	RealEpicsForCurrentPreSaleFeature []Epic
}

// PreSaleFeature describes relations between pre sale feature in real jira epics
// which were used to did this feature
type EpicAnalyze struct {
	// Name is pre sale feature name
	Name string
	// Estimate is pre sale feature estimate
	Estimate float64

	JiraKey string

	SpentSum float64
}
