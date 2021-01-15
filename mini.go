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
	"log"

	"github.com/fastwego/miniprogram/apis/operation"

	"github.com/gin-gonic/gin"
)

func MiniDemo(c *gin.Context) {

	// 已授权 小程序 appid
	appid := c.Request.URL.Query().Get("appid")

	if len(appid) == 0 {
		log.Println("appid not found")
		return
	}

	mini, err := myPlatform.NewMiniprogram(appid)
	if err != nil {
		return
	}

	feedback, err := operation.GetFeedback(mini)
	log.Println(string(feedback), err)

	c.Writer.Write(feedback)
}
