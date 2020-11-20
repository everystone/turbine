package main

type repoConfig struct {
	Name      string `json:"name"`
	Branch    string `json:"branch"`
	Build     string `json:"build"`
	Run       string `json:"run"`
	BuildTime int    `json:"buildTime"`
}

type config struct {
	repos []repoConfig
}

func (c *config) load() {
	load("config.json", &c.repos)

	// add default master branch if no branch specified
	for i, r := range c.repos {
		if r.Branch == "" {
			c.repos[i].Branch = "master"
		}
	}
}

func (c *config) save() {
	save("config.json", c.repos)
}

func (c *config) get(name string, branch string) (bool, *repoConfig) {
	for _, r := range c.repos {
		if r.Name == name && r.Branch == branch {
			return true, &r
		}
	}
	return false, nil
}
