package notification

import (
 "github.com/bluele/slack"
)

var SlackNotificationHookURL string
var SlackNotificationChannel string

func SendMessage (title string, message string) {

  hook := slack.NewWebHook(SlackNotificationHookURL)
  err := hook.PostMessage(&slack.WebHookPostPayload{
    Text: title,
    Channel: SlackNotificationChannel,
    Attachments: []*slack.Attachment{
      {Text: message, Color: "good", MarkdownIn: []string{"pretext", "text", "fields"}},
    },
  })
  if err != nil {
    panic(err)
  }
}
