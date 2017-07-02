package main

import (
	"./requeststore"
	"fmt"
	"net/http"
	"os"
)

var mux http.Handler

//func hanlderToFunc(cache http.Handler, proxy http.Handler) http.HandlerFunc {
func handlerfunctoHandlerfunc(proxy http.HandlerFunc) http.HandlerFunc {
	//return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" || req.Method == "HEAD" {
			if msUpdatesDomainRequest(req) && cachableUpdatesRequest(req) {
				err := store.StoreRequest(*req, req.Method+":"+req.URL.String(), false)
				if verbose {
					switch {
					case err == requeststore.ErrFoundInStore:
						fmt.Fprintln(os.Stderr, "Store error:", err, "Request Already exist")
					case err == requeststore.ErrFoundInStorePrivate:
						fmt.Fprintln(os.Stderr, "Store error:", err, "Request In-transit lock exist")
					case err != nil:
						fmt.Fprintln(os.Stderr, "Store error:", err)
					default:

					}
				}
				switch {
				case req.Method == "GET":
					exists, err := store.RetrieveResponse(req.Method + ":" + req.URL.Scheme + "://msupdates.ngtech.internal" + req.URL.Path)
					if err == nil {
						if verbose {
							fmt.Fprintln(os.Stderr, "Store.RetrieveResponse error:", err)
						}
						exists.Close()
						mux.ServeHTTP(res, req)
						if verbose {
							fmt.Fprintln(os.Stderr, "Store : After mux")
						}
						return
					}
				case req.Method == "HEAD":
					_, err := store.RetrieveResponseHeader(req.Method + ":" + req.URL.Scheme + "://msupdates.ngtech.internal" + req.URL.Path)
					if err == nil {
						if verbose {
							fmt.Fprintln(os.Stderr, "Store.RetrieveResponseHeader error:", err)
						}

						mux.ServeHTTP(res, req)
						if verbose {
							fmt.Fprintln(os.Stderr, "Store : After mux")
						}
						return
					}
				default:
					//All other cases should be proxied
				}
			}
			proxy.ServeHTTP(res, req)
			return

		} else {
			proxy.ServeHTTP(res, req)
			return
		}
	})
}

type MyServer struct {
	i int
}

func (myServer *MyServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("Serving locally---------------")
	var object *requeststore.Response
	var err error
	switch {
	case req.Method == "GET":
		object, err = store.RetrieveResponse(req.Method + ":" + req.URL.Scheme + "://msupdates.ngtech.internal" + req.URL.Path)
		if err != nil {
			// Handle error
			if verbose {
				fmt.Printf("%s not found\n", req.URL.Path)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "<html><body style='font-size:100px'>four-oh-four</body></html>")
			return
		}
		defer object.Close()
	case req.Method == "HEAD":
	default:
		//
	}
	if verbose {
		fmt.Printf("Serving locally %s\n", req.URL.String())
		//t := fileStat.ModTime()
		//fmt.Printf("time %+v\n", t)
		fmt.Println(req.Header.Get("Range"))
	}

	// Copy original Headers from the header of the request
	headerObject, err := store.RetrieveResponseHeader(req.Method + ":" + req.URL.Scheme + "://msupdates.ngtech.internal" + req.URL.Path)
	req.Header.Del("If-Unmodified-Since")
	if err == nil {
		requeststore.CopyHeaders(headerObject.Header, w.Header())
	}
	if addXCacheHeader {
		w.Header().Set("X-Cache", "HIT from "+hostname)
	}
	switch {
	case req.Method == "GET":
		http.ServeContent(w, req, req.URL.Path, object.ResponseTime, object)
	case req.Method == "HEAD":
		w.WriteHeader(headerObject.StatusCode)
	default:
		//
	}
	if verbose {
		fmt.Printf("Served locally %s\n", req.URL.String())
	}
	switch {
	case req.Method == "HEAD":
		fmt.Printf("Served locally %s %s %s\n", req.Method, req.URL.String(), "-1")
	default:
		fmt.Printf("Served locally %s %s %s\n", req.Method, req.URL.String(), w.Header().Get("Content-Length"))
	}
}
