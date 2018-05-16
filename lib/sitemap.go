package vegeta

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func SitemapToFile(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	r := regexp.MustCompile(`<loc>(.+)?<\/loc>`)

	res := r.FindAllString(string(html), -1)

	fname := os.TempDir() + "targets.txt"
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, v := range res {
		f.WriteString(prepareUrl(v))
	}

	return fname
}

func prepareUrl(s string) string {
	rawUrl := strings.NewReplacer(
		`<loc>`, "",
		`</loc>`, "",
		`%`, "%25",
	).Replace(
		s,
	)
	tmp := strings.Split(rawUrl, "://")
	if len(tmp) < 2 {
		log.Fatalf("Incorrect url - %s", rawUrl)
	}

	// u, err := url.Parse(url.QueryEscape(tmp[1]))
	// if err != nil {
	// 	panic(err)
	// }

	return "GET " + rawUrl + "\r\n"
}
