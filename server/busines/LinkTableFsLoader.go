package busines

import (
	"encoding/json"
	"github.com/LastSprint/FeaturePriceCalculator/models"
	"github.com/pkg/errors"
	"os"
)

type LinkTableFsLoader struct {
	PathToFile string
}

type linkTableJsonView struct {
	Name          string   `json:"name"`
	Estimate      float64  `json:"estimate"`
	JiraEpicLinks []string `json:"jiraEpicLinks"`
}

func (l *LinkTableFsLoader) Load() ([]models.PreSaleFeatureRawModel, error) {
	file, err := os.OpenFile(l.PathToFile, os.O_RDONLY, os.ModePerm)

	if err != nil {
		return nil, errors.WithMessagef(err, "while reading file at %s", l.PathToFile)
	}

	var res []linkTableJsonView

	if err = json.NewDecoder(file).Decode(&res); err != nil {
		return nil, errors.WithMessagef(err, "while decoding link table from %s to json", &l.PathToFile)
	}

	result := make([]models.PreSaleFeatureRawModel, len(res))

	for i, it := range res {
		result[i] = models.PreSaleFeatureRawModel{
			Name:          it.Name,
			Estimate:      it.Estimate,
			JiraEpicLinks: it.JiraEpicLinks,
		}
	}

	return result, nil
}
