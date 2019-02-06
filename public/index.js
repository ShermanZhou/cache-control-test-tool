var headers = new Headers();
function zzFetch() {
    fetch("/api/get", {
        headers: headers,
        cache: "no-cache"
    });
}
setInterval(function(){
    zzFetch();
},5000);