package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie2222对对对2@gmail.com"
	c.Data["Title"] = "feic点点滴滴hangh反反复复付a"
	c.TplName = "index.tpl"
}
