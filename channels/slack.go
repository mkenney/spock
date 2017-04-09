package channels

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	slack "github.com/kcmerrill/slack-go-webhook"
	"github.com/kcmerrill/spock/info"
)

// Slack sends a messages as an incoming webhook to the slack api
func Slack(stdin io.Reader, args []string) (string, error) {
	// lets get started ...
	var webhook, channel string

	// our flags
	f := flag.NewFlagSet("slack", flag.ContinueOnError)
	f.StringVar(&webhook, "webhook", "", "The integration endpoint")
	f.StringVar(&channel, "channel", "", "The channel to be used")

	// set flags
	f.Parse(args)

	if webhook != "" {
		// yay! we have a webhook!
		in, _ := ioutil.ReadAll(stdin)
		cInfo := info.New(in)

		text := "[" + cInfo.ID + "] failed"
		if cInfo.Template != "" {
			text = cInfo.Template
		}

		attachment := slack.Attachment{Color: "danger"}
		attachment.
			AddField(slack.Field{Title: "Check", Value: cInfo.ID}).
			AddField(slack.Field{Title: "Error", Value: cInfo.Error})
		payload := slack.Payload{
			Text:        text,
			Attachments: []slack.Attachment{attachment},
		}
		err := slack.Send(webhook, "", payload)
		if len(err) > 0 {
			return err[0].Error(), nil
		}

		return "Notified failure via Slack", nil
	}

	return "", fmt.Errorf("Param 'webhook' needs to be set for the 'slack' integration")
}
