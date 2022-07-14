package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	hbot "github.com/neurosnap/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	keywords = []string{"erock", "pico.sh", "picosh"}
	dms      = []string{"erock", "#pico.sh"}
	deny     = []string{"erock", "SaslServ", "NickServ"}
)

func resetTimer() time.Time {
	return time.Now().Add(
		time.Minute * time.Duration(10),
	)
}

func send(auth sasl.Client, subject string, body string) {
	to := []string{"irc@erock.io"}
	content := "From: bot@erock.io\r\n" +
		"To: irc@erock.io\r\n" +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		fmt.Sprintf("%s\r\n", body)
	msg := strings.NewReader(content)
	err := smtp.SendMail(
		"smtp.fastmail.com:587",
		auth,
		"bot@erock.io",
		to,
		msg,
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	ircPass := os.Getenv("IRC_PASS")
	smtPass := os.Getenv("IRC_SMTP_PASS")
	auth := sasl.NewPlainClient("", "me@erock.io", smtPass)

	timer := resetTimer()
	isAway := false

	saslOption := func(bot *hbot.Bot) {
		bot.SASL = true
		bot.SSL = true
		bot.Password = ircPass
	}

	bot, err := hbot.NewBot("irc.erock.io:6697", "erock/irc.libera.chat@bot", saslOption)
	if err != nil {
		panic(err)
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
			subject := fmt.Sprintf("%s - irc bot", m.From)

			channel := m.From
			if strings.Contains(m.To, "#") {
				channel = m.To
			}

			body := fmt.Sprintf(
				"%s\r\n---\r\nirc://irc.libera.chat/%s\r\nfrom: %s\r\nto: %s",
				m.Content,
				channel,
				m.From,
				m.To,
			)
			send(auth, subject, body)
			return false
		},
	}

	away := hbot.Trigger{
		Condition: func(b *hbot.Bot, m *hbot.Message) bool {
			return m.From == "erock"
		},
		Action: func(b *hbot.Bot, m *hbot.Message) bool {
			timer = resetTimer()
			if isAway {
				bot.Info("MARKING USER AS ACTIVE")
				// this removes the away status
				// b.Send("AWAY")
				isAway = false
			}
			return false
		},
	}

	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	bot.Logger.SetHandler(logHandler)
	bot.AddTrigger(notify)
	bot.AddTrigger(away)
	bot.Run()
	fmt.Println("Bot shutting down.")
}
