package bot

import (
	"shopping-cart/service/constant"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	err    error
	client *linebot.Client
)

func InitLineBot() error {
	client, err = linebot.New(constant.LineBotSecret, constant.LineBotToken)
	if err != nil {
		return err
	}

	return nil
}

func GetLineBotClient() *linebot.Client {
	return client
}
