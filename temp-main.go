package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		bytes, _ := ioutil.ReadFile("/Users/tahir/workspace/workhorse-logs/test-app/55027f17-7546-438c-90da-8dcc2c39e7ed/logs.txt")
		fmt.Fprintf(w, string(bytes))
	})

	http.ListenAndServe(":8081", nil)
}
