package gophish

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

// Build gathers the resources of index.raw.html and creates a new servable index.html
func Build(web *Website) {

	rescounter := 0

	for i := 0; i < len(web.RawFiles); i++ {

		filein, errin := os.Open(web.GetFile(web.RawFiles[i]))
		fileout, errout := os.Create(web.GetFile(
			strings.Replace(web.RawFiles[i], ".raw.", ".", 1)))

		defer filein.Close()
		defer fileout.Close()

		if errin != nil {
			log.Fatal(errin)
		}
		if errout != nil {
			log.Fatal(errout)
		}

		lineRe := regexp.MustCompile(web.LineMatchRe)
		srcRe := regexp.MustCompile(web.ResMatchRe)

		scanner := bufio.NewScanner(filein)
		for scanner.Scan() {
			line := scanner.Text()
			rescounter++

			if lineRe.MatchString(line) {

				src := srcRe.FindStringSubmatch(line)

				if len(src) > 0 {

					link := strings.Split(src[0], `"`)[1]
					_ext := strings.Split(link, `.`)
					ext := _ext[len(_ext)-1]

					respath := web.GetResource(fmt.Sprintf("%d.%s", rescounter, ext)) // path.Join(common.ResFolder, fmt.Sprintf("%d.%s", rescounter, ext))
					rawlink := link

					if !strings.Contains(link, "http://") && !strings.Contains(link, "http://") {
						rawlink = web.WebLink + "/" + link
					}

					log.Printf("Downloading: %s", rawlink)

					err := DownloadFile(respath, rawlink)
					if err != nil {
						println(err)
						log.Fatal(err)
					}

					line = strings.Replace(line, link, "/"+respath, -1)
					println(line)
				}
			}

			fmt.Fprintln(fileout, line)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}
}

// DownloadFile will download a url to a local file.
// Writes as it downloads and does not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	os.MkdirAll(path.Dir(filepath), os.ModePerm)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
