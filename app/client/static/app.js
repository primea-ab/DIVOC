function makeRequest(url, body, callback) {
    const oReq = new XMLHttpRequest();
    oReq.onload = function() {
        callback(oReq.response);
    }
    oReq.responseType = "json";

    if (body === null) {
        oReq.open("GET", url);
        oReq.send();
    } else {
        oReq.open("POST", url);
        oReq.send(JSON.stringify(body));
    }
}

makeRequest("/search?query=real", null, function(response) {
    makeRequest("/start-download", response["Results"][0], function() {
        console.log("download started");
    });
});