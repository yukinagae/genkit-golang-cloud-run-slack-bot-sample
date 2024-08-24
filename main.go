package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/yukinagae/genkit-golang-cloud-run-slack-bot-sample/flow"
)

func main() {
	ctx := context.Background()

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	slackSigningSecret := os.Getenv("SLACK_SIGNING_SECRET")

	var api = slack.New(slackBotToken)

	f := flow.DefineFlow(ctx)

	// Define the http handler for the Slack Events API
	// see: https://github.com/slack-go/slack/blob/master/examples/eventsapi/events.go
	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slack.NewSecretsVerifier(r.Header, slackSigningSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			if _, err := w.Write([]byte(r.Challenge)); err != nil {
				log.Fatal(err)
			}
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				// TODO: handle retry

				// skip if bot mentions itself
				botID := ev.BotID
				if botID != "" {
					return
				}

				// thinking...
				ts := getTimestamp(ev)
				_, botMessageTimestamp, err := api.PostMessage(
					ev.Channel,                              //
					slack.MsgOptionText("typing...", false), //
					slack.MsgOptionTS(ts),                   //
				)
				if err != nil {
					log.Fatal(err)
					return
				}
				// skip if failed to send message
				if botMessageTimestamp == "" {
					return
				}

				// delete mention
				input := deleteMention(ev.Text)

				// run the flow to get the answer
				answer, err := f.Run(ctx, input)
				if err != nil {
					log.Fatal(err)
					return
				}

				log.Printf("💖answer: %v", answer) // TODO: debug

				if _, _, _, err := api.UpdateMessage(
					ev.Channel,                         //
					botMessageTimestamp,                //
					slack.MsgOptionText(answer, false), //
				); err != nil {
					log.Fatal(err)
				}
			}
		}
	})
	log.Println("Server listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func getTimestamp(event *slackevents.AppMentionEvent) string {
	if event.ThreadTimeStamp != "" {
		return event.ThreadTimeStamp
	}
	return event.TimeStamp
}

func deleteMention(rawInput string) string {
	re := regexp.MustCompile(`<@.*?>`)
	input := re.ReplaceAllString(rawInput, "")
	input = strings.TrimSpace(input)
	return input
}
