package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

var (
	proxies map[string]*httputil.ReverseProxy
)

type Config struct {
	Host   string
	Server string
}

func main() {

	cfgs := []Config{
		Config{
			Host:   "blog.sketchground.dk",
			Server: "http://127.0.0.1:9900",
		},
		Config{
			Host:   "journal.sketchground.dk",
			Server: "http://127.0.0.1:9900",
		},
	}
	hosts := []string{}

	// Load services...
	proxies = map[string]*httputil.ReverseProxy{}
	for _, cfg := range cfgs {
		u, _ := url.Parse(cfg.Server)
		log.Printf("Initializing proxy connection for %v -> %v\n", cfg.Host, cfg.Server)
		proxies[cfg.Host] = httputil.NewSingleHostReverseProxy(u)
		hosts = append(hosts, cfg.Host)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Println("Starting reverse proxy for ssl connections")
		log.Fatal(http.Serve(autocert.NewListener(hosts...), &P{secure: true}))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Println("Starting reverse proxy for http connections")
		log.Fatal(http.ListenAndServe(":80", &P{})) // port 80
		wg.Done()
	}()
	wg.Wait()
}

type P struct {
	secure bool
}

func (p *P) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, ok := proxies[req.Host]
	if p.secure && ok {
		proxies[req.Host].ServeHTTP(rw, req)
		return
	}
	if !p.secure && ok { // Redirect http connections to https variant
		u := fmt.Sprintf("https://%v%v", req.Host, req.URL.Path)
		http.Redirect(rw, req, u, http.StatusFound)
		return
	}
	fmt.Fprintf(rw, "Nothing here. Go elsewhere.")
}
