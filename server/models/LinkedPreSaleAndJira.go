package models

// LinkedPreSaleAndJira is a result of links pre sale features and real jira epics
type LinkedPreSaleAndJira struct {
	// LinkedFeatures array of links between pre sale features and jira epics
	LinkedFeatures []PreSaleFeature

	// EpicsWithoutPreSaleFeatures array of epics which don't have pre sale features to link
	EpicsWithoutPreSaleFeatures []Epic
}

type EpicsAnalytics struct {
	LinkedFeatures              []EpicAnalyze
	EpicsWithoutPreSaleFeatures []EpicAnalyze
}
