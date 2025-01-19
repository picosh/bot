package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	hbot "github.com/neurosnap/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	awayNick   = "erock"
	fromEmail  = "bot@erock.io"
	toEmail    = "irc@erock.io"
	emailLogin = "me@erock.io"
	smtpAddr   = "smtp.fastmail.com:587"
	botNick    = "erock/irc.libera.chat@bot-tmp"
	ircHost    = "irc.pico.sh:6697"
	keywords   = []string{"erock", "pico.sh", "picosh"}
	dms        = []string{"erock", "#pico.sh", "#pico.sh-ops", "#pico.sh-+"}
	deny       = []string{"erock", "SaslServ", "NickServ", "ChanServ", "irc.pico.sh"}
)

func resetTimer() time.Time {
	return time.Now().Add(
		time.Minute * time.Duration(10),
	)
}

func send(auth sasl.Client, subject string, body string) {
	to := []string{toEmail}
	content := fmt.Sprintf("From: %s\r\n", fromEmail) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		fmt.Sprintf("%s\r\n", body)
	msg := strings.NewReader(content)
	err := smtp.SendMail(
		smtpAddr,
		auth,
		fromEmail,
		to,
		msg,
	)
	if err != nil {
		fmt.Println(err)
	}
}

func msgToEmail(m hbot.Message) string {
	channel := m.From
	if strings.Contains(m.To, "#") {
		channel = m.To
	}

	body := fmt.Sprintf(
		"%s\r\n---\r\n%s\r\nfrom: %s\r\nto: %s",
		m.Content,
		channel,
		m.From,
		m.To,
	)

	return body
}

func main() {
	ircSecret := os.Getenv("IRC_SECRET")
	ircPass := os.Getenv("IRC_PASS")
	smtPass := os.Getenv("IRC_SMTP_PASS")
	auth := sasl.NewPlainClient("", emailLogin, smtPass)

	timer := resetTimer()
	isAway := false
	queue := []hbot.Message{}

	saslOption := func(bot *hbot.Bot) {
		bot.SASL = true
		bot.SSL = true
		bot.Password = ircPass
	}

	bot, err := hbot.NewBot(ircHost, botNick, saslOption)
	if err != nil {
		bot.Error(err.Error())
		return
	}
	// remove default channels from bot since I'm connecting to a bouncer
	bot.Channels = []string{}
	// extend ping timeout so quiet irc setups don't keep disconnecting
	bot.PingTimeout = 4 * time.Hour

	go func() {
		for {
			now := time.Now()
			if !isAway && now.After(timer) {
				bot.Info("MARKING USER AS AWAY")
				isAway = true
				// bot.Send("AWAY idle")
			}
			time.Sleep(15 * time.Second)
		}
	}()

	go func() {
		for {
			if len(queue) > 0 {
				subject := fmt.Sprintf("%d messages -- irc bot", len(queue))
				body := ""
				for _, m := range queue {
					body += fmt.Sprintf("%s\r\n\r\n", msgToEmail(m))
				}
				send(auth, subject, body)
				// reset queue
				queue = make([]hbot.Message, 0)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	notify := hbot.Trigger{
		Condition: func(b *hbot.Bot, m *hbot.Message) bool {
			mentioned := false
			for _, key := range keywords {
				if strings.Contains(m.Content, strings.TrimSpace(key)) {
					mentioned = true
					break
				}
			}

			dmed := false
			for _, to := range dms {
				if m.To == to {
					dmed = true
					break
				}
			}

			denied := false
			for _, d := range deny {
				if m.From == d {
					denied = true
					break
				}
			}

			deniedContent := false
			for _, to := range dms {
				// these are weird edge cases where the content of the message is identical
				// to the `to` which means it was a JOIN event or
				// `from` which is a weird event that triggers when I click on a DM
				if m.Content == to || m.Content == m.From {
					deniedContent = true
					break
				}
			}

			return isAway && !denied && !deniedContent && (mentioned || dmed)
		},
		Action: func(b *hbot.Bot, m *hbot.Message) bool {
			bot.Info(fmt.Sprintf("NOTIFY ACTION FROM (%s) TO (%s)", m.From, m.To))
			queue = append(queue, *m)
			return false
		},
	}

	away := hbot.Trigger{
		Condition: func(b *hbot.Bot, m *hbot.Message) bool {
			return m.From == awayNick
		},
		Action: func(b *hbot.Bot, m *hbot.Message) bool {
			timer = resetTimer()
			// we're looking at the messages so clear the queue
			queue = make([]hbot.Message, 0)
			if isAway {
				bot.Info("MARKING USER AS ACTIVE")
				// this removes the away status
				// b.Send("AWAY")
				isAway = false
			}
			return false
		},
	}

	// Start an http server that sends a message to a user based on the body
	http.HandleFunc("POST /send", func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if auth != ircSecret {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		to := r.FormValue("to")
		message := r.FormValue("message")
		if to == "" || message == "" {
			http.Error(w, "missing required params", http.StatusBadRequest)
			return
		}

		bot.Msg(to, message)
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		http.ListenAndServe(":8080", nil)
		cancel()
	}()

	go func() {
		logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
		bot.Logger.SetHandler(logHandler)
		bot.AddTrigger(notify)
		bot.AddTrigger(away)
		bot.Run()
		cancel()
	}()

	<-ctx.Done()

	// send email when bot shuts down which could mean our bouncer is shutdown
	send(auth, "irc bot shutdown", "irc bot shutdown!")
	fmt.Println("Bot shutting down.")
}
