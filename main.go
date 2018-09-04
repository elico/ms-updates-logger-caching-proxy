package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/elico/requeststore-v1"
)

// A very simple http proxy

var httpClient *http.Client
var tr *http.Transport

const (
	MaxIdleConnections int = 50
	RequestTimeout     int = 90
	defaultListen          = ":8080"
	defaultDir             = "./storedata"
)

var (
	listen          string
	hostname        string
	addXCacheHeader bool
	roundtripper    bool
	verbose         bool
	retries         int
	hashSum         bool
	useDisk         bool
	// private      bool
	dir       string
	dumpHTTP  bool
	cacheHead bool
)

var store requeststore.Store

func init() {
	httpClient = createHTTPClient()
	tr = newTransport()

	flag.StringVar(&listen, "listen", defaultListen, "the host and port to bind to")
	flag.StringVar(&hostname, "hostname", "wupdate-cacher", "The hostname to showup in the X-Cache header")
	flag.BoolVar(&verbose, "v", false, "show verbose output and debugging")
	flag.BoolVar(&roundtripper, "rt", true, "Use GoLang RoundTripper instead of httpClient")
	flag.BoolVar(&cacheHead, "head", false, "Cache HEAD request, FOR debug only")
	flag.BoolVar(&addXCacheHeader, "add-x-cache", true, "Add X-Cache HIT header for responses served from localy from cache")
	flag.IntVar(&retries, "retries", 4, "The number of http and connect retries")

	flag.BoolVar(&hashSum, "hash-sum", true, "Calculate the sha256 digest of a response content enabled by default.")

	flag.StringVar(&dir, "dir", defaultDir, "the dir to store cache data in, implies -disk")
	//	flag.BoolVar(&useDisk, "disk", false, "whether to store cache data to disk")
	//	flag.BoolVar(&private, "private", false, "make the cache private")
	flag.BoolVar(&dumpHTTP, "dumpHTTP", false, "dumps http requests and responses to stdout")

	flag.Parse()
}

func main() {
	// dir := "./storedata"

	log.Printf("storing cached resources in %s", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatal(err)
	}
	var err error
	store, err = requeststore.NewDiskStore(dir)
	if err != nil {
		log.Fatal(err)
	}

	mux = &MyServer{i: 1}

	http.Handle("/", handlerfunctoHandlerfunc(simpleProxyHandlerFunc))
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		L.Println("ListenAndServe: ", err)
	}
}
