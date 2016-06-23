package ship

import (
	"encoding/json"
	"log"
	"net/http"
)

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

var currentConfig *Config

func Serve(config *Config) error {
	currentConfig = config
	http.HandleFunc("/", reqHandler)
	return http.ListenAndServe("0.0.0.0:8080", Log(http.DefaultServeMux))
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RemoteAddr, r.Method)
		handler.ServeHTTP(w, r)
	})
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	if authReqestByToken(r) {
		decoder := json.NewDecoder(r.Body)
		var imgConfig IncomingMessage

		err := decoder.Decode(&imgConfig)
		if err != nil {
			http.Error(w, "Could not decode json", 500)
			log.Print(err)
			return
		}
		go handleIncomingMsg(imgConfig)
		return
	}
	http.Error(w, "Not Authorized", 401)
}

func authReqestByToken(r *http.Request) bool {
	key := r.URL.Query().Get("token")
	return currentConfig.Serve.Token == "" || key == currentConfig.Serve.Token
}

func handleIncomingMsg(img IncomingMessage) {
	//arukas api call

}
