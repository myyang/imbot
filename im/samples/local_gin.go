package main

import (
	"github.com/gin-gonic/gin"
	"github.com/myyang/imbot/im/lark"
	"github.com/myyang/imbot/im/slack"
)

func main() {
	engine := gin.New()

	lark.RegisterGin(engine)
	slack.RegisterGin(engine)

	engine.Run(":8080")
}
