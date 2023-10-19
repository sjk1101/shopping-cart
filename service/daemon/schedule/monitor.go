package schedule

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"shopping-cart/service/thirdparty/bot"
	"shopping-cart/service/thirdparty/database"

	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

func newMonitorCronJob(in scheduleIn) JobInterface {
	return &monitorCronJob{
		in: in,
	}
}

type monitorCronJob struct {
	in scheduleIn
}

func (j *monitorCronJob) GetExpired() string {
	return "1 * * * * *"
}

func (j *monitorCronJob) Exec(ctx context.Context) {
	db := database.Session()
	client := bot.GetLineBotClient()
	products, err := j.in.ProductRepo.Find(ctx, db, func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("quantity < ?", 5)
		return tx
	})
	if err != nil {
		log.Println(err.Error())
	}

	segmentCount := int(math.Ceil(float64(len(products)) / float64(10)))

	for i := 0; i < segmentCount; i++ {
		start := i * 10
		end := (i + 1) * 10
		if end > len(products) {
			end = len(products)
		}
		messages := []string{}
		for _, v := range products[start:end] {
			messages = append(messages, fmt.Sprintf("the inventory of %s is: %d", v.Name, v.Quantity))
		}
		msg := linebot.NewTextMessage(strings.Join(messages, "\n"))
		if _, err := client.BroadcastMessage(msg).Do(); err != nil {
			log.Fatal(err)
		}
	}
}
