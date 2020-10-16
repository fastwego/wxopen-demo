// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/fastwego/offiaccount/type/type_message"
	"github.com/gin-gonic/gin"
)

func MsgDemo(c *gin.Context) {

	// 区分不同账号
	appid := c.Param("appid")

	offiAccount, err := myPlatform.NewOffiAccount(appid)

	// 调用相应公众号服务
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))

	message, err := offiAccount.Server.ParseXML(body)
	if err != nil {
		log.Println(err)
	}

	var output interface{}
	switch message.(type) {
	case type_message.MessageText: // 文本 消息
		msg := message.(type_message.MessageText)

		// 回复文本消息
		output = type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeText,
			},
			Content: type_message.CDATA(msg.Content),
		}
	}

	myPlatform.Server.Response(c.Writer, c.Request, output)
}
