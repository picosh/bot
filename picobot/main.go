package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	hbot "github.com/neurosnap/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

type RepoUpdate struct {
	Data struct {
		Webhook struct {
			Date       string `json:"date"`
			Event      string `json:"event"`
			Repository struct {
				Name string `json:"name"`
				Rev  struct {
					ID        string `json:"id"`
					ShortID   string `json:"shortId"`
					Message   string `json:"message"`
					Committer struct {
						Name string `json:"name"`
					} `json:"committer"`
					Author struct {
						Name string `json:"name"`
					} `json:"author"`
				} `json:"revparse_single"`
			} `json:"repository"`
		} `json:"webhook"`
	} `json:"data"`
}

func main() {
	ircPass := os.Getenv("IRC_PICO_PASS")
	port := os.Getenv("IRC_WEB_PORT")
	if port == "" {
		port = "80"
	}

	saslOption := func(bot *hbot.Bot) {
		bot.SASL = true
		bot.SSL = true
		bot.Password = ircPass
	}

	opts := func(bot *hbot.Bot) {
		bot.Channels = []string{"#pico.sh"}
		// extend ping timeout so quiet irc setups don't keep disconnecting
		bot.PingTimeout = 4 * time.Hour
	}
	bot, err := hbot.NewBot(
		"irc.erock.io:6697",
		"picobot/irc.libera.chat@picobot",
		saslOption,
		opts,
	)
	if err != nil {
		panic(err)
	}

	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	bot.Logger.SetHandler(logHandler)

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("pong"))
	})

	http.HandleFunc("/push", func(resp http.ResponseWriter, req *http.Request) {
		/*
			var j interface{}
			err = json.NewDecoder(req.Body).Decode(&j)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", j)
		*/

		var ru RepoUpdate
		err := json.NewDecoder(req.Body).Decode(&ru)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		resp.Write([]byte("sending message to channel"))
		url := fmt.Sprintf(
			"https://git.sr.ht/~erock/%s/commit/%s",
			ru.Data.Webhook.Repository.Name,
			ru.Data.Webhook.Repository.Rev.ID,
		)
		bot.Msg("erock", fmt.Sprintf(
			"[sr.ht] %s -- (%s) %s",
			url,
			ru.Data.Webhook.Repository.Rev.Committer.Name,
			ru.Data.Webhook.Repository.Rev.Message,
		))
	})

	go func() {
		fmt.Printf("Starting web server on on %s\n", port)
		err = http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic(err)
		}
	}()

	bot.Run()
	fmt.Println("Bot shutting down.")
}
