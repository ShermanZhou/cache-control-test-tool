package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

const port string = "8123"

//should have thread lock here.
var lastUpdate = time.Now().Local()
var LOCATION, _ = time.LoadLocation("Local")

func main() {
	router := httprouter.New()
	router.GET("/api/get", Get)
	router.NotFound = http.FileServer(http.Dir("public"))
	secondsTimer := time.NewTicker(time.Second * 20)
	defer secondsTimer.Stop()
	go func() {
		for {
			<-secondsTimer.C
			lastUpdate = time.Now().Local()
			log.Println("                   debug: lastUpdate renewed")
		}
	}()
	log.Printf("Serve on port: %s \n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

var Hdrs = struct {
	Etag            string
	LastModified    string
	Pragma          string
	IfModifiedSince string
	IfNoneMatch     string
	Expires         string
	CacheControl    string
}{
	Etag:            "Etag",
	LastModified:    "Last-Modified",
	Pragma:          "Pragma",
	IfModifiedSince: "If-Modified-Since",
	IfNoneMatch:     "If-None-Match",
	Expires:         "Expires",
	CacheControl:    "Cache-Control",
}

func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mtimeClient := r.Header.Get(Hdrs.IfModifiedSince)
	var mtimeClientT time.Time

	if mtimeClient != "" {
		mtimeClientT, _ = time.ParseInLocation(http.TimeFormat, mtimeClient, LOCATION)
	}

	tagClient := r.Header.Get(Hdrs.IfNoneMatch)
	content := time.Now().Local().Format(time.RFC3339Nano)
	tDiff := lastUpdate.Sub(mtimeClientT)
	fmt.Printf("log: modified: Svr: %s\tClt: %s\tDiff: %v\n", lastUpdate.Format(http.TimeFormat), mtimeClientT.Format(http.TimeFormat), tDiff)

	//server will issue 'using cache' based on IfModifiedSince
	if tDiff <= time.Second*1 || mtimeClientT.After(lastUpdate) {
		w.WriteHeader(http.StatusNotModified)
		log.Println("log: Send 304 due to modified date")
		return
	}
	//server will issue 'using cache' based on etag
	if false && tagClient == "abcdefg" {
		w.WriteHeader(http.StatusNotModified)
		log.Println("log: Send 304 due etag unchanged")
		return
	}

	log.Println("log: Send 200, no condition to use 304")

	w.Header().Set(Hdrs.CacheControl, "no-cache")
	w.Header().Set(Hdrs.LastModified, lastUpdate.Format(http.TimeFormat))
	w.Header().Set(Hdrs.Etag, "abcdefg")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, content)

}
