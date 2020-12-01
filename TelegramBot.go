package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
	"time"
)

type TelegramBot struct {
	trustedChatID int64
	botApi        *tgbotapi.BotAPI
	handlers      map[string]func(tgbotapi.Update)
	apiKey        string
}

func (bot *TelegramBot) Init() {
	bot.apiKey = os.Getenv("TELEGRAM_API_KEY")
	log.Printf("TELEGRAM_API_KEY = %s", bot.apiKey)
	bot.handlers = make(map[string]func(tgbotapi.Update))
	log.Println("Validating API Key")
	botApi, err := tgbotapi.NewBotAPI(bot.apiKey)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Valid API Key")
	}

	bot.botApi = botApi
}

func (bot *TelegramBot) Run() {
	log.Println("Running...")
	me, _ := bot.botApi.GetMe()
	log.Printf("ID=%d User=%s", me.ID, me.FirstName)

	messageLoop := func() {
		for {
			u := tgbotapi.NewUpdate(0)
			u.Timeout = 60

			updates, _ := bot.botApi.GetUpdatesChan(u)
			for update := range updates {
				if update.Message == nil { // ignore any non-Message Updates
					continue
				}

				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				bot.HandleMessage(update)
			}
			time.Sleep(1000)
		}
	}

	go messageLoop()
}

func (bot *TelegramBot) RegisterMessageHandler(handlerTrigger string, action func(tgbotapi.Update)) {
	bot.handlers[handlerTrigger] = action
}

func (bot *TelegramBot) HandleMessage(update tgbotapi.Update) {
	if update.Message.Text == bot.apiKey {
		bot.trustedChatID = update.Message.Chat.ID
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You may now communicate with the bot.")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.botApi.Send(msg)
		return
	}

	if bot.trustedChatID > 0 && update.Message.Chat.ID == bot.trustedChatID {
		findAny := false
		msgLower := strings.ToLower(update.Message.Text)
		for k := range bot.handlers {
			kLower := strings.ToLower(k)
			if strings.EqualFold(k,"*") || strings.Index(msgLower, kLower) == 0 {
				findAny = true
				go bot.handlers[k](update)
			}
		}

		if findAny == false {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The message was unhandled.")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.botApi.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send the bot's api token to prove your identity.")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.botApi.Send(msg)
	}
}
