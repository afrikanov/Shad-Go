// +build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	resultTime float64
	wg         sync.WaitGroup
)

func getRequest(request string) {
	startTime := time.Now()
	resp, err := http.Get(request)
	requestTime := time.Since(startTime)
	if err == nil {
		fmt.Print("Request = ", request, " ")
		fmt.Print(resp.Status, " ")
		fmt.Println(requestTime.Seconds(), " seconds")
		resultTime += requestTime.Seconds()
	} else {
		fmt.Println("Get ", request, " failed")
	}
	wg.Done()
}

func main() {
	args := os.Args[1:]
	for _, request := range args {
		go getRequest(request)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Println(resultTime, " elapsed")
}
