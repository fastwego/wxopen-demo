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
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/fastwego/offiaccount/util"

	"github.com/fastwego/offiaccount/apis/oauth"

	"github.com/gin-gonic/gin"
)

func Demo(c *gin.Context) {

	appid := c.Param("appid")

	// 优先从环缓存获取
	jsapi_ticket, err := myPlatform.Cache.Fetch("jsapi_ticket:" + appid)
	if len(jsapi_ticket) == 0 {

		// 创建 公众号
		offiAccount, err := myPlatform.NewOffiAccount(appid)
		if err != nil {
			return
		}

		var ttl int64
		jsapi_ticket, ttl, err = oauth.GetJSApiTicket(offiAccount) // 调用 公众号 jsapi_ticket 接口
		if err != nil {
			return
		}

		err = myPlatform.Cache.Save("jsapi_ticket:"+appid, jsapi_ticket, time.Duration(ttl)*time.Second)
		if err != nil {
			return
		}
	}

	// 生成签名
	nonceStr := util.GetRandString(6)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	pageUrl := "http://" + c.Request.Host + c.Request.RequestURI
	plain := "jsapi_ticket=" + jsapi_ticket + "&noncestr=" + nonceStr + "&timestamp=" + timestamp + "&url=" + pageUrl

	signature := fmt.Sprintf("%x", sha1.Sum([]byte(plain)))
	fmt.Println(plain, signature)

	configMap := map[string]string{
		"url":       pageUrl,
		"nonceStr":  nonceStr,
		"appid":     appid,
		"timestamp": timestamp,
		"signature": signature,
	}

	marshal, err := json.Marshal(configMap)
	if err != nil {
		return
	}

	config := template.JS(marshal)

	if err != nil {
		fmt.Println(err)
		return
	}

	t1, err := template.ParseFiles("jssdk/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	t1.Execute(c.Writer, config)
	return
}
