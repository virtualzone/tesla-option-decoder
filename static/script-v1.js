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
        Array.prototype.slice.call(document.getElementsByTagName("input")).forEach(input => input.disabled = disable);
        Array.prototype.slice.call(document.getElementsByTagName("textarea")).forEach(input => input.disabled = disable);
    }

    document.getElementById("mode-codes").onchange = function(e) {
        document.getElementById("url").required = false;
        document.getElementById("codes").required = true;
    }

    document.getElementById("mode-url").onchange = function(e) {
        document.getElementById("url").required = true;
        document.getElementById("codes").required = false;
    }

    function handleResponse(res) {
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
    }
    
    document.getElementsByTagName("form").item(0).onsubmit = function (e) {
        e.preventDefault();
        disableForm(true);
        if (document.getElementById("mode-url").checked) {
            let payload = {
                url: document.getElementById("url").value
            };
            queryAjax("POST", "/api/optioncodes/url", payload).then(res => handleResponse(res)).catch(e => {
                console.log(e);
                alert("Could not load option codes.");
                disableForm(false);
            });
        } else {
            let payload = {
                codes: document.getElementById("codes").value.split(/[,\.(\r\n)\r\n\t\s|]/).filter(el => el !== "")
            };
            queryAjax("POST", "/api/optioncodes/codes", payload).then(res => handleResponse(res)).catch(e => {
                console.log(e);
                alert("Could not load option codes.");
                disableForm(false);
            });
        }
    };
}());