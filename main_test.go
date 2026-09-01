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
	"slices"
	"testing"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// memberJoinedPayload is a member joined event as the LINE platform sends it.
// Note the source: member joined/left only happen in group and multi-person
// chats, so the source is never a UserSource.
const memberJoinedPayload = `{
  "type": "memberJoined",
  "timestamp": 1462629479859,
  "mode": "active",
  "webhookEventId": "01FZ74A0TDDPYRVKNK77XKC3ZR",
  "deliveryContext": { "isRedelivery": false },
  "source": { "type": "group", "groupId": "Ca56f94637c...", "userId": "U4af4980629..." },
  "replyToken": "0f3779fba3b349968c5d07db31eab56f",
  "joined": {
    "members": [
      { "type": "user", "userId": "U4af4980629..." },
      { "type": "user", "userId": "U91eeaf62d9..." }
    ]
  }
}`

// TestMemberJoinedIsEventNotMessageContent pins down the two facts the previous
// handler got wrong: a memberJoined payload decodes to an *event*, so it can
// only be matched on the outer switch over cb.Events, and its source is a
// GroupSource, so asserting webhook.UserSource on it panics.
func TestMemberJoinedIsEventNotMessageContent(t *testing.T) {
	event, err := webhook.UnmarshalEvent([]byte(memberJoinedPayload))
	if err != nil {
		t.Fatalf("UnmarshalEvent: %v", err)
	}

	e, ok := event.(webhook.MemberJoinedEvent)
	if !ok {
		t.Fatalf("got %T, want webhook.MemberJoinedEvent", event)
	}
	if _, isMessage := event.(webhook.MessageEvent); isMessage {
		t.Error("memberJoined decoded as a MessageEvent; it can never appear in the e.Message switch")
	}
	if _, isUser := e.Source.(webhook.UserSource); isUser {
		t.Error("source decoded as UserSource; the group case would not be exercised")
	}

	if got, want := sourceID(e.Source), "group Ca56f94637c..."; got != want {
		t.Errorf("sourceID = %q, want %q", got, want)
	}
	want := []string{"U4af4980629...", "U91eeaf62d9..."}
	if got := userIDs(orZero(e.Joined).Members); !slices.Equal(got, want) {
		t.Errorf("userIDs = %v, want %v", got, want)
	}
}

func TestSourceID(t *testing.T) {
	tests := []struct {
		name string
		src  webhook.SourceInterface
		want string
	}{
		{"user", webhook.UserSource{UserId: "U1"}, "user U1"},
		{"group", webhook.GroupSource{GroupId: "C1", UserId: "U1"}, "group C1"},
		{"room", webhook.RoomSource{RoomId: "R1", UserId: "U1"}, "room R1"},
		{"nil", nil, "unknown source (<nil>)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sourceID(tt.src); got != tt.want {
				t.Errorf("sourceID = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUserIDs(t *testing.T) {
	got := userIDs([]webhook.UserSource{{UserId: "U1"}, {UserId: "U2"}})
	if want := []string{"U1", "U2"}; !slices.Equal(got, want) {
		t.Errorf("userIDs = %v, want %v", got, want)
	}
	if got := userIDs(nil); len(got) != 0 {
		t.Errorf("userIDs(nil) = %v, want empty", got)
	}
}

// TestOrZeroNil covers the payloads that omit an optional sub-object: the
// handler must log them, not panic.
func TestOrZeroNil(t *testing.T) {
	if got := userIDs(orZero((*webhook.JoinedMembers)(nil)).Members); len(got) != 0 {
		t.Errorf("userIDs = %v, want empty", got)
	}
	if got := orZero((*webhook.BeaconContent)(nil)).Hwid; got != "" {
		t.Errorf("Hwid = %q, want empty", got)
	}
}
