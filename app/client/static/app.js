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
    document.getElementById("music").play();

    makeRequest("/search?query=" + document.getElementById("query").value, null, function(response) {        
        searchResults = response["Results"]
            .sort((a, b) => stringSort(a.Names[0], b.Names[0]))
            .map(res => ({...res, isHeader: false}));
        
        if (searchResults && searchResults.length !== 0) {
            searchResults = [
                {isHeader: true}, ...searchResults
            ];
        }
        

        document.getElementById("results").replaceChildren(...searchResults.map((result, i) => {
            let resultDiv = document.createElement("div");

            if (result.isHeader) {
                resultDiv.innerHTML = `<div class="bold">Number of Seeders</div><div class="bold">File Size</div><div class="bold">File Name</div><div></div>`;
                return resultDiv
            }

            resultDiv.innerHTML = `<div>${result['Clients'].length}</div><div>${nicerSize(result['Size'])}</div><div>${result['Names'][0]}</div>`;
            let downloadButton = document.createElement("div");
            downloadButton.className = 'button';
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

                if (status['Progress'] === 1) {
                    statusDiv.style.background = '#03b503';
                } else if (status['Progress']) {
                    statusDiv.style.background = `linear-gradient(to right, #03e303 ${status['Progress'] * 100}%, transparent ${status['Progress'] * 100}%)`;
                }

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

function searchOnEnter(e) {
    if (e.keyCode === 13) {
        search()
    }
}

let r = 250;
let g = 50;
let b = 50;

function start() {
    setInterval(changeColor, 25);
}

function changeColor() {
    rFactor = r > 220 ? 0.7 : (r < 25 ? 0.3 : 0.5)
    gFactor = g > 220 ? 0.7 : (g < 25 ? 0.3 : 0.5)
    bFactor = b > 220 ? 0.7 : (b < 25 ? 0.3 : 0.5)

    const diffR = Math.round((Math.random() - rFactor) * 20);
    const diffG = Math.round((Math.random() - gFactor) * 20);
    const diffB = Math.round((Math.random() - bFactor) * 20);

    r = Math.min(Math.max(0, r + diffR), 255)
    g = Math.min(Math.max(0, g + diffG), 255)
    b = Math.min(Math.max(0, b + diffB), 255)

    document.getElementById("header").style.color = `rgb(${r}, ${g}, ${b})`;
}