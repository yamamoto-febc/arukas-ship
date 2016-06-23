package message

type IncomingMessage struct {
	Repository struct {
		Status    string
		RepoUrl   string `json:"repo_url"`
		Owner     string
		IsPrivate bool `json:"is_private"`
		Name      string
		StarCount int    `json:"star_count"`
		RepoName  string `json:"repo_name"`
	}

	Push_data struct {
		PushedAt int `json:"pushed_at"`
		Images   []string
		Pusher   string
	}
}
