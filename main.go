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
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var silent bool;
var echoMap = make(map[string]bool);

var bot *linebot.Client


func tellTime(event *linebot.Event, timeString string){
	if timeString == "" {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("現在時間是: " + time.Now().Format("2006-01-02 15:04:05"))).Do();
	} else if silent != true {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("自動報時: " + timeString)).Do();
	}				
}

func routineDog(event *linebot.Event){
	for {
		time.Sleep(15 * 60 * 1000 * time.Millisecond) //time.Sleep(100 * time.Millisecond)
		tellTime(event, time.Now().Format("2006-01-02 15:04:05"));
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
	
		var eventSource = event.Source; //EventSource
		var sourceId String;
		sourceId := eventSource.GroupID;
		if sourceId== nil {
			sourceId := eventSource.RoomID
		}
		if sourceId != nil {
			if v, ok := echoMap[sourceId]; ok {
			} else {
				log.Print("New routineDog added: " + sourceId)
				echoMap[sourceId] = true
				go routineDog(event)
			}
		}

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Print("TextMessage: ID(" + message.ID + "), Text(" + message.Text  + "), current silent status=" + strconv.FormatBool(silent) )
				//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				//	log.Print(err)
				//}
				
				//log.Print("現在幾點 == " + message.Text + " is " +strconv.FormatBool("現在幾點" == message.Text)) 
				
				if "你閉嘴" == message.Text {
					silent = true;
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("(X!)")).Do();
				} else if "說吧" == message.Text {
					silent = false;
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("麥克風測試，1、2、3... OK")).Do();
				} else if "現在幾點" == message.Text {
					tellTime(event, "");
				} else if silent != true {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("嗯嗯，呵呵，我要去洗澡了")).Do();
				}
			case *linebot.ImageMessage :
				log.Print("ImageMessage: ID(" + message.ID + "), OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("傳這甚麼廢圖？你有認真在分享嗎？")).Do();
				}
			case *linebot.VideoMessage :
				log.Print("VideoMessage: ID(" + message.ID + "), OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("看甚麼影片，不知道流量快用光了嗎？")).Do();
				}
			case *linebot.AudioMessage :
				log.Print("AudioMessage: ID(" + message.ID + "), OriginalContentURL(" + message.OriginalContentURL + "), Duration(" + strconv.Itoa(message.Duration) + ")" )
				if silent != true {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("說的比唱的好聽，唱得鬼哭神號，是要嚇唬誰？")).Do();
				}
			case *linebot.LocationMessage:
				log.Print("LocationMessage: ID(" + message.ID + "), Title(" + message.Title  + "), Address(" + message.Address + "), Latitude(" + strconv.FormatFloat(message.Latitude, 'f', -1, 64) + "), Longitude(" + strconv.FormatFloat(message.Longitude, 'f', -1, 64) + ")" )
				if silent != true {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("這是哪裡啊？火星嗎？")).Do();
				}
			case *linebot.StickerMessage :
				log.Print("StickerMessage: ID(" + message.ID + "), PackageID(" + message.PackageID + "), StickerID(" + message.StickerID + ")" )
				if silent != true {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("腳踏實地打字好嗎？傳這甚麼貼圖？")).Do();
				}
			}
		} else if event.Type == linebot.EventTypePostback {
		} else if event.Type == linebot.EventTypeBeacon {
		}
	}
}
