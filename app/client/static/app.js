function makeRequest(url, body, callback) {
    const oReq = new XMLHttpRequest();

    if (callback !== null) {
        oReq.onload = function() {
            callback(oReq.response);
        }
    }

    oReq.responseType = "json";

    if (body) {
        oReq.open("POST", url);
        oReq.send(JSON.stringify(body));        
    } else {        
        oReq.open("GET", url);
        oReq.send();
    }
}

let searchResults = [];

function search() {
    makeRequest("/search?query=" + document.getElementById("query").value, null, function(response) {        
        searchResults = response["Results"]
        .sort((a, b) => stringSort(a.Names[0], b.Names[0]))
        .map(res => ({...res, isHeader: false}));

        searchResults = [{isHeader: true}, ...searchResults]
        document.getElementById("results").replaceChildren(...searchResults.map((result, i) => {
            let resultDiv = document.createElement("div");
            if (result.isHeader) {
                resultDiv.innerHTML = `<div class="bold">Number of Seeders</div><div class="bold">File Size</div><div class="bold">File Name</div><div></div>`;
                return resultDiv
            }

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

function stringSort(a, b) {    
    return a.toLowerCase() > b.toLowerCase()
}

function startDownload(i) {
    makeRequest("/start-download", searchResults[i], null);
}

function getProgress() {
    makeRequest("/progress", null, function(response) {
        document.getElementById("status").replaceChildren(...response
            .sort((a, b) => stringSort(a.Name, b.Name))
            .map((status, i) => {
                let statusDiv = document.createElement("div");        
                statusDiv.style.background = `linear-gradient(to right, #03e303 ${status['Progress'] * 100}%, white ${status['Progress'] * 100}%)`;
                statusDiv.innerText = status['Name'];        
                return statusDiv;
        }));
    });
}

setInterval(getProgress, 100);
getProgress();

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