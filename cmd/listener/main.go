package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	listener "github.com/vincent-pli/tekton-listener/pkg/listener"
)

func handler(w http.ResponseWriter, r *http.Request) {
	event := map[string]interface{}{}

	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &event); err != nil {
		return
	}

	l, err := listener.New()
	if err != nil {
		fmt.Println("Err, Err")
	}
	l.HandleEvent(event)
}

func main() {
	log.Print("Hello world sample started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
