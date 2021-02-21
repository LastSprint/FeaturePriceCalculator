package models

// PreSaleFeatureRawModel describes input data format
// which contains pre sale features linked with jira epics
type PreSaleFeatureRawModel struct {
	Name string
	// Estimate in seconds
	Estimate float64
	// JiraEpicLinks array of URL to jira epics
	// can't be empty
	JiraEpicLinks []string
}
