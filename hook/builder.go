package main

type builder struct {
	config *repoConfig
}

func (b *builder) build(branch string) {

	// load config
	//b.repos = load()

	// if folder does not exist, run git clone

	// checkout branch

	// run git pull

	// start timer

	// execute shell script that builds repo ( in goroutine ), read from json config

	// if return code is zero, update config with elapsed time
}

func newBuilder(config *repoConfig) *builder {
	return &builder{config}
}
