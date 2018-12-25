package main

import (
	"flag"

	"./build"
	"./serve"
)

const defaulPort = 8080
const defaultLogFile = "log.txt"

var website = map[string]func(){
	"clip.unl.pt": build.ClipUnlPt,
}

func main() {

	webPtr := flag.String("serve", "clip.unl.pt", "Website to be built.")
	buildPtr := flag.Bool("build", false, "Setup the server before serving the specified website.")
	logFilePtr := flag.String("logFile", defaultLogFile, "The file to write logs.")
	portPtr := flag.Uint("port", defaulPort, "Port to serve on.")

	flag.Parse()
	web, ok := website[*webPtr]

	if !ok {
		return
	}

	if *buildPtr {
		web()
	}

	serve.Serve(*webPtr, *portPtr)
}
