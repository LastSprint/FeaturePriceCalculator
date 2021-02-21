package busines

import (
	"fmt"
	"github.com/LastSprint/FeaturePriceCalculator/models"
	models2 "github.com/LastSprint/JiraGoIssues/models"
	"github.com/LastSprint/JiraGoIssues/services"
	"github.com/pkg/errors"
	"strings"
)

type JiraService struct {
	Loader *services.JiraIssueLoader
}

func (j *JiraService) GetEpicsIssues(epicJiraKey string) ([]models.Issue, error) {
	res, err := j.Loader.LoadIssues(services.SearchRequest{
		IncludedTypes:           []string{models2.IssueTypeBug, models2.IssueTypeTask, models2.IssueTypeServiceTask},
		EpicLink:                epicJiraKey,
	})

	if err != nil {
		return nil, errors.WithMessagef(err, "while loading epic %s from jira", epicJiraKey)
	}

	mapped := make([]models.Issue, len(res.Issues))

	for i, it := range res.Issues {
		mapped[i] = models.Issue{
			Name:      it.Fields.Summary,
			TimeSpent: float64(it.Fields.TimeSpend),
		}
	}

	return mapped, nil
}

type JQL string

func (J JQL) GetUseOnlyAdditionalFields() bool {
	return false
}

func (J JQL) MakeJiraRequest() string {
	return string(J)
}

func (J JQL) GetAdditionFields() []services.JiraField {
	return []services.JiraField{}
}

func (j *JiraService) GetAllEpicsExcept(project, board string, epicsJiraKeys []string) ([]models.JiraIssueShortDescription, error) {

	epicsSetForJql := strings.Join(epicsJiraKeys, ",")

	jql := fmt.Sprintf("project = %s and board = %s and key not in (%s) and issuetype = Epic", project, board, epicsSetForJql)

	epics, err := j.Loader.LoadIssues(JQL(jql))

	if err != nil {
		return nil, errors.WithMessagef(err, "while loading epics without excluded")
	}

	keys := make([]models.JiraIssueShortDescription, len(epics.Issues))

	for i, it := range epics.Issues {
		keys[i] = models.JiraIssueShortDescription{
			Key:  it.Key,
			Name: it.Fields.Summary,
		}
	}
	return keys, nil
}
