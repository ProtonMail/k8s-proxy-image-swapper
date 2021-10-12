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
	"gopkg.in/yaml.v2"
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
	_, err = w.Write(mutated)
	if err != nil {
		log.Println("Failed writing HTTP response")
	}
}

type Config struct {
	TLSCertPath string `yaml:"tlscertpath"`
	TLSKeyPath  string `yaml:"tlskeypath"`
	Port        string `yaml:"port"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Errorf("Usage : %v $CONFIG_FILE_PATH", os.Args[0])
		return
	}
	configFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Errorf("Error opening %v : %v", os.Args[1], err)
		return
	}
	defer configFile.Close()

	var config Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Errorf("Error reading config : %v", err)
		return
	}

	log.Println("Starting server ...")

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/mutate", handleMutation)

	s := &http.Server{
		Addr:           config.Port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1_048_576
	}

	log.Fatal(s.ListenAndServeTLS(config.TLSCertPath, config.TLSKeyPath))
}
