package main

import (
	"fmt"
	"net/http"
)

func main() {
	runCheckStatus()
}

func checkStatus(
	done <-chan interface{},
	urls ...string,
) <-chan *http.Response {
	responses := make(chan *http.Response)
	go func() {
		defer close(responses)
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err) //not much you can do other than print the error
				continue
			}
			select {
			case <-done:
				return
			case responses <- resp:
			}

		}
	}()
	return responses
}

func runCheckStatus() {
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}
