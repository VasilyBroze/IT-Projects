package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bnResp struct {
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

type yfResp struct {
	QuoteSummary struct {
		Result []struct {
			Price struct {
				RegularMarketPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketPrice"`
			} `json:"price"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteSummary"`
}

/*type yfResp struct {
	//Price  float64 `json:"quoteSummary.result[0].price.regularMarketPrice.raw"
	Price float64 `json:"quoteSummary.result[0].price.regularMarketPrice.raw,string"`
	//Result string `json:"quoteSummary.result"`
}*/

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

				delete(db[update.Message.Chat.ID], command[1])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Тикер удалён"))
			}
		case "SHOW":
			msg := ""
			var sum float64
			for key, value := range db[update.Message.Chat.ID] {
				price, _ := getPrice(key)
				if price == 0 {
					price, _ = getPrice2(key)
				}
				sum += value * price
				msg += fmt.Sprintf("%s: %v [%.2f USD]\n", key, value, value*price)
			}
			msg += fmt.Sprintf("Общий балланс: %.2f USD\n", sum)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

		case "SHOW2":
			msg := ""
			var sum float64
			for key, value := range db[update.Message.Chat.ID] {
				price, _ := getPrice2(key)
				sum += value * price
				msg += fmt.Sprintf("%s: %v [%.2f USD]\n", key, value, value*price)
			}
			msg += fmt.Sprintf("Общий балланс: %.2f USD\n", sum)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

		case "/description":
			msg := fmt.Sprintf("Описание комманд:\nADD (тикер) (количество) - добавить\nSUB (тикер) (количество) - отнять\nDEL (тикер) - удалить\nSHOW - баланс")
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не найдена"))
		}
	}
}

func getPrice(symbol string) (price float64, err error) {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var jsonResp bnResp

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return
	}
	if jsonResp.Code != 0 {
		err = errors.New("Неверный символ")
	}

	price = jsonResp.Price
	return
}

func getPrice2(symbol string) (price2 float64, err error) { //РУБЛЁВЫЕ АКЦИИ
	resp, _ := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s.ME?modules=price", symbol))

	if err != nil {
		return
	}

	resp2, _ := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/RUB=X?modules=price"))

	defer resp.Body.Close()
	defer resp2.Body.Close()

	var jsonResp yfResp
	var jsonRespUSD yfResp

	body, err := ioutil.ReadAll(resp.Body)
	/*generatedjson := (string(body))
	fmt.Println(generatedjson)*/
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		panic(err)
	}

	body2, err := ioutil.ReadAll(resp2.Body)
	/*generatedjson := (string(body2))
	fmt.Println("НАЧАЛО ДЖЕСОНА")
	fmt.Println(generatedjson)
	fmt.Println("КОНЕЦ ДЖЕСОНА")*/
	if err := json.Unmarshal(body2, &jsonRespUSD); err != nil {
		panic(err)
	}

	price2 = (jsonResp.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw) / (jsonRespUSD.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw)
	if jsonResp.QuoteSummary.Error != "null" {
		err = errors.New("Неверный тикер")
	}

	return
}
