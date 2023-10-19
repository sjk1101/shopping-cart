package controller

import (
	"log"

	"shopping-cart/service/thirdparty/bot"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

type BotControllerInterface interface {
	Repeat(ctx *gin.Context)
}

type botController struct {
	in ctrlIn
}

func newBotController(in ctrlIn) BotControllerInterface {
	return &botController{
		in: in,
	}
}

func (ctrl *botController) Repeat(ctx *gin.Context) {

	client := bot.GetLineBotClient()
	// 接收請求
	events, err := client.ParseRequest(ctx.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Println(err.Error())
			return
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// 回覆訊息
				if _, err = client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}
