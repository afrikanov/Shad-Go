// +build !solution

package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	mapKeyURL map[string]string
	mapURLKey map[string]string
	randMax   = int(1e9)
)

type GetURL struct {
	URL string `json:"url"`
}

type ResponseURL struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var url GetURL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "JSON is invalid", 400)
		return
	}
	_ = r.Body.Close()
	if mapURLKey[url.URL] != "" {
		marshaled, _ := json.Marshal(ResponseURL{
			URL: url.URL,
			Key: mapURLKey[url.URL],
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(marshaled)
		return
	}
	rand.Seed(time.Now().UnixNano())
	key := rand.Intn(randMax)
	keyString := strconv.Itoa(key)
	mapKeyURL[keyString] = url.URL
	mapURLKey[url.URL] = keyString
	marshaled, _ := json.Marshal(ResponseURL{
		URL: url.URL,
		Key: keyString,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshaled)
}

func goHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.String()[4:]
	if mapKeyURL[key] == "" {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	url := mapKeyURL[key]
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	portPtr := flag.Int("port", 8000, "port string")
	flag.Parse()
	portNumber := *portPtr
	mapKeyURL = make(map[string]string)
	mapURLKey = make(map[string]string)

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/go/", goHandler)

	localAddress := "localhost:" + strconv.Itoa(portNumber)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}
