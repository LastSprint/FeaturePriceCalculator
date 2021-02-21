function init() {

    var vue = new Vue({
        el: '#app',
        data: {
            loading: true
        }
    })

    loadAllProjects(function (result) {
        vue.loading = false
        vue.loaded = result
        setTimeout(function () {
            // runVue(result)
            console.log(result.EpicsWithoutPreSaleFeatures)
        }, 30)
    })
}
