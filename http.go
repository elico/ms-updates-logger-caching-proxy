package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// resp implements http.ResponseWriter writing
type dummyResp struct {
	io.Writer
	h int
}

func newDummyResp() http.ResponseWriter {
	return &dummyResp{Writer: &bytes.Buffer{}}
}

func (w *dummyResp) Header() http.Header {
	return make(http.Header)
}

func (w *dummyResp) WriteHeader(h int) {
	w.h = h
}

func (w *dummyResp) String() string {
	return fmt.Sprintf("[%v] %q", w.h, w.Writer)
}

func (w *dummyResp) Write(buf []byte) (int, error) {
	w.WriteHeader(len(buf))
	return len(buf), nil
}

/*
type HttpInterceptor struct {
    origWriter http.ResponseWriter
    overridden bool
}

func (i *HttpInterceptor) WriteHeader(rc int) {
        i.origWriter.WriteHeader(rc)
        return

    // if the default case didn't execute (and return) we must have overridden the output
    //i.overridden = true
    //log.Println(i.overridden)
}

func (i *HttpInterceptor) Write(b []byte) (int, error) {
    if !i.overridden {
        return i.origWriter.Write(b)
    }

    // Return nothing if we've overriden the response.
    return 0, nil
}

func (i *HttpInterceptor) Header() http.Header {
    return i.origWriter.Header()
}
*/
