package main

import (
	"github.com/LastSprint/FeaturePriceCalculator/apicontroller"
	"github.com/LastSprint/FeaturePriceCalculator/busines"
	"github.com/LastSprint/JiraGoIssues/services"
	"os"
)

const (
	JiraBaseUrl   string = "JIRA_BASE_URL"
	JiraPass      string = "JIRA_PASSWORD"
	JiraLogin     string = "JIRA_LOGIN"
	PathToWeb     string = "PATH_TO_WEB"
	PathToCert    string = "FPC_TLS_CERT_PATH"
	PathToKey     string = "FPC_TLS_KEY_PATH"
	ListenAddress string = "FPC_LISTEN_ADDRESS"
)

func main() {

	jiraUrl := envOrCurrent(JiraBaseUrl, "https://jira.surfstudio.ru/rest/api/2/search", false)
	jiraPass := envOrCurrent(JiraPass, "", true)
	jiraLogin := envOrCurrent(JiraLogin, "", true)

	controller := &apicontroller.Api{
		PreSaleToJiraMapper: &busines.JiraEpicsAnalytics{
			JiraDataProvider: &busines.JiraService{Loader: services.NewJiraIssueLoader(jiraUrl, jiraLogin, jiraPass)},
		},
		BaseUrl:       "/project_price_validator",
		ListenAddress: envOrCurrent(ListenAddress, "", false),
		PathToWeb:     envOrCurrent(PathToWeb, "web-front", false),
		CertPath:      envOrCurrent(PathToCert, "", false),
		KeyPath:       envOrCurrent(PathToKey, "", false),
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
