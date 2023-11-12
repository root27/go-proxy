package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
)

var Servers = []string{
	"localhost:8081",
	"localhost:8082",
	"localhost:8083",
}

var customTransport = http.DefaultTransport

func chooseServer() string {

	r := rand.New(rand.NewSource(99))

	randomNumber := r.Intn(len(Servers))

	return Servers[randomNumber]
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Choose server
	server := chooseServer()

	// Create request
	req, err := http.NewRequest(r.Method, "http://"+server+r.URL.String(), r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Copy headers
	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}

	// Make request
	resp, err := customTransport.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	// Copy headers
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy body

	io.Copy(w, resp.Body)

	log.Print("Proxying request to " + server)
	log.Print("Request: " + r.URL.String())
	log.Print("Response: " + resp.Status)

}

func main() {

	http.HandleFunc("/", proxyHandler)

	log.Println("Starting proxy server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
