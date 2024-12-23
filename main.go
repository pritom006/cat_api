package main

import (
	_ "catapigo/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    // Set static path for serving static files
    beego.SetStaticPath("/static", "./static")
}

func main() {
	beego.Run()
}

