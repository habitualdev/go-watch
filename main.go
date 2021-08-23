package main

import (
	"flag"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

//Simple http POST function. Easily extendable to add auth headers

func sendData(url string, filename string) string {
	client := &http.Client{}
	var errortype string
	defer func() {
		if r:= recover(); r!=nil {
			log.Println("ERROR: " + errortype )

		}
	}()
	errortype = "Open File failed"
	data, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	errortype = "Request Creation Failed"
	req, err := http.NewRequest("POST", url, data)
	//req.Header.Set(STUFF HERE)
	//req.SetBasicAuth(USER,PASSWORD)
	///req.AddCookie(COOKIE)
	if err != nil {
		panic(err)
	}
	errortype = "POST failed"
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp.Status
}

func main() {
	var watchpath string
	var extension string
	var url string

	// Create Flag conditions

	flag.StringVar(&watchpath, "p", ".", "Set path to watch")
	flag.StringVar(&extension, "e", ".*", "Set the extension to watch")
	flag.StringVar(&url, "u", "http://127.0.0.1", "Set the url to post to")
	flag.Parse()

	// Create Watcher and apply watch filters

	w := watcher.New()
	w.FilterOps(watcher.Move, watcher.Create)
	r := regexp.MustCompile(extension)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				response := sendData(url,event.Path)
 				log.Println("{url: " + strconv.Quote(url) + ", response:" + strconv.Quote(response) + ", file: " + strconv.Quote(event.Path) + "}")
			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Add(watchpath); err != nil {
			log.Fatalln(err)
		}

	for path, f := range w.WatchedFiles() {
		log.Printf("%s: %s\n", path, f.Name())
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Println(err)
		}

	}
}
