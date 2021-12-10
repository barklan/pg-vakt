package pkg

import (
	"log"
	"os"

	"github.com/caarlos0/env"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Data struct {
	B         *tb.Bot
	Chat      *tb.Chat
	Config    Config
	MediaPath string
}

type Config struct { // TODO
	TGBotToken      string `env:"PG_VAKT_TG_BOT_TOKEN"`
	MainChatId      int64  `env:"PG_VAKT_TG_CHAT_ID"`
	IntervalMinutes int    `env:"PG_VAKT_INTERVAL_MINUTES"`
	Continuous      bool   `env:"PG_VAKT_CONTINUOUS"`
	ContinuousPath  string `env:"PG_VAKT_CONTINUOUS_PATH"`
	SSHUser         string `env:"PG_VAKT_SSH_USER"`
	SSHHostname     string `env:"PG_VAKT_SSH_HOSTNAME"`
	SSHKeyFilename  string `env:"PG_VAKT_SSH_KEY_FILENAME"`
	ContainerName   string `env:"PG_VAKT_CONTAINER_NAME"`
	Database        string `env:"PG_VAKT_DATABASE"`
}

func InitData() *Data {
	data := Data{}

	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panic(err)
	}
	log.Println(cfg)
	data.Config = cfg

	data.MediaPath = "/app/media"

	// TODO
	data.Chat = GetMainChat(data.Config.MainChatId)

	// TODO
	data.B = Bot(data.Config.TGBotToken)

	return &data
}

// CreateMediaDirIfNotExists creates the directory in default media path.
// It can accept nested directory path, but all parent directories must
// exist. Returns full directory path.
func (d *Data) CreateMediaDirIfNotExists(dirname string) string {
	fullDirname := d.MediaPath + "/" + dirname
	_, err := os.Stat(fullDirname)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(fullDirname, 0755)
		if errDir != nil {
			log.Panic(err)
		}
	}

	return fullDirname
}

func (d *Data) baseSend(to tb.Recipient, msg interface{}, options ...interface{}) (*tb.Message, error) {
	m, err := d.B.Send(to, msg, options...)
	if err != nil {
		log.Printf("Failed to send tg message. %v", err)
		return nil, err
	}
	log.Printf("Send TG message %v\n", msg)
	return m, err
}

func (d *Data) JustSend(to tb.Recipient, msg interface{}, options ...interface{}) {
	go func() {
		d.baseSend(to, msg, options...)
	}()
}

func (d *Data) SendSync(msg interface{}, options ...interface{}) (*tb.Message, error) {
	return d.baseSend(d.Chat, msg, options...)
}

func (d *Data) Send(msg interface{}, options ...interface{}) {
	go func() {
		d.baseSend(d.Chat, msg, options...)
	}()
}
