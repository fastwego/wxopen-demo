## wxopen-demo 

A demo project for fastwego/wxopen

[![powered-by-fastwego](https://img.shields.io/badge/Powered%20By-fastwego-brightgreen)](https://github.com/fastwego)

### Install
- checkout project `git clone https://github.com/fastwego/wxopen-demo.git`
- install fastwego/wxopen `go get -u github.com/fastwego/wxopen`
- build `go build`
- edit config in `.env.dist` file and rename to `.env`
- run `wxopen-demo` & view `http://localhost`
- that's all & good luck ;)

### use case demo

- [请求授权](auth.go)
- [代 公众号 调用接口](menu.go)
- [处理公众号 消息/通知](msg.go)
- [代 公众号 发起网页授权](oauth.go)
- [代代 公众号 使用 js-sdk](jssdk.go)
- [代 小程序 调用接口](mini.go)

### tips

使用前请阅读 官方文档 关于 验证票据（component_verify_ticket）：

https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/component_verify_ticket.html