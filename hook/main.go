package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	configuration *config
	port          = flag.String("p", "666", "listen port")
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func handleHook(w http.ResponseWriter, r *http.Request) {
	var p pushPayload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Github payload: %s", p.Ref)

	// check if we have config
	if ok, repo := configuration.get(p.Repository.Name); ok {

		build := newBuilder(repo)
		go build.build(p.Repository.Fullname)
	}

}

func main() {
	flag.Parse()
	conf := &config{}
	conf.load()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/api/", handler)
	http.HandleFunc("/payload", handleHook)
	log.Printf("Turbine spinning at port %s", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
