
const SERVER_URL = "http://lastsprint.dev:6656"

function loadAllProjects(project, board, callback) {

    const url= SERVER_URL + `/project_price_validator/analytics/${project}/${board}`;

    console.log(url)

    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);

    xhr.onload = function (e) {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                callback(JSON.parse(xhr.response))
            } else {
                console.error(xhr.statusText);
            }
        }
    };
    xhr.onerror = function (e) {
        console.error(xhr.statusText);
    };

    xhr.timeout = 100 * 10000

    xhr.send(null);

    // http.open("GET", url, false);
    // http.send(null);
    //
    // console.log()
    //
    // return http.response
}