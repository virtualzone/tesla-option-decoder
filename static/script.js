(function () {
    "use strict";

    function queryAjax(method, url, data) {
        let headers = new Headers();
        if (data) {
            headers.append("Content-Type", "application/json");
        }
        let options = {
            method: method,
            mode: "cors",
            cache: "no-cache",
            credentials: "same-origin",
            headers: headers
        };
        if (data) {
            options.body = JSON.stringify(data);
        }
        return window.fetch(url, options);
    }

    function disableForm(disable) {
        let inputList = Array.prototype.slice.call(document.getElementsByTagName("input"));
        inputList.forEach(input => input.disabled = disable);
    }

    document.getElementsByTagName("form").item(0).onsubmit = function (e) {
        e.preventDefault();
        disableForm(true);
        let payload = {
            url: document.getElementById("url").value
        };
        queryAjax("POST", "/api/optioncodes", payload).then(res => {
            if (res.status !== 200) {
                alert("Could not load option codes.");
                disableForm(false);
                return;
            }
            res.json().then(json => {
                for (let code in json) {
                    let item = json[code];
                    let tr = document.createElement("tr");
                    let td1 = document.createElement("td");
                    let td2 = document.createElement("td");
                    let td3 = document.createElement("td");
                    td1.textContent = code;
                    td2.textContent = item.title;
                    td3.textContent = item.description;
                    tr.append(td1, td2, td3);
                    document.getElementsByTagName("tbody").item(0).append(tr);
                }
                document.getElementsByTagName("table").item(0).style.display = "table";
            });
        }).catch(e => {
            console.log(e);
            alert("Could not load option codes.");
            disableForm(false);
        });
    };
}());