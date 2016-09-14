package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"crypto/sha256"
)

// A very simple http proxy

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		//Timeout:       time.Duration(RequestTimeout) * time.Second,
		CheckRedirect: noRedirect,
	}

	return client
}

func newTransport() *http.Transport {
	return http.DefaultTransport.(*http.Transport)
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("REDIRECT!!!")
}

func simpleProxyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	L.Printf("%s %s %s %s\n", r.Method, r.RequestURI, r.Proto, r.RemoteAddr)

	redirect := false
	if msUpdatesDomainRequest(r)  {
		// Schedual caching of the full object
		if verbose {
			L.Println("Allowed WindowsUpdates host", r.Host)
		}
		RemoveHopHeaders(r.Header)
		if !r.URL.IsAbs() {
			r.RequestURI = "http://"+ r.Host + r.URL.String()
			r.URL, _= url.Parse(r.RequestURI)
		}

	} else {
		http.Error(w, "Request is not allowed\n", http.StatusUnauthorized)
		return

	}
/*
	} else if r.URL.IsAbs() {
		// This is an error if is not empty on Client
		r.RequestURI = ""
		RemoveHopHeaders(r.Header)
		//	} else if r.URL.Path == "/reload" {
		//		self.reload(w, r)
	} else {
		if verbose {
			L.Printf("%s is not a full URL path\n", r.RequestURI)
		}
		http.Error(w, r.RequestURI + " is not a full URL path\n", http.StatusInternalServerError)
		return
	}
*/
	start := time.Now()
	var err error
	var resp *http.Response

	for retry := 0; retry < retries; retry++ {
		if roundtripper {
			resp, err = tr.RoundTrip(r)
		} else {
			resp, err = httpClient.Do(r)
		}
		if err == nil {
			break
		} else {
			time.Sleep(1 << uint(retry) * time.Second)
		}
	}
	switch {
	case err != nil && strings.Contains(err.Error(), "REDIRECT!!!"):
		redirect = true
		_ = redirect
	case err != nil:
		switch {
		case strings.Contains(err.Error(), "REDIRECT!!!"):

		default:
			if verbose {
				L.Println(r.URL, " ", err.Error())
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:

	}

	CopyHeader(w, resp)
	w.WriteHeader(resp.StatusCode)
	n := int64(0)
	// copy content
	s := sha256.New()
	switch {
	case r.Method == "HEAD":

	case resp.Body != nil:
		defer resp.Body.Close()
		if hashSum {
			multiw := io.MultiWriter(w, s)
			n, err = io.Copy(multiw, resp.Body)
		} else {
			n, err = io.Copy(w, resp.Body)
		}
		if err != nil {
			if verbose {
				L.Printf("Copy: %s\n", err.Error())
			}
			// There is a possibility that at this stage we cannot write the header of the http response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		L.Println(r.URL, r.Method, "unknonw response Body case")
	}
	d := BeautifyDuration(time.Since(start))
	ndtos := BeautifySize(n)
	if hashSum {
		digest := s.Sum(nil)
		L.Printf("RESPONSE %s %s in %s <-%s sha256 => %x\n", r.URL, resp.Status, d, ndtos, digest)
	}
	if verbose {
		L.Printf("RESPONSE %s %s in %s <-%s\n", r.URL.Host, resp.Status, d, ndtos)
	}
}

func copyHeaders(dst, src http.Header) {
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
	if dst.Get("Content-Type") == "" {
		dst.Add("Content-Type", "    ")
	}
}

// copy and overwrite headers from r to w
func CopyHeader(w http.ResponseWriter, r *http.Response) {
        // copy headers
        dst, src := w.Header(), r.Header
        for k, _ := range dst {
                dst.Del(k)
        }
        for k, vs := range src {
                for _, v := range vs {
                        dst.Add(k, v)
                }
        }

        // A hack to disable the defaults of GoLang ot add text/html content-type on an empty response
        if dst.Get("Content-Type") == "" {
                dst.Add("Content-Type", "    ")
        }

}

var hopHeaders = []string{
        // If no Accept-Encoding header exists, Transport will add the headers it can accept
        // and would wrap the response body with the relevant reader.
        "Accept-Encoding",
        "Connection",
        "Keep-Alive",
        "Proxy-Authenticate",
        "Proxy-Authorization",
        "Te", // canonicalized version of "TE"
        "Trailers",
        "Transfer-Encoding",
        "Upgrade",
        "Proxy-Connection", // added by CURL  http://homepage.ntlworld.com/jonathan.deboynepollard/FGA/web-proxy-connection-header.html
}

func RemoveHopHeaders(h http.Header) {
        for _, k := range hopHeaders {
                h.Del(k)
        }
}
