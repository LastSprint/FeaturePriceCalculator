package apicontroller

import (
	"encoding/json"
	"github.com/LastSprint/FeaturePriceCalculator/models"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"log"
	"net/http"
	"sync"
	"time"
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

	h := cors.Default().Handler(r)

	http.ListenAndServe(a.ListenAddress, h)
}

func (a *Api) registerMethods(router *httprouter.Router) {
	router.GET(a.BaseUrl+"/analytics/:project/:board", a.analyticsHandler)
}

type ResponseCache struct {
	data models.LinkedPreSaleAndJira
	time time.Time
}

var cache = ResponseCache{}
var mutex = &sync.Mutex{}

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

	if cache.time.Add(time.Minute * 30).After(time.Now()) {
		if err := json.NewEncoder(w).Encode(cache.data); err != nil {

		}
		return
	}

	res, err := a.PreSaleToJiraMapper.Run(project, board)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"msg": err.Error()})
		return
	}

	mutex.Lock()
	cache.time = time.Now()
	cache.data = *res
	mutex.Unlock()

	s := 0.0

	for _, t := range res.EpicsWithoutPreSaleFeatures {
		s += t.TimeSpendSum
	}

	log.Printf("[INFO] Unlinked time sum %v", int(s))

	json.NewEncoder(w).Encode(res)
	return
}
