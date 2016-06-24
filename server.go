package ship

import (
	"encoding/json"
	"fmt"
	"github.com/yamamoto-febc/arukas-ship/arukas"
	"github.com/yamamoto-febc/arukas-ship/message"
	"log"
	"net/http"
)

var currentConfig *Config

func Serve(config *Config) error {
	currentConfig = config
	http.HandleFunc("/", reqHandler)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Serve.Port), Log(http.DefaultServeMux))
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
		var imgConfig message.IncomingMessage

		err := decoder.Decode(&imgConfig)
		if err != nil {
			http.Error(w, "Could not decode json", 500)
			log.Print(err)
			return
		}

		appName := r.URL.Query().Get("app")
		if appName == "" {
			http.Error(w, "Could not read AppName. Please set 'app' parameter.", 500)
			log.Print("URL Parameter 'app' is empty.")
			return
		}

		go handleIncomingMsg(appName, imgConfig)
		return
	}
	http.Error(w, "Not Authorized", 401)
}

func authReqestByToken(r *http.Request) bool {
	key := r.URL.Query().Get("token")
	return currentConfig.Serve.Token == "" || key == currentConfig.Serve.Token
}

func handleIncomingMsg(appName string, img message.IncomingMessage) {

	client, err := arukas.NewArukasClient()
	if err != nil {
		log.Fatal(err)
	}
	err = client.HandleRequest(appName, &img)
	if err != nil {
		log.Fatal(err)
	}

}
