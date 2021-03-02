package busines

import (
	"github.com/LastSprint/FeaturePriceCalculator/models"
	models2 "github.com/LastSprint/JiraGoIssues/models"
	"github.com/pkg/errors"
	"sync"
)

// JiraEpicsProvider can load jira epics
type JiraDataProvider interface {
	GetAllEpics(project, board string) ([]models2.IssueEntity, error)
	// GetEpicsIssues load issues for one epic
	GetEpicsIssues(epicJiraKey string) ([]models.Issue, error)
}

type JiraEpicsAnalytics struct {
	JiraDataProvider
}

func (a *JiraEpicsAnalytics) Run(project, board string) (*models.EpicsAnalytics, error) {
	epics, err := a.JiraDataProvider.GetAllEpics(project, board)

	if err != nil {
		return nil, errors.WithMessage(err, "can't load link table")
	}

	var estimated = []models2.IssueEntity{}
	var notEstimated = []models2.IssueEntity{}

	for _, it := range epics {
		if it.Fields.Estimate <= 60*60 {
			notEstimated = append(notEstimated, it)
		} else {
			estimated = append(estimated, it)
		}
	}

	var features []models.EpicAnalyze
	var unlinked []models.EpicAnalyze

	var featuresErr error
	var unlinkedErr error

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		features, featuresErr = a.analyze(estimated)
		wg.Done()
	}()

	go func() {
		unlinked, unlinkedErr = a.analyze(notEstimated)
		wg.Done()
	}()

	wg.Wait()

	if featuresErr != nil {
		return nil, featuresErr
	}

	if unlinkedErr != nil {
		return nil, unlinkedErr
	}

	return &models.EpicsAnalytics{
		LinkedFeatures:              features,
		EpicsWithoutPreSaleFeatures: unlinked,
	}, nil
}

func (a *JiraEpicsAnalytics) analyze(issues []models2.IssueEntity) ([]models.EpicAnalyze, error) {
	arr := []models.EpicAnalyze{}

	for _, it := range issues {
		issues, err := a.JiraDataProvider.GetEpicsIssues(it.Key)

		if err != nil {
			return nil, err
		}

		arr = append(arr, models.EpicAnalyze{
			Name:     it.Fields.Summary,
			Estimate: float64(it.Fields.Estimate),
			JiraKey:  it.Key,
			SpentSum: sumIssues(issues),
		})
	}

	return arr, nil
}
