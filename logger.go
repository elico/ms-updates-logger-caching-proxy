package main

import (
	"log"
	"os"
)

// global logger
var L = log.New(os.Stdout, "ms-updates-cacher: ", log.Lshortfile|log.LstdFlags)
