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
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"net/url"
  	"io/ioutil"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	// fixie
	fixieUrl, err := url.Parse(os.Getenv("FIXIE_URL"))
  	customClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(fixieUrl)}}
  	resp, err := customClient.Get("http://welcome.usefixie.com")
  	if (err != nil) {
    		println(err.Error())
    		return
  	}
  	defer resp.Body.Close()
  	body, err := ioutil.ReadAll(resp.Body)
  	println(string(body))
	
	// end fixie
	
	// line bot
	strID := os.Getenv("ChannelID")
	numID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		log.Fatal("Wrong environment setting about ChannelID")
	}

	bot, err = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	received, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, result := range received.Results {
		content := result.Content()
		log.Println("-->", content)
		// add eggyo geo test
		resp, err := http.Get("http://eggyo-geo-node.herokuapp.com/geo/" + content.Text)
		if (err != nil) {
    			println(err.Error())
    			return
  		}
  		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		println(string(body))

		//Log detail receive content
		if content != nil {
			log.Println("RECEIVE Msg:", content.IsMessage, " OP:", content.IsOperation, " type:", content.ContentType, " from:", content.From, "to:", content.To, " ID:", content.ID)
		}
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
			text, err := content.TextContent()
			_, err = bot.SendText([]string{content.From}, "OK "+text.Text)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
