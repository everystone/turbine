package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type builder struct {
	config     *repoConfig
	buildStart time.Time
	runner     *exec.Cmd
}

func (b *builder) fetchCode() error {
	// if folder does not exist, run git clone

	repo, err := git.PlainClone("/repos/", false, &git.CloneOptions{
		URL:      fmt.Sprintf("git@github.com:%s.git", b.config.Name),
		Progress: os.Stdout,
	})
	if err != nil {
		log.Printf("Failed to clone %s", err)
		if repo == nil {
			log.Printf("repo is nil, aborting.")
			return err
		}
		// repo exists, continue and checkout branch.
		log.Printf("repo exists")
	}
	// checkout correct branch
	w, _ := repo.Worktree()
	ref := plumbing.NewBranchReferenceName(fmt.Sprintf("refs/heads/%s", b.config.Branch))

	log.Printf("Checking out branch %s", ref)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: ref,
	})
	if err != nil {
		log.Printf("Failed to checkout branch %v", ref)
		return err
	}
	return nil

}

func (b *builder) build() error {
	// execute shell script that builds repo, read from json config
	cmd := exec.Command(b.config.Build)
	err := cmd.Run()
	t := time.Now()
	elapsed := t.Sub(b.buildStart)
	if err != nil {
		log.Printf("Build of %s completed in %v", b.config.Name, elapsed)
		b.config.BuildTime = int(elapsed.Seconds())
		configuration.save()
	}

	log.Printf("Build of %s failed after %v", b.config.Name, elapsed)
	log.Printf("error: %v", err)
	return err
}

func (b *builder) status() string {
	t := time.Now()
	elapsed := t.Sub(b.buildStart)
	// TODO: calculate percentage based on b.config.buildTime
	return fmt.Sprintf("Build has been running for %v", elapsed)
}

func (b *builder) deploy() {
	// execute command from config in new process
	// or as a managed child process that we can controll from api?
	// support rest api for status, stop, start, restart of service..

	b.runner = exec.Command(b.config.Run)
	log.Printf("Process started %s", b.config.Run)
	b.runner.Run()

}

func (b *builder) run(branch string) {

	err := b.fetchCode()
	if err != nil {
		return
	}
	err = b.build()
	if err != nil {
		return
	}
	b.deploy()

}

func newBuilder(config *repoConfig) *builder {
	return &builder{config, time.Now()}
}
