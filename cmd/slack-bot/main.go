package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nlopes/slack"
)

const (
	botName = "slack-bot"
)

func doCommand(command string, option ...string) (string, error) {
	cmd := exec.Command(command, option...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func sendReply(rtm *slack.RTM, channel string, out string, err error) {
	if err != nil {
		out += fmt.Sprintf("%s", err)
	}
	rtm.SendMessage(rtm.NewOutgoingMessage(out, channel))
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Println("SLACK_TOKEN environment variable is empty. please set token.")
		os.Exit(1)
	}

	api := slack.New(token)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// nop

		case *slack.ConnectedEvent:
			log.Println("connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			user := rtm.GetInfo().GetUserByID(ev.Msg.User)
			if user == nil || user.Name == botName {
				// nop
				continue
			}

			text := ev.Msg.Text
			commands := strings.Split(text, " ")
			var (
				out = ""
				err error
			)
			switch commands[0] {
			case "do":
				sendReply(rtm, ev.Channel, commands[1]+" command execution acknowledged...", nil)
				out, err = doCommand(commands[1], commands[2:]...)
			default:
				out, err = "unknown operation...", nil
			}
			// TODO: error check

			if len(out) < 4000 {
				sendReply(rtm, ev.Channel, out, err)
			} else {
				params := slack.FileUploadParameters{
					Channels: []string{ev.Channel},
					Content:  out,
				}
				_, e := api.UploadFile(params)
				if e != nil {
					log.Printf("failed to upload file: %s", e.Error())
				}
			}

		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials")
			os.Exit(1)
		default:
			// nop
		}
	}
}
