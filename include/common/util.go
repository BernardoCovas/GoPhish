package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// HandleLogin extracts the "user" and "pass" fields
// (if available) from a request.
func HandleLogin(r *http.Request) (string, string) {

	if r.Method != "POST" {
		return "", ""
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return "", ""
	}

	values, err := url.ParseQuery(string(body))

	if err != nil {
		log.Print(err)
		return "", ""
	}

	if len(values["user"]) != 1 || len(values["pass"]) != 1 {
		return "", ""
	}

	user := values["user"][0]
	pass := values["pass"][0]

	return user, pass
}

// ServeFile adds a file serving handler to
// the default mutex and logs to terminal.
func ServeFile(
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
