var headers = new Headers();
headers.append("sessionid", "12345");
function zzFetch() {
    fetch("/api/get", {
        headers: headers,
        cache: "no-cache"
    }).then((value)=>{
        value.text().then(e=>console.log(e));
    })
}
function zzXHR() {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/api/get");
    xhr.setRequestHeader("Cache-Control","no-cache");
    xhr.setRequestHeader("SessionID","12345");
    xhr.setRequestHeader("Vary", "SessionId");
    xhr.onreadystatechange = (e)=>{
        console.log(xhr.responseText);
    };
    xhr.send();
}