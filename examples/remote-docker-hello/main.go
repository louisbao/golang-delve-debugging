package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = "8080"

func main() {
	http.HandleFunc("/hello", hello)

	fmt.Println("hello server is runing on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	hostName, _ := os.Hostname()
	fmt.Fprintf(w, "HostName: %s", hostName)
}
