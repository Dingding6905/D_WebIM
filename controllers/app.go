package controllers

import (
	"strings"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

var langTypes []string // Languages that are supported.

func init() {
	fmt.Println("=====app init()")
	// Initialize language type list.
	langTypes = strings.Split(beego.AppConfig.String("lang_types"), "|")

	// Load locale files according to language types.
	for _, lang := range langTypes {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file:", err)
			return
		}
	}
}

type baseController struct {
	beego.Controller
	i18n.Locale
}

func (this *baseController) Prepare() {
	fmt.Println("=====app (this *baseController) Prepare()")
	this.Lang = ""

	al := this.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5]
		if i18n.IsExist(al) {
			this.Lang = al
		}
	}

	if len(this.Lang) == 0 {
		this.Lang = "en-US"
	}

	this.Data["Lang"] = this.Lang
}

type AppController struct {
	baseController
}

func (this *AppController) Get() {
	fmt.Println("=====app (this *AppController) Get()")
	this.TplName = "welcome.html"
}

func (this *AppController) Join() {
	fmt.Println("=====app (this *AppController) Join()")
	uname := this.GetString("uname")
	tech := this.GetString("tech")

	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	switch tech {
	case "longpolling":
		this.Redirect("/lp?uname="+uname, 302)
	case "websocket":
		this.Redirect("/ws?uname="+uname, 302)
	default:
		this.Redirect("/", 302)
	}
	return
}
