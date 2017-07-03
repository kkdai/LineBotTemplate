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
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Constants
var timeFormat = "01/02 03:04:05PM"
var user_zchien = "U696bcb700dfc9254b27605374b86968b"
var user_yaoming = "U3aaab6c6248bb38f194134948c60f757"
var user_jackal = "U3effab06ddf5bcf0b46c1c60bcd39ef5"
var user_shane = "U2ade7ac4456cb3ca99ffdf9d7257329a"

// Global Settings
var tellTimeInterval int = 15
var answers_TextMessage = []string{
		"人被殺，就會死。",
		"凡是每天喝水的人，有高機率在100年內死去",
		"今年中秋節剛好是滿月、今年七夕恰逢鬼月、今年母親節正好是星期日",
		"只要每天省下買一杯奶茶的錢，十天後就能買十杯奶茶",
		"台灣人在睡覺時，大多數的美國人都在工作",
		"台灣競爭力低落，在美國就連小學生都會說流利的英語",
		"在非洲，每六十秒，就有一分鐘過去",
		"每呼吸60秒，就減少一分鐘的壽命",
		"身高170cm的女生看起來和身高1米7的女生一樣高",
		"英國研究證實，全世界的人口中，減去瘦子的人口數後，剩下來的都是胖子。",
		"張開你的眼睛！否則，你將什麼都看不見。",
		"嗯嗯，呵呵，我要去洗澡了",
		"當一個便當吃不飽時.你可以吃兩個",
		"當你吃下吃下廿碗白飯，換算竟相當於吃下了二十碗白飯的熱量",
		"當你的左臉被人打，那你的左臉就會痛",
		"當蝴蝶在南半球拍了兩下翅膀，牠就會稍微飛高一點點",
		"誰能想的到，這名16歲少女，在四年前，只是一名12歲少女",
		"據統計，未婚生子的人數中有高機率為女性",
		"在非洲，每一分鐘，就有六十秒過去。",
		"在你的面前閉氣的話，就會不能呼吸喔。",
		"廢話,日本,推特,氏くん,理所當然的詩",
		"跟你在一起時，回憶一天前的事，就像回想昨天的事情。",
		"你不在的這十二個月，對我來說如同一年般長。",
		"不知道為什麼，把眼睛矇上後什麼都看不到。",
		"出生時，大家都是裸體的喔。",
	}


var silentMap = make(map[string]bool) // [UserID/GroupID/RoomID]:bool

//var echoMap = make(map[string]bool)

var loc, _ = time.LoadLocation("Asia/Taipei")
var bot *linebot.Client


func tellTime(replyToken string, doTell bool){
var silent = false
	now := time.Now().In(loc)
	nowString := now.Format(timeFormat)
	
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
		now := time.Now().In(loc)
		log.Println("time to tell time to : " + sourceId + ", " + now.Format(timeFormat))
		tellTime(sourceId, false)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			now := time.Now().In(loc)
			log.Println("keep alive at : " + now.Format(timeFormat))
			http.Get("https://line-talking-bot-go.herokuapp.com")
			time.Sleep(time.Duration(rand.Int31n(59)) * time.Minute)
		}
	}()
	
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)

}

func getSourceId(event *linebot.Event) string {
	var source = event.Source //EventSource
	
	var sourceId = source.UserID
	if sourceId != "" {
		//log.Print("source UserID: " + sourceId)
		return sourceId
	}

	sourceId = source.GroupID
	if sourceId != "" {
		//log.Print("source GroupID: " + sourceId)
		return sourceId
	}

	sourceId = source.RoomID
	if sourceId != "" {
		//log.Print("source RoomID: " + sourceId)
		return sourceId
	}

	log.Print("Unknown source: " + sourceId)
	return sourceId
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	log.Print("URL:"  + r.URL.String())
	
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
		
		var userId = source.UserID
		var groupId = source.GroupID
		var roomId = source.RoomID
		log.Print("callbackHandler to source UserID/GroupID/RoomID: " + userId + "/" + groupId + "/" + roomId)
		
		var silent bool = false

		if event.Type == linebot.EventTypeMessage {
			silent, ok := silentMap[sourceId]
			
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
				} else if strings.Contains(message.Text, "現在幾點")  {
					tellTime(replyToken, true)
				} else if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_TextMessage[rand.Intn(len(answers_TextMessage))])).Do()
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
