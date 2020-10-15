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
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fastwego/wxopen/apis/auth"

	"github.com/gin-gonic/gin"
)

func AuthDemo(c *gin.Context) {

	auth_code := c.Request.URL.Query().Get("auth_code")
	if len(auth_code) == 0 {

		redirect_uri := "http://" + c.Request.Host + c.Request.RequestURI

		// 获取 预授权 码 pre_auth_code
		preAuthCodeParams := struct {
			ComponentAppid string `json:"component_appid"`
		}{
			ComponentAppid: myPlatform.Config.AppId,
		}
		payload, err := json.Marshal(preAuthCodeParams)
		if err != nil {
			fmt.Println(err)
			return
		}

		preauthCode, err := auth.CreatePreauthCode(myPlatform, payload)
		if err != nil {
			fmt.Println(err)
			return
		}

		preAuthCodeJson := struct {
			PreAuthCode string `json:"pre_auth_code"`
			ExpiresIn   int    `json:"expires_in"`
		}{}
		err = json.Unmarshal(preauthCode, &preAuthCodeJson)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(preAuthCodeJson)

		// 获取 跳转链接
		params := url.Values{}
		// component_appid=xxxx&pre_auth_code=xxxxx&redirect_uri=xxxx&auth_type=xxx
		params.Add("component_appid", myPlatform.Config.AppId)
		params.Add("pre_auth_code", preAuthCodeJson.PreAuthCode)
		params.Add("redirect_uri", redirect_uri)
		params.Add("auth_type", "3")
		uri := auth.GetAuthorizationRedirectUri(params)

		//action=bindcomponent&auth_type=3&no_scan=1&component_appid=xxxx&pre_auth_code=xxxxx&redirect_uri=xxxx&auth_type=xxx&biz_appid=xxxx#wechat_redirect
		//params.Add("action", "bindcomponent")
		//params.Add("auth_type", "3")
		//params.Add("no_scan", "1")
		//params.Add("component_appid", myPlatform.Config.AppId)
		//params.Add("pre_auth_code", preAuthCodeJson.PreAuthCode)
		//params.Add("redirect_uri", redirect_uri)
		//
		//uri := auth.GetAuthorizationRedirectUri2(params)

		c.Header("content-type", "text/html; charset=utf-8")
		c.Writer.WriteString("请点击授权：<a href='" + uri + "'>" + uri + "</a>")
		return
	}

	// 使用授权码获取授权信息
	apiQueryAuthParams := struct {
		ComponentAppid    string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}{
		ComponentAppid:    myPlatform.Config.AppId,
		AuthorizationCode: auth_code,
	}
	payload, err := json.Marshal(apiQueryAuthParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := auth.ApiQueryAuth(myPlatform, payload)

	apiQueryAuthResp := struct {
		AuthorizationInfo struct {
			AuthorizerAppid        string `json:"authorizer_appid"`
			AuthorizerAccessToken  string `json:"authorizer_access_token"`
			ExpiresIn              int    `json:"expires_in"`
			AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
			FuncInfo               []struct {
				FuncscopeCategory struct {
					ID int `json:"id"`
				} `json:"funcscope_category"`
			} `json:"func_info"`
		} `json:"authorization_info"`
	}{}

	err = json.Unmarshal(data, &apiQueryAuthResp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(apiQueryAuthResp)

	// demo 只将 authorizer_access_token 缓存到本地
	// 实际业务建议存放到 数据库
	err = myPlatform.Cache.Save("authorizer_access_token:"+apiQueryAuthResp.AuthorizationInfo.AuthorizerAppid, apiQueryAuthResp.AuthorizationInfo.AuthorizerAccessToken, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = myPlatform.Cache.Save("authorizer_refresh_token:"+apiQueryAuthResp.AuthorizationInfo.AuthorizerAppid, apiQueryAuthResp.AuthorizationInfo.AuthorizerRefreshToken, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Header("content-type", "text/html; charset=utf-8")
	c.Writer.WriteString("authorizer_access_token success, <a href='/api/wxopen/menu?action=/menu/get&appid=" + apiQueryAuthResp.AuthorizationInfo.AuthorizerAppid + "'> try it 公众号 </a>")

	c.Writer.WriteString("<hr />authorizer_access_token success, <a href='/api/wxopen/mini?appid=" + apiQueryAuthResp.AuthorizationInfo.AuthorizerAppid + "'> try it 小程序 </a>")

	// 获取授权账号信息
	params := struct {
		ComponentAppid  string `json:"component_appid"`
		AuthorizerAppid string `json:"authorizer_appid"`
	}{
		ComponentAppid:  myPlatform.Config.AppId,
		AuthorizerAppid: apiQueryAuthResp.AuthorizationInfo.AuthorizerAppid,
	}
	payload, err = json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err = auth.ApiGetAuthorizerInfo(myPlatform, payload)

	fmt.Println(string(data), err)
	c.Writer.WriteString(string(data))
	return
}
