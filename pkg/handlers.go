package pkg

import (
	"fmt"
	"math/rand"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var randSrc = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func registerTemporaryNginxHandler(b *tb.Bot, data *Data) {
	b.Handle("/nginx", func(m *tb.Message) {
		ok := VerifySender(data, m)
		if !ok {
			return
		}

		minutes := 5
		basicAuthUsername := "nginx"
		basicAuthPassword := RandString(24)

		addressChan := make(chan string)
		go TemporaryNginx(
			data,
			addressChan,
			minutes,
			basicAuthUsername,
			basicAuthPassword,
		)

		address := <-addressChan

		wgetCmd := fmt.Sprintf(
			"`wget --no-verbose --no-parent --recursive --level=1 --no-directories --http-user=%s --http-password=%s %s`",
			basicAuthUsername,
			basicAuthPassword,
			address,
		)

		data.Send(
			fmt.Sprintf(
				`Temporary nginx started at %s for %d minutes.
To download all files use something like this:
%s`,
				address,
				minutes,
				wgetCmd,
			),
			tb.ModeMarkdown,
		)
	})
}

func RegisterHandlers(b *tb.Bot, data *Data) {
	registerTemporaryNginxHandler(b, data)
}

func VerifySender(data *Data, m *tb.Message) bool {
	if m.Chat.ID == data.Chat.ID {
		return true
	}
	data.JustSend(m.Chat, "I am not registered for this chat.")
	return false
}
