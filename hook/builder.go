package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type builder struct {
	config     *repoConfig
	buildStart time.Time
	runner     *exec.Cmd
}

func (b *builder) fetchCode() error {
	// if folder does not exist, run git clone
	home, err := os.UserHomeDir()
	publicKeys, err := ssh.NewPublicKeysFromFile("git", fmt.Sprintf("%s/.ssh/id_rsa", home), "")
	if err != nil {
		log.Printf("generate publickeys failed: %s\n", err.Error())
	}

	path := fmt.Sprintf("./repos/%s", b.config.Name)
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      fmt.Sprintf("git@github.com:%s.git", b.config.Name),
		Auth:     publicKeys,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Printf("Failed to clone %s", err)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("and folder does not exist, aborting.")
			return err
		}
	}
	// checkout correct branch
	w, _ := repo.Worktree()
	ref := plumbing.NewBranchReferenceName(b.config.Branch)

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

func (b *builder) run() {
	log.Printf("Slave started %s branch: %s", b.config.Name, b.config.Branch)
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
	return &builder{config, time.Now(), nil}
}
