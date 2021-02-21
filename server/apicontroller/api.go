package apicontroller

import (
	"encoding/json"
	"github.com/LastSprint/FeaturePriceCalculator/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type PreSaleToJiraMapper interface {
	Run(project, board string) (*models.LinkedPreSaleAndJira, error)
}

type Api struct {
	PreSaleToJiraMapper

	BaseUrl string
	ListenAddress string
}

func (a *Api) Start() {
	r := httprouter.New()

	a.registerMethods(r)

	http.ListenAndServe(a.ListenAddress, r)
}

func (a *Api) registerMethods(router *httprouter.Router) {
	router.GET(a.BaseUrl+"/analytics/:project/:board", a.analyticsHandler)
}

func (a *Api) analyticsHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {

	project := p.ByName("project")
	board := p.ByName("board")

	if len(project) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "project parameter should be set"})
		return
	}

	if len(board) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "board parameter should be set"})
		return
	}

	res, err := a.PreSaleToJiraMapper.Run(project, board)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"msg": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(res)
	return
}
