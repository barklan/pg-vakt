package pkg

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func Bot(botToken string) *tb.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panic(err)
		return nil
	}
	return b
}

func GetMainChat(chatID int64) *tb.Chat {
	return &tb.Chat{ID: chatID}
}
