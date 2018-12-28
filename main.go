package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	gophish "./lib"
)

const defaulPort = 8080
const defaultLogFile = "log.txt"

func main() {

	webPtr := flag.String("web", "clip.unl.pt", "Website to be served.")
	buildPtr := flag.Bool("build", false, "Build the website before serving.")
	servePtr := flag.Bool("serve", false, "Wait for connections at the specified port.")
	portPtr := flag.Uint("port", defaulPort, "Port to serve on.")
	targetsPtr := flag.String("targets", "", "Filename of a text file containing one target username per line.")

	flag.Parse()
	_website, ok := gophish.WebsiteMap[*webPtr]

	if !ok {

		keys := []string{}
		for k := range gophish.WebsiteMap {
			keys = append(keys, k)
		}

		log.Fatal("Not a valid website. Expecting: " + strings.Join(keys, ", "))
	}

	website := _website()

	if *targetsPtr != "" {

		contents, err := ioutil.ReadFile(*targetsPtr)
		if err != nil {
			log.Fatal(err)
		}

		targets := strings.Split(strings.Replace(string(contents), "\r", "", -1), "\n")
		website.Targets = targets

		log.Printf("Using targets: %s", strings.Join(targets, ", "))
	}

	if !*buildPtr && !*servePtr {
		flag.PrintDefaults()
		return
	}

	if *buildPtr {
		gophish.Build(website)
	}

	if *servePtr {
		website.Serve(*portPtr)
	}
}
