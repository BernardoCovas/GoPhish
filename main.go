package main

import (
    "flag"
    "fmt"
    "net/http"
    "log"
)

const INDEX_RAW_HTML = "index.raw.html"
const INDEX_HTML = "index.html"

func main() {

    buildPtr := flag.Bool("build", false, "Wether to gather the resources of " + INDEX_RAW_HTML + ". This will exit after build, and not serve.")
    resFolderPtr := flag.String("resFolder", "__res__", "The resource folder.")
    portPtr := flag.Uint("port", 8080, "Port to serve on.")

    flag.Parse()

    if (*buildPtr) {
        build()
    } else {
        serve(*portPtr, *resFolderPtr)
    }

}

func serve(port uint, resFolder string) {

    http.HandleFunc("/",
        func(w http.ResponseWriter, r *http.Request) {
            serveFile(w, r, INDEX_HTML)
    })

    http.HandleFunc("/post/",
        func(w http.ResponseWriter, r *http.Request) {

            fmt.Println(r.Method)
            if (r.Method != "POST") {
               fmt.Fprintln(w, "404 Not found.")
               return
            }

            fmt.Fprintf(w, r.RemoteAddr)
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

