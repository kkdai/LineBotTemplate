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
// limitations under the License.11

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
	

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

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
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				output := sqlConnect(message.Text)
				// fmt.printf("%q", output)
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(output)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func sqlConnect(currency string)(output string){
	// var output string
	var (
		cashbuy float64
		cashsell float64
		ratebuy float64
		ratesell float64
		datetime string
	)
	
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	checkErr(err)
	
	// rows, err := db.Query("SELECT * FROM $1 ORDER BY id DESC LIMIT 1;", currency)
	rows, err := db.Query("SELECT cashbuy, cashsell, ratebuy, ratesell, datetime FROM bot_"+currency+" ORDER BY id DESC LIMIT 1;")
	checkErr(err)
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&cashbuy, &cashsell, &ratebuy, &ratesell, &datetime)
		checkErr(err)
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, datetime)

		output = "台銀"+currency+"即時匯率:"+
					"\n 現金買入:"+strconv.FormatFloat(cashbuy, 'f', 4, 64)+
					"\n 現金賣出:"+strconv.FormatFloat(cashsell, 'f', 4, 64)+
					"\n 即期買入:"+strconv.FormatFloat(ratebuy, 'f', 4, 64)+
					"\n 即期賣出:"+strconv.FormatFloat(ratesell, 'f', 4, 64)+
					"\n 更新時間("+t.Format("2006/01/02-15:04:05")+")"
	}
	return
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
