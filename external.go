package main

import (
	"net/http"
//	"net/url"
	"strings"
)

func msUpdatesDomainRequest(req *http.Request) bool {
	if strings.HasSuffix(req.Host, "download.windowsupdate.com") {
		return true
	}
	if strings.HasSuffix(req.Host, "download.microsoft.com") {
		return true
	}
	return false
}

func cachableUpdatesRequest(req *http.Request) bool {
	if ((req.Method == "GET" ||( cacheHead && req.Method == "HEAD") )&& msUpdatesDomainRequest(req)) {
		switch {
		// Blacklisting any option to cache antivirus definntions updates
		case (strings.Contains(req.URL.Path, "/DefinitionUpdates/")):
			return false
		case (strings.HasPrefix(req.URL.Path, "/d/") || strings.HasPrefix(req.URL.Path, "/c/")):
			return true
		case (strings.HasSuffix(req.URL.Path, ".exe") || strings.HasSuffix(req.URL.Path, ".cab") || strings.HasSuffix(req.URL.Path, ".msi") || strings.HasSuffix(req.URL.Path, ".psf") || strings.HasSuffix(req.URL.Path, ".esd") || strings.HasSuffix(req.URL.Path, ".msu")):
			return true
		default:
			return false
		}
	} else {
		return false
	}
}
