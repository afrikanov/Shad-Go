// +build !solution

package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	for _, request := range args {
		resp, err := http.Get(request)
		if err != nil {
			os.Exit(1)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		_, _ = os.Stdout.Write(body)
	}
}
