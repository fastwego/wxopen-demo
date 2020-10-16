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
	"fmt"

	"github.com/fastwego/wxopen/apis/oauth"

	"github.com/gin-gonic/gin"
)

func OauthDemo(c *gin.Context) {

	// code=CODE&state=STATE&appid=APPID

	// 区分不同账号
	appid := c.Param("appid")

	code := c.Request.URL.Query().Get("code")
	state := c.Request.URL.Query().Get("state")

	if len(code) == 0 && len(state) == 0 {

		redirect_uri := "http://" + c.Request.Host + c.Request.RequestURI

		authorizeUrl, _ := oauth.GetAuthorizeUrl(myPlatform, appid, redirect_uri, "snsapi_userinfo", "STATE")

		c.Header("content-type", "text/html; charset=utf-8")
		c.Writer.WriteString("请点击授权：<a href='" + authorizeUrl + "'>" + authorizeUrl + "</a>")
		return
	}

	// 使用授权码获取授权信息
	accessToken, err := oauth.GetAccessToken(myPlatform, appid, code)
	if err != nil {
		return
	}

	// 获取用户信息
	userInfo, err := oauth.GetUserInfo(myPlatform, accessToken.AccessToken, accessToken.Openid)
	if err != nil {
		return
	}

	fmt.Println(userInfo, err)

	c.Writer.WriteString(fmt.Sprintf("%v %v", userInfo, err))
	return
}
