package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	hbot "github.com/neurosnap/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

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
		resp.Write([]byte("ok"))
		fmt.Println("received request")
		bot.Msg("erock", "git commit")
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
