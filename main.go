// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func main() {
	channelSecret := os.Getenv("ChannelSecret")
	bot, err := messaging_api.NewMessagingApiAPI(
		os.Getenv("ChannelAccessToken"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		log.Println("/callback called...")

		cb, err := webhook.ParseRequest(channelSecret, req)
		if err != nil {
			log.Printf("Cannot parse request: %+v\n", err)
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		log.Println("Handling events...")
		for _, event := range cb.Events {
			log.Printf("/callback called%+v...\n", event)

			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: message.Text,
								},
							},
						},
					); err != nil {
						log.Print(err)
					} else {
						log.Println("Sent text reply.")
					}
				case webhook.StickerMessageContent:
					replyMessage := fmt.Sprintf(
						"貼圖訊息: sticker id is %s, stickerResourceType is %s", message.StickerId, message.StickerResourceType)
					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: replyMessage,
								},
							},
						}); err != nil {
						log.Print(err)
					} else {
						log.Println("Sent sticker reply.")
					}
				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			case webhook.FollowEvent:
				log.Printf("Followed by %s\n", sourceID(e.Source))
			case webhook.MemberJoinedEvent:
				log.Printf("Members joined %s: %s\n", sourceID(e.Source), userIDs(orZero(e.Joined).Members))
			case webhook.MemberLeftEvent:
				log.Printf("Members left %s: %s\n", sourceID(e.Source), userIDs(orZero(e.Left).Members))
			case webhook.BeaconEvent:
				beacon := orZero(e.Beacon)
				log.Printf("Beacon %q from %s: hwid=%s\n", beacon.Type, sourceID(e.Source), beacon.Hwid)
			default:
				log.Printf("Unsupported event: %T\n", event)
			}
		}
	})

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("http://localhost:" + port + "/")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// sourceID describes where an event came from. The concrete type behind
// webhook.SourceInterface depends on the chat: 1:1 chats carry a UserSource,
// while group and multi-person chats carry a GroupSource or RoomSource, so this
// must never assume UserSource.
func sourceID(src webhook.SourceInterface) string {
	switch s := src.(type) {
	case webhook.UserSource:
		return "user " + s.UserId
	case webhook.GroupSource:
		return "group " + s.GroupId
	case webhook.RoomSource:
		return "room " + s.RoomId
	default:
		return fmt.Sprintf("unknown source (%T)", src)
	}
}

// orZero dereferences p, returning the zero value when it is nil. The webhook
// models optional sub-objects as pointers, so payloads that omit them must not
// crash the handler.
func orZero[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// userIDs collects the user IDs out of the member list carried by member
// joined/left events.
func userIDs(members []webhook.UserSource) []string {
	ids := make([]string, len(members))
	for i, m := range members {
		ids[i] = m.UserId
	}
	return ids
}
