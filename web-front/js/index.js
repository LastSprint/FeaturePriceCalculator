

function init(project, board) {

    var vue = new Vue({
        el: '#app',
        data: {
            loading: true
        }
    })

    loadAllProjects(project, board, function (result) {
        vue.loading = false
        vue.loaded = result
        setTimeout(function () {
            // runVue(result)
            console.log(result)
        }, 30)
    })
}
