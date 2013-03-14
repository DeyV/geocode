package main

import (
	"fmt"
	"log"
	"net/http"
)

const listenAddr = "localhost:8100"

func main() {
	fmt.Println("Server run on: ", listenAddr)

	http.Handle("/json", http.HandlerFunc(codeHandler))
	http.Handle("/text", http.HandlerFunc(codeTextHandler))
	http.Handle("/textcache", http.HandlerFunc(codeTextCacheHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
