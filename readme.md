# pkg - Simple Load Balancer

This package provides a basic implementation of a load balancer in Go. It uses a round-robin algorithm to distribute incoming HTTP requests across a set of backend servers.

## Key Components

* **`Server` Interface:**
   * Defines the contract for backend servers that the load balancer can work with.
   * Requires the following methods:
      * `Address() string`: Returns the server's address.
      * `IsAlive() bool`: Checks if the server is currently available.
      * `Serve(w http.ResponseWriter, r *http.Request)`: Handles an incoming request.

* **`simpleServer` Struct:**
   * A basic implementation of the `Server` interface.
   * Acts as a reverse proxy to forward requests to a single target server.
   * Constructor: `NewSimpleServer(addr string)`

* **`loadBalancer` Struct:**
   * The core load balancer component.
   * Maintains a list of backend servers (`servers`).
   * Uses `roundRobinCount` to keep track of the next server to be used.
   * Constructor: `NewLoadBalancer(port string, servers []Server)`
   * Key methods:
      * `getNextAvailableServer() Server`: Finds the next available server using round-robin.
      * `ServeProxy(w http.ResponseWriter, r *http.Request)`: Handles incoming requests, forwarding them to an available server.

## Usage Example

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	lb "github.com/Nathene/load_balancer/pkg" 
)

func main() {
	// Define backend servers
	servers := []lb.Server{
		lb.NewSimpleServer("https://www.facebook.com"),
		lb.NewSimpleServer("https://www.google.com"),
		lb.NewSimpleServer("https://www.youtube.com"),
	}

	// Create load balancer
	loadBalancer := lb.NewLoadBalancer("8000", servers)

	// Define a simple handler for the load balancer
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.ServeProxy(w, r)
	})

	// Start the load balancer
	fmt.Println("Load Balancer started at :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
```