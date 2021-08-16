package nknovh_engine

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (o *NKNOVH) tgPoll() error {
	//Unsupported in v1.1.0
	return nil
	bot, err := tgbotapi.NewBotAPI(o.conf.Messengers.Telegram.Token)
	if err != nil {
		return err
	}
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//chat_id
		//link to hash

		// /reg hash_id
		// /help
		// /start instruct
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
	return nil
}