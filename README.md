# LINE BotTemplate

A minimal Golang LINE Bot template built on the official [line-bot-sdk-go v8](https://github.com/line/line-bot-sdk-go). Clone it, deploy it, and you have a working LINE Bot in a few minutes.

[![GoDoc](https://godoc.org/github.com/kkdai/LineBotTemplate.svg?status.svg)](https://godoc.org/github.com/kkdai/LineBotTemplate)
![Go](https://github.com/kkdai/LineBotTemplate/workflows/Go/badge.svg)
[![goreportcard.com](https://goreportcard.com/badge/github.com/kkdai/LineBotTemplate)](https://goreportcard.com/report/github.com/kkdai/LineBotTemplate)

## What it does

- Serves a `/callback` webhook endpoint and verifies the LINE signature
- Echoes back any text message it receives
- Replies with the sticker ID / resource type when it receives a sticker
- Logs every incoming event so you can see the payload shape while developing

That is the whole app — one file, [`main.go`](main.go), roughly 100 lines. It is meant to be read and then rewritten into your own bot.

## Requirements

| | |
|---|---|
| Go | 1.25 or later |
| SDK | `github.com/line/line-bot-sdk-go/v8` v8.22.0 |

## Environment variables

Both are required. You get them from the [LINE Developers Console](https://developers.line.biz/console/).

| Name | Where to find it |
|---|---|
| `ChannelSecret` | Messaging API channel → **Basic settings** tab |
| `ChannelAccessToken` | Messaging API channel → **Messaging API** tab → issue a long-lived token |

`PORT` is optional and defaults to `5000`. Most PaaS providers set it for you.

## Getting started

### 1. Create a LINE Messaging API channel

1. Sign in to the [LINE Developers Console](https://developers.line.biz/console/).
2. Create a new **Messaging API** channel.
3. Copy the `Channel Secret` from the **Basic settings** tab.
4. Issue a `Channel Access Token` on the **Messaging API** tab.
5. Open **LINE Official Account Manager** from the **Basic settings** tab, go to *Response settings*, and turn **Webhook** on (turn *Auto-reply messages* off, otherwise the OA replies before your bot does).

### 2. Deploy

Pick whichever platform you prefer.

#### Render

[![Deploy to Render](http://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)

Render reads [`render.yaml`](render.yaml). Fill in `ChannelSecret` and `ChannelAccessToken` when prompted.

#### Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Heroku reads [`app.json`](app.json) and [`Procfile`](Procfile). Fill in the same two variables, and note the app name — you need it for the webhook URL.

### 3. Point LINE at your webhook

Back in the LINE Developers Console, set the **Webhook URL** to your deployment plus `/callback`:

```
https://<your-app>.onrender.com/callback
https://<your-app>.herokuapp.com/callback
```

Hit **Verify**. If it returns success, add the bot as a friend and send it a message.

## Running locally

```bash
git clone https://github.com/kkdai/LineBotTemplate.git
cd LineBotTemplate

export ChannelSecret=<your channel secret>
export ChannelAccessToken=<your channel access token>

go run main.go
# http://localhost:5000/
```

LINE only calls HTTPS webhooks, so expose the local port with a tunnel:

```bash
ngrok http 5000
```

Then set the webhook URL to `https://<subdomain>.ngrok-free.app/callback`.

## Making it your own

All the interesting code lives in the event loop in `main.go`:

```go
for _, event := range cb.Events {
    switch e := event.(type) {
    case webhook.MessageEvent:
        switch message := e.Message.(type) {
        case webhook.TextMessageContent:
            // your logic here
        }
    }
}
```

- **Reply to different message types** — add cases for `webhook.ImageMessageContent`, `webhook.LocationMessageContent`, `webhook.AudioMessageContent`, and so on.
- **Reply with something richer than text** — swap `messaging_api.TextMessage` for `StickerMessage`, `ImageMessage`, `FlexMessage`, or `TemplateMessage` in the `Messages` slice.
- **Handle non-message events** — `webhook.FollowEvent`, `webhook.JoinEvent`, `webhook.PostbackEvent`, and `webhook.BeaconEvent` are cases on the *outer* switch, alongside `webhook.MessageEvent`.
- **Push instead of reply** — a reply token is single-use and short-lived; use `bot.PushMessage` when you need to message a user outside of a reply window.

See the [Messaging API reference](https://developers.line.biz/en/reference/messaging-api/) for the full surface.

## Tutorials

- Video: [How to deploy LINE BotTemplate](https://www.youtube.com/watch?v=0BIknEz1f8k)
- Video: [How to modify your LINE BotTemplate code](https://www.youtube.com/watch?v=ckij73sIRik)
- 中文教學：[用 Golang 打造你的 LINE Bot](http://www.evanlin.com/create-your-line-bot-golang/)

## Inspired by

- [Golang (heroku) で LINE Bot 作ってみる](http://qiita.com/dongri/items/ba150f04a98e96b160e7)
- [LINE BOT をとりあえずタダで Heroku で動かす](http://qiita.com/yuya_takeyama/items/0660a59d13e2cd0b2516)
- [阿美語萌典 BOT](https://github.com/miaoski/amis-linebot)

## Project52

It is one of my [project 52](https://github.com/kkdai/project52).

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
