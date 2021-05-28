function makeRequest(url, body, callback) {
    const oReq = new XMLHttpRequest();

    if (callback !== null) {
        oReq.onload = function() {
            callback(oReq.response);
        }
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

let searchResults = [];

function search() {
    makeRequest("/search?query=" + document.getElementById("query").value, null, function(response) {
        searchResults = response["Results"];

        document.getElementById("results").replaceChildren(...searchResults.map((result, i) => {
            let resultDiv = document.createElement("div");

            resultDiv.innerHTML = `<div>${result['Clients'].length}</div><div>${nicerSize(result['Size'])}</div><div>${result['Names'][0]}</div>`;

            let downloadButton = document.createElement("button");
            downloadButton.innerText = 'Download';
            downloadButton.onclick = function() {
                search(startDownload(i))
            }

            let downloadButtonContainer = document.createElement("div");
            downloadButtonContainer.appendChild(downloadButton);

            resultDiv.appendChild(downloadButtonContainer);

            return resultDiv;
        }));
    });
}

function startDownload(i) {
    makeRequest("/start-download", searchResults[i], null);
}

function nicerSize(size) {
    if (size > 1000 ** 3) {
        return (size / (1000 ** 3)) + ' GB';
    } else if (size > 1000 ** 2) {
        return (size / (1000 ** 2)) + ' MB';
    } else if (size > 1000 ** 1) {
        return (size / (1000 ** 1)) + ' KB';
    } else {
        return size + ' B';
    }
}