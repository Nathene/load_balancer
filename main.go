package main

import (
	"log"
	"net/http"

	lb "github.com/Nathene/load_balancer/pkg"
)

func main() {
	servers := []lb.Server{
		lb.NewSimpleServer("https://www.facebook.com"),
		lb.NewSimpleServer("https://www.google.com"),
		lb.NewSimpleServer("https://www.youtube.com"),
	}
	lb := lb.NewLoadBalancer("8000", servers)
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(w, r)
	}
	http.HandleFunc("/", handleRedirect)

	log.Printf("Serving requests from localhost:%s\n", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}
