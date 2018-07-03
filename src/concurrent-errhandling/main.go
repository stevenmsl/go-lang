/*
Separate the concerns of error handling from the producer goroutines as what they can do with the errors are quite limited.
Let the goroutine that spawned the producer goroutines make the decisions in how to deal with the errors.
*/

package main

import (
	"fmt"
	"net/http"
)

//Result ...
//encompass both the response and the error
type Result struct {
	Error    error
	Response *http.Response
}

func main() {
	runCheckStatus3()
	//runCheckStatus2()
	//runCheckStatus()
}

//Separate the concerns of error handing from the producer goroutine
func checkStatus2(
	done <-chan interface{},
	urls ...string,
) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)
		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)
			result = Result{Error: err, Response: resp}
			select {
			case <-done:
				return
			case results <- result:
			}

		}
	}()
	return results
}

func runCheckStatus3() {
	done := make(chan interface{})
	defer close(done)
	errCount := 0
	urls := []string{"https://www.google.com", "b", "c", "d"}
	for result := range checkStatus2(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

//You can make better decisions dealing with errors here.
func runCheckStatus2() {
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus2(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
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
