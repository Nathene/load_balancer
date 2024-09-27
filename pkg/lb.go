package pkg

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func (ss *simpleServer) Address() string {
	return ss.addr
}

func (ss *simpleServer) IsAlive() bool {
	return true
}

func (ss *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	ss.proxy.ServeHTTP(w, r)
}

type loadBalancer struct {
	Port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *loadBalancer {
	return &loadBalancer{
		Port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func (lb *loadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *loadBalancer) ServeProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(w, r)

}

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

func NewSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}
