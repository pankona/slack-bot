package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/exec"

	"github.com/nlopes/slack"
)

const (
	botName = "slack-spirit"
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
	log.Printf("reply: %s\n", out)
}

func main() {
	logger, err := syslog.New(syslog.LOG_LOCAL0, botName)
	if err != nil {
		log.Println("failed to configure logger. exit.")
		os.Exit(1)
	}
	log.SetOutput(logger)

	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Println("SLACK_TOKEN environment variable is empty. please set token.")
		os.Exit(1)
	}
	userName := os.Getenv("HACK_SPIRIT_USERNAME")
	if token == "" {
		log.Println("HACK_SPIRIT_USERNAME environment variable is empty. please set token.")
		os.Exit(1)
	}
	password := os.Getenv("HACK_SPIRIT_PASSWORD")
	if token == "" {
		log.Println("HACK_SPIRIT_PASSWORD environment variable is empty. please set token.")
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
			var out = ""
			var err error
			switch text {
			case "shukkin":
				sendReply(rtm, ev.Channel, "shukkin acknowledged...", nil)
				out, err = doCommand("hack-spirit", "start_work", "-u", userName, "-p", password)
			case "taikin":
				sendReply(rtm, ev.Channel, "taikin acknowledged...", nil)
				out, err = doCommand("hack-spirit", "finish_work", "-u", userName, "-p", password)
			case "status":
				sendReply(rtm, ev.Channel, "status acknowledged...", nil)
				out, err = doCommand("hack-spirit", "work_status", "-u", userName, "-p", password)
			default:
				out, err = "unknown operation...", nil
			}
			sendReply(rtm, ev.Channel, out, err)

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
