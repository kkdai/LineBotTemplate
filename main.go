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
// https://github.com/line/line-bot-sdk-go/tree/master/linebot

package main

import (
	"strconv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var silent bool = false
var tellTimeInterval int = 15
var echoMap = make(map[string]bool)

var loc, _ = time.LoadLocation("Asia/Taipei")
var	now = time.Now().In(loc)
var bot *linebot.Client


func tellTime(replyToken string, doTell bool){
	nowString := now.Format("2006-01-01 15:04:05")
	
	if doTell {
		log.Println("現在時間(台北): " + nowString)
		bot.ReplyMessage(replyToken, linebot.NewTextMessage("現在時間(台北): " + nowString)).Do()
	} else if silent != true {
		log.Println("自動報時(台北): " + nowString)
		bot.PushMessage(replyToken, linebot.NewTextMessage("自動報時(台北): " + nowString)).Do()
	} else {
		log.Println("tell time misfired")
	}
}

func routineDog(sourceId string) {
	for {
		time.Sleep(time.Duration(tellTimeInterval) * time.Minute)
		log.Println("time to tell time to : " + sourceId + ", " + now.Format("2006-01-02 15:04:05"))
		tellTime(sourceId, false)
	}
}

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			log.Println("keep alive at : " + now.Format("2006-01-02 15:04:05"))
			http.Get("https://line-talking-bot-go.herokuapp.com")
		}
	}()
}

func getSourceId(event *linebot.Event) string {
	var source = event.Source //EventSource
	var sourceId = source.UserID
	if sourceId != "" {
		log.Print("source UserID: " + sourceId)
		return sourceId
	}

	sourceId = source.GroupID
	if sourceId != "" {
		log.Print("source GroupID: " + sourceId)
		return sourceId
	}

	sourceId = source.RoomID
	if sourceId != "" {
		log.Print("source RoomID: " + sourceId)
		return sourceId
	}

	log.Print("Unknown source: " + sourceId)
	return sourceId
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		var replyToken = event.ReplyToken
		var sourceId = getSourceId(event)
		log.Print("callbackHandler to source id: " + sourceId)

		if sourceId != "" {
			if _, ok := echoMap[sourceId]; ok {
				//log.Print(sourceId + ": " + v)
			} else {
				log.Print("New routineDog added: " + sourceId)
				echoMap[sourceId] = true
				go routineDog(sourceId)
			}
		}

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Print("ReplyToken[" + replyToken + "] TextMessage: ID(" + message.ID + "), Text(" + message.Text  + "), current silent status=" + strconv.FormatBool(silent) )
				//if _, err = bot.ReplyMessage(replyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				//	log.Print(err)
				//}
				
				if strings.Contains(message.Text, "你閉嘴") {
					silent = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
				} else if "說吧" == message.Text {
					silent = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("麥克風測試，1、2、3... OK")).Do()
				} else if "time1" == message.Text {
					tellTimeInterval = 1					
				} else if "time15" == message.Text {
					tellTimeInterval = 15
				} else if strings.Contains(message.Text, "現在幾點")  {
					tellTime(replyToken, true)
				} else if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("嗯嗯，呵呵，我要去洗澡了")).Do()
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("> //////////////////// <")).Do()
				}
			case *linebot.ImageMessage :
				log.Print("ReplyToken[" + replyToken + "] ImageMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("傳這甚麼廢圖？你有認真在分享嗎？")).Do()
				}
			case *linebot.VideoMessage :
				log.Print("ReplyToken[" + replyToken + "] VideoMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("看甚麼影片，不知道流量快用光了嗎？")).Do()
				}
			case *linebot.AudioMessage :
				log.Print("ReplyToken[" + replyToken + "] AudioMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), Duration(" + strconv.Itoa(message.Duration) + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("說的比唱的好聽，唱得鬼哭神號，是要嚇唬誰？")).Do()
				}
			case *linebot.LocationMessage:
				log.Print("ReplyToken[" + replyToken + "] LocationMessage[" + message.ID + "] Title(" + message.Title  + "), Address(" + message.Address + "), Latitude(" + strconv.FormatFloat(message.Latitude, 'f', -1, 64) + "), Longitude(" + strconv.FormatFloat(message.Longitude, 'f', -1, 64) + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("這是哪裡啊？火星嗎？")).Do()
				}
			case *linebot.StickerMessage :
				log.Print("ReplyToken[" + replyToken + "] StickerMessage[" + message.ID + "] PackageID(" + message.PackageID + "), StickerID(" + message.StickerID + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("腳踏實地打字好嗎？傳這甚麼貼圖？")).Do()
				}
			}
		} else if event.Type == linebot.EventTypePostback {
		} else if event.Type == linebot.EventTypeBeacon {
		}
	}
}
