package busines

import (
	"github.com/LastSprint/FeaturePriceCalculator/models"
	"github.com/pkg/errors"
	"sync"
)

// LinkTableProvider provides link table. What is it? Look at PreSaleToJiraLinker docs
type LinkTableProvider interface {
	Load() ([]models.PreSaleFeatureRawModel, error)
}

// JiraEpicsProvider can load jira epics
type JiraEpicsProvider interface {
	// GetEpicsIssues load issues for one epic
	GetEpicsIssues(epicJiraKey string) ([]models.Issue, error)
	// GetAllEpicsExcept load all epics which are not contained in `epicsJiraKeys` array
	GetAllEpicsExcept(project, board string, epicsJiraKeys []string) ([]models.JiraIssueShortDescription, error)
}

// PreSaleToJiraLinker do next things:
// 1. Load (wherever) link table ($pre_sale_feature$ 1 -> m $jira_epic$)
// 2. Load all issues for each epic and sum time spent for it
// 3. Load epics that weren't written in link table
// 4. Make operation 2 for them
// 5. Returns results
type PreSaleToJiraLinker struct {
	LinkTableProvider
	JiraEpicsProvider
}

func (a *PreSaleToJiraLinker) Run(project, board string) (*models.LinkedPreSaleAndJira, error) {
	linkTable, err := a.LinkTableProvider.Load()

	if err != nil {
		return nil, errors.WithMessage(err, "can't load link table")
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	var features []models.PreSaleFeature
	var linkErr error

	go func() {
		features, linkErr = a.analyzeLinkTable(linkTable)
		wg.Done()
	}()

	var unlinked []models.Epic
	var unLinkErr error

	go func() {
		unlinked, unLinkErr = a.handleUnlinkedEpics(linkTable, board, project)
		wg.Done()
	}()

	wg.Wait()

	if linkErr != nil {
		return nil, linkErr
	}

	if unLinkErr != nil {
		return nil, unLinkErr
	}

	return &models.LinkedPreSaleAndJira{
		LinkedFeatures:              features,
		EpicsWithoutPreSaleFeatures: unlinked,
	}, nil
}

func (a *PreSaleToJiraLinker) analyzeLinkTable(linkTable []models.PreSaleFeatureRawModel) ([]models.PreSaleFeature, error) {
	features := make([]models.PreSaleFeature, len(linkTable))

	for i, item := range linkTable {

		mapped := make([]models.JiraIssueShortDescription, len(item.JiraEpicLinks))

		for j, it := range item.JiraEpicLinks {
			mapped[j] = models.JiraIssueShortDescription{
				Key:  it,
				Name: "",
			}
		}

		epics, err := a.loadEpics(mapped)

		if err != nil {
			return nil, errors.WithMessagef(err, "while handling pre sale feature %s", item.Name)
		}

		features[i] = models.PreSaleFeature{
			Name:                              item.Name,
			Estimate:                          item.Estimate,
			RealEpicsForCurrentPreSaleFeature: epics,
		}
	}

	return features, nil
}

func (a *PreSaleToJiraLinker) loadEpics(epics []models.JiraIssueShortDescription) ([]models.Epic, error) {

	result := make([]models.Epic, len(epics))

	for i, epic := range epics {
		issues, err := a.JiraEpicsProvider.GetEpicsIssues(epic.Key)

		if err != nil {
			 return nil, errors.WithMessagef(err, "while loading epics %s", epic)
		}

		result[i] = models.Epic{
			// TODO := Seems like that it needs only for debugging
			Name:         epic.Name,
			JiraKey:      epic.Key,
			JiraLink:     "",
			TimeSpendSum: sumIssues(issues),
		}
	}

	return result, nil
}

func (a *PreSaleToJiraLinker) handleUnlinkedEpics(linkTable []models.PreSaleFeatureRawModel, board, project string) ([]models.Epic, error) {
	exclude := []string{}

	for _, item := range linkTable {
		for _, epic := range item.JiraEpicLinks {
			exclude = append(exclude, epic)
		}
	}

	epics, err := a.JiraEpicsProvider.GetAllEpicsExcept(project, board, exclude)

	if err != nil {
		return nil, err
	}

	return a.loadEpics(epics)
}

func sumIssues(issues []models.Issue) float64 {
	result := 0.0

	for _, item := range issues {
		result += item.TimeSpent
	}

	return result
}
