package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"

	"samples/WebIM/controllers"
)

const (
	APP_VER = "0.1.1.0227"
)

func main() {
	beego.Info(beego.BConfig.AppName, APP_VER)

	beego.Router("/", &controllers.AppController{})
	// 自定义路由，如果是 post 方法走 Join 函数，发送 longpolling.html 页面
	beego.Router("/join", &controllers.AppController{}, "post:Join")

	// (this *AppController) Join()
	beego.Router("/lp", &controllers.LongPollingController{}, "get:Join")
	beego.Router("lp/post", &controllers.LongPollingController{})
	beego.Router("lp/fetch", &controllers.LongPollingController{}, "get:Fetch")

	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.Run()
}
