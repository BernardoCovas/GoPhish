package main

import (
	"flag"
	"log"
	"strings"

	"./build"
	"./serve"
)

const defaulPort = 8080
const defaultLogFile = "log.txt"

func main() {

	webPtr := flag.String("web", "clip.unl.pt", "Website to be served.")
	buildPtr := flag.Bool("build", false, "Build the website before serving.")
	servePtr := flag.Bool("serve", false, "Wait for connections at the specified port.")
	// logFilePtr := flag.String("logFile", defaultLogFile, "The file to write logs.")
	portPtr := flag.Uint("port", defaulPort, "Port to serve on.")

	flag.Parse()
	website, ok := build.WebsiteMap[*webPtr]

	if !ok {

		keys := []string{}
		for k := range build.WebsiteMap {
			keys = append(keys, k)
		}

		log.Fatal("Not a valid website. Expecting: " + strings.Join(keys, ", "))
	}

	if *buildPtr {
		website.Build()
	}

	if *servePtr {
		serve.Serve(website, *portPtr)
	}
}
