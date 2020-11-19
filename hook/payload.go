package main

type repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Fullname string `json:"full_name"`
	Owner    user   `json:"owner"`
}

type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type pushPayload struct {
	Ref        string     `json:"ref"`
	Sender     user       `json:"sender"`
	Pusher     user       `json:"pusher"`
	Repository repository `json:"repository"`
}
