package main

type repoConfig struct {
	Name      string `json:"name"`
	BuildTime int    `json:"buildTime"`
}

type config struct {
	repos []repoConfig
}

func (c *config) load() {
	load("config.json", &c.repos)
}

func (c *config) get(name string) (bool, *repoConfig) {
	for _, r := range c.repos {
		if r.Name == name {
			return true, &r
		}
	}
	return false, nil
}
