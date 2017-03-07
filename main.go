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
		id int
		cashbuy float64
	)
	
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	checkErr(err)
	
	// rows, err := db.Query("SELECT * FROM $1 ORDER BY id DESC LIMIT 1;", currency)
	rows, err := db.Query("SELECT * FROM "+currency+" ORDER BY id DESC LIMIT 1;")
	checkErr(err)
	defer rows.Close()
	for rows.Next(){
		// var id int
		// var cashbuy float32
		// var cashsell float32
		// var ratebuy float32
		// var ratesell float32
		// var datetime string
		// err = rows.Scan(&id, &cashbuy, &cashsell, &ratebuy, &ratesell, &datetime)
		err := rows.Scan(&id, &cashbuy)
		checkErr(err)
		output = strconv.FormatFloat(cashbuy, 'f', 4, 64)
	}
	return
	// for rows.Next(){
		// output = "日幣現金賣出:"+rows.+""
	// }
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
