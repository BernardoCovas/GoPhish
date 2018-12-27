package main

import (
	"flag"
	"log"
	"strings"

	"./include/build"
	"./include/common"
)

const defaulPort = 8080
const defaultLogFile = "log.txt"

func main() {

	webPtr := flag.String("web", "clip.unl.pt", "Website to be served.")
	buildPtr := flag.Bool("build", false, "Build the website before serving.")
	servePtr := flag.Bool("serve", false, "Wait for connections at the specified port.")
	portPtr := flag.Uint("port", defaulPort, "Port to serve on.")

	flag.Parse()
	_website, ok := common.WebsiteMap[*webPtr]

	if !ok {

		keys := []string{}
		for k := range common.WebsiteMap {
			keys = append(keys, k)
		}

		log.Fatal("Not a valid website. Expecting: " + strings.Join(keys, ", "))
	}

	website := _website()
	if !*buildPtr && !*servePtr {
		flag.PrintDefaults()
		return
	}

	if *buildPtr {
		build.Build(website)
	}

	if *servePtr {
		website.Serve(*portPtr)
	}
}
