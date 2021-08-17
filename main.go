package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	message   string
	token     string
	header    string
	frequency int
)

func init() {
	flag.StringVar(&message, "m", "", "Bot message")
	flag.StringVar(&token, "t", "", "Bot token")
	flag.StringVar(&header, "h", "", "Nickname header")
	flag.IntVar(&frequency, "f", 1, "Frequency of updates, in seconds")
	flag.Parse()
}

func main() {

	// create discord connection
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("Creating Discord session: %s", err)
		return
	}

	// show as online
	err = dg.Open()
	if err != nil {
		fmt.Printf("Opening discord connection: %s", err)
		return
	}

	// get guides for bot
	guilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		fmt.Printf("Error etting guilds: %s", err)
		return
	}

	// split message and measure
	messages := strings.Split(message, "")
	cap := len(messages)

	// create timer
	timer := time.NewTicker(time.Duration(frequency) * time.Second)

	// forever
	for {
		// starting at each indix once
		for index, _ := range messages {
			// every period of our timer
			select {
			case <-timer.C:
				// add header and current index
				display := header + messages[index]
				counter := index + 1
				for next := 1; next < 31-len(header); next++ {
					// start message over from the beginning
					if counter == cap {
						display += " "
						counter = 0
					}

					// append the next character
					display += messages[counter]
					counter++
				}

				// Update nickname in guilds
				for _, g := range guilds {
					err = dg.GuildMemberNickname(g.ID, "@me", display)
					if err != nil {
						fmt.Printf("Error updating nickname: %s\n", err)
						continue
					}
					fmt.Printf("Updating nickname: %s in %s\n", display, g.Name)
				}
			}
		}
	}
}
