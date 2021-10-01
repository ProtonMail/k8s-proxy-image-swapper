// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	m "github.com/Polyconseil/k8s-proxy-image-swapper/mutate"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
}

func handleMutation(w http.ResponseWriter, r *http.Request) {
	debug := os.Getenv("LOGLEVEL") == "DEBUG"
	if debug {
		fmt.Println("handleMutation called")
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}

	registry := os.Getenv("REGISTRY_URL")

	mutated, err := m.Mutate(body, debug, registry)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func main() {
	log.Println("Starting server ...")

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/mutate", handleMutation)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Fatal(s.ListenAndServeTLS("/tls/cert.pem", "/tls/key.pem"))
}
