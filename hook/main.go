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
	slaves        map[string]*builder
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
	if p.Zen != "" {
		log.Printf("Ping payload detected, skipping.")
		fmt.Fprintf(w, "hello friend")
		return
	}
	log.Printf("Github payload: %v", p)
	ref := strings.Split(p.Ref, "/")
	branch := ref[len(ref)-1]
	log.Printf("Webhook event %s branch: %s", p.Repository.Fullname, branch)
	// check if we have config for repo / branch
	if ok, repo := configuration.get(p.Repository.Fullname, branch); ok {
		// we have config, check if we have an existing slave/builder
		if slave, ok := slaves[repo.Name]; ok {

			log.Printf("Killing existing slave..")
			// todo: let it finish if building?
			slave.runner.Process.Kill()
			slave = newBuilder(repo)
			go slave.run()
			return
		}
		build := newBuilder(repo)
		slaves[repo.Name] = build
		go build.run()
	} else {
		log.Printf("No config matches: %s  branch: %s", p.Repository.Fullname, branch)
	}
	fmt.Fprintf(w, "build started")
}

func main() {
	flag.Parse()
	configuration = &config{}
	configuration.load()
	slaves = make(map[string]*builder)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/api/", handler)
	http.HandleFunc("/payload", handleHook)
	log.Printf("Turbine spinning at port %s", *port)
	log.Printf("Github hook: /payload")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
