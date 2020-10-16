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
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastwego/wxopen/type/type_platform"

	"github.com/fastwego/wxopen"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

var myPlatform *wxopen.Platform

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	myPlatform = wxopen.NewPlatform(wxopen.PlatformConfig{
		AppId:     viper.GetString("APPID"),
		AppSecret: viper.GetString("APPSECRET"),
		Token:     viper.GetString("TOKEN"),
		AesKey:    viper.GetString("AESKEY"),
	})
}

func HandleEvent(c *gin.Context) {

	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))

	message, err := myPlatform.Server.ParseXML(body)
	if err != nil {
		log.Println(err)
	}

	var output interface{}
	switch message.(type) {
	case type_platform.EventComponentVerifyTicket:
		msg := message.(type_platform.EventComponentVerifyTicket) // 存储 ComponentVerifyTicket

		err := myPlatform.ReceiveComponentVerifyTicketHandler(myPlatform, msg.ComponentVerifyTicket)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = myPlatform.Server.Response(c.Writer, c.Request, output) // 响应 success
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.POST("/wechat/Auth/index", HandleEvent)

	// 请求授权
	router.GET("/api/wxopen/auth", AuthDemo)

	// 代 公众号 调用接口
	router.GET("/api/wxopen/menu", MenuDemo)

	// 处理公众号 消息/通知
	router.POST("/wechat/Message/index/:appid", MsgDemo)

	// 代 公众号 发起网页授权
	router.GET("/api/wxopen/oauth/:appid", OauthDemo)

	// 代 公众号 使用 js-sdk
	router.GET("/api/wxopen/jssdk/:appid", Demo)

	// 代 小程序 调用接口
	router.GET("/api/wxopen/mini", MiniDemo)

	svr := &http.Server{
		Addr:    viper.GetString("LISTEN"),
		Handler: router,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
