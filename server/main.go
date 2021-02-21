package main

import (
	"github.com/LastSprint/FeaturePriceCalculator/apicontroller"
	"github.com/LastSprint/FeaturePriceCalculator/busines"
	"github.com/LastSprint/JiraGoIssues/services"
	"os"
)

const (
	JiraBaseUrl string = "JIRA_BASE_URL"
	JiraPass string = "JIRA_PASSWORD"
	JiraLogin string = "JIRA_LOGIN"
)

func main() {

	jiraUrl := envOrCurrent(JiraBaseUrl, "https://jira.surfstudio.ru/rest/api/2/search", false)
	jiraPass := envOrCurrent(JiraPass, "", true)
	jiraLogin := envOrCurrent(JiraLogin, "", true)

	controller := &apicontroller.Api{
		PreSaleToJiraMapper: &busines.PreSaleToJiraLinker{
			LinkTableProvider: &busines.LinkTableFsLoader{PathToFile: "server/db.json"},
			JiraEpicsProvider: &busines.JiraService{Loader: services.NewJiraIssueLoader(jiraUrl, jiraLogin, jiraPass)},
		},
		BaseUrl:             "/project_price_validator",
		ListenAddress:       ":6656",
	}

	controller.Start()
}

func envOrCurrent(key string, def string, unset bool) string {

	env := os.Getenv(key)

	if len(env) == 0 {
		return def
	}

	if unset {
		os.Unsetenv(key)
	}

	return env
}