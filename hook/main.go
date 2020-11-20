package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	log.Printf("Github payload: %v", p)
	ref := strings.Split(p.Ref, "/")
	branch := ref[len(ref)-1]
	log.Printf("Webhook event %s branch: %s", p.Repository.Fullname, branch)
	// check if we have config for repo / branch
	if ok, repo := configuration.get(p.Repository.Fullname, branch); ok {
		build := newBuilder(repo)
		go build.run(p.Ref)
	} else {
		log.Printf("No config matches.")
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
	log.Printf("Github hook: /payload")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
