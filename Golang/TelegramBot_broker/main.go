package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type wallet map[string]float64

var db = map[int64]wallet{}

func main() {
	bot, err := tgbotapi.NewBotAPI("5405522760:AAFqA15HEI8tn--bRzzEd-TQiobMIv2AAEo")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // If we got a message
			continue
		}
		command := strings.Split(update.Message.Text, " ")

		switch command[0] {
		case "ADD":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			} else {
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				}
				if _, ok := db[update.Message.Chat.ID]; !ok {
					db[update.Message.Chat.ID] = wallet{}
				}
				db[update.Message.Chat.ID][command[1]] += amount
				balanceText := fmt.Sprintf("Баланс %v: %v", command[1], db[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
			}
		case "SUB":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			} else {
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				}
				if _, ok := db[update.Message.Chat.ID]; !ok {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка. Заданый тикер отсутствует."))
					//db[update.Message.Chat.ID] = wallet{}
				} else {
					db[update.Message.Chat.ID][command[1]] -= amount
					balanceText := fmt.Sprintf("Баланс %v: %v", command[1], db[update.Message.Chat.ID][command[1]])
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
				}
			}
		case "DEL":
			if len(command) != 2 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Акции добавлены"))
			}
		case "SHOW":
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не найдена"))

		}

	}
}
