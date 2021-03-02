package apicontroller

import (
	"encoding/json"
	"fmt"
	"github.com/LastSprint/FeaturePriceCalculator/models"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

type PreSaleToJiraMapper interface {
	Run(project, board string) (*models.EpicsAnalytics, error)
}

type Api struct {
	PreSaleToJiraMapper

	BaseUrl       string
	ListenAddress string
	PathToWeb     string

	CertPath string
	KeyPath  string
}

func (a *Api) Start() {
	r := httprouter.New()
	a.registerMethods(r)

	r.ServeFiles("/web/wf/*filepath", http.Dir(a.PathToWeb))

	h := cors.Default().Handler(r)

	fmt.Println("BaseURL", a.BaseUrl)
	fmt.Println("ListenAddress", a.ListenAddress)
	fmt.Println("PathToWeb", a.PathToWeb)
	fmt.Println("CertPath", a.CertPath)
	fmt.Println("KeyPath", a.KeyPath)

	if len(a.CertPath) == 0 && len(a.KeyPath) == 0 {
		log.Fatal(http.ListenAndServe(a.ListenAddress, h))
	} else {
		log.Fatal(http.ListenAndServeTLS(a.ListenAddress, a.CertPath, a.KeyPath, h))
	}
}

func (a *Api) registerMethods(router *httprouter.Router) {
	router.GET(a.BaseUrl+"/analytics/:project/:board", a.analyticsHandler)
	router.GET("/web/analytcs", a.webShowProjectInputHandler)
}

type ResponseCache struct {
	data models.EpicsAnalytics
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

	//if cache.time.Add(time.Minute * 30).After(time.Now()) {
	//	if err := json.NewEncoder(w).Encode(cache.data); err != nil {
	//
	//	}
	//	return
	//}

	res, err := a.PreSaleToJiraMapper.Run(project, board)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"msg": err.Error()})
		return
	}

	//mutex.Lock()
	//cache.time = time.Now()
	//cache.data = *res
	//mutex.Unlock()

	json.NewEncoder(w).Encode(res)
	return
}

func (a *Api) webShowProjectInputHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pr := r.FormValue("project")
	board := r.FormValue("board")

	tmpl, err := template.New("analatics").Parse(tmplString)

	if err != nil {
		log.Fatal("[ERR]", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl.Execute(w, map[string]string{"project": pr, "board": board})
}

var tmplString string = `
<!DOCTYPE html>
<html>
<head>
    <link href="https://cdn.jsdelivr.net/npm/vuesax/dist/vuesax.css" rel="stylesheet">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no, minimal-ui">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/tonsky/FiraCode@4/distr/fira_code.css">
    <script src="https://cdn.jsdelivr.net/npm/chart.js@2.9.4/dist/Chart.min.js"></script>
    <script src="https://unpkg.com/vue/dist/vue.js"></script>
    <script src="https://unpkg.com/vuesax"></script>
    <script src="wf/js/network.js"></script>
    <script src="wf/js/index.js"></script>
    <script src="wf/js/cmp/UnlinkedCardCmp.js"></script>
    <script src="wf/js/cmp/LinkedTableCmp.js"></script>
    <link href="wf/all.css" rel="stylesheet">
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons"
          rel="stylesheet">
</head>
<body style="background-color: #525252">

<div id="app">

    <div v-if="loading" class="rbcenter">
        <div style="text-align: center">
            <div class="trinity-rings-spinner" style="position: relative; left: 20%">
                <div class="circle"></div>
                <div class="circle"></div>
                <div class="circle"></div>
            </div>
            <p style="font-family: 'Fira Code', monospace; color: #f5f3f1;margin-top: 16px;">Loading...</p>
            <p style="font-family: 'Fira Code', monospace; color: #f5f3f1">It's not as easy as you think 🥲...</p>
        </div>
    </div>

    <div v-else>
        <fpc-main v-bind:reply="loaded"></fpc-main>
    </div>
</div>

<script>
    Vue.component(TableCmpName, TableCmp)
    Vue.component(UnlinkedJiraEpicLineCmpName, UnlinkedJiraEpicLineCmp)
    Vue.component(UnlinkedCardCmpName, UnlinkedCardCmp)
    Vue.component(MainCmpName, MainCmp)
    init({{.project}}, {{.board}})

</script>
</body>
</html>
`
