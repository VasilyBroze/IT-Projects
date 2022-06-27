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

type bnResp struct { //BINANCE
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

type yfResp struct { //YAHOO FINANCE
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

		case "ADD": //ДОБАВИТЬ ТИКЕР
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			} else {
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверное количество"))
				}
				if _, ok := db[update.Message.Chat.ID]; !ok {
					db[update.Message.Chat.ID] = wallet{}
				}
				db[update.Message.Chat.ID][command[1]] += amount
				balanceText := fmt.Sprintf("Баланс %v: %v", command[1], db[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
			}

		case "SUB": //ОТНЯТЬ ТИКЕР
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			} else {
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверное количество"))
				}
				if _, ok := db[update.Message.Chat.ID]; !ok {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка. Заданый тикер отсутствует."))
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
					if price == 0 {
						price, _ = getPrice3(key)
					}
				}
				sum += value * price
				if price != 0 {
					msg += fmt.Sprintf("%s: %v [%.2f USD]\n", key, value, value*price)
				} else {
					msg += fmt.Sprintf("%s: %v [%.2f USD (Тикер не найден)]\n", key, value, value*price)
				}
			}
			msg += fmt.Sprintf("Общий балланс: %.2f USD\n", sum)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

		case "SHOWRUB":
			msg := ""
			var sum float64
			usd, _ := getPriceUSD()
			for key, value := range db[update.Message.Chat.ID] {
				price, _ := getPrice(key)
				if price == 0 {
					price, _ = getPrice2(key)
					if price == 0 {
						price, _ = getPrice3(key)
					}
				}
				sum += value * price * usd
				if price != 0 {
					msg += fmt.Sprintf("%s: %v [%.2f RUB]\n", key, value, value*price*usd)
				} else {
					msg += fmt.Sprintf("%s: %v [%.2f RUB (Тикер не найден)]\n", key, value, value*price)
				}
			}
			msg += fmt.Sprintf("Общий балланс: %.2f RUB\n", sum)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

		case "/description":
			msg := fmt.Sprintf("Описание комманд:\nADD (тикер) (количество) - добавить\nSUB (тикер) (количество) - отнять\nDEL (тикер) - удалить\nSHOW - баланс (USD)\nSHOWRUB - баланс (RUB)")
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

	if err := json.Unmarshal(body, &jsonResp); err != nil {
		panic(err)
	}

	body2, err := ioutil.ReadAll(resp2.Body)

	if err := json.Unmarshal(body2, &jsonRespUSD); err != nil {
		panic(err)
	}

	if jsonResp.QuoteSummary.Error != nil {
		return
	}

	price2 = (jsonResp.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw) / (jsonRespUSD.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw)

	return
}

func getPrice3(symbol string) (price3 float64, err error) { //АМЕРИКАНСКИЕ АКЦИИ
	resp, _ := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=price", symbol))

	if err != nil {
		return
	}

	defer resp.Body.Close()

	var jsonResp yfResp

	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &jsonResp); err != nil {
		panic(err)
	}

	if jsonResp.QuoteSummary.Error != nil {
		return
	}

	price3 = (jsonResp.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw)

	return
}

func getPriceUSD() (price4 float64, err error) {

	resp2, _ := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/RUB=X?modules=price"))

	defer resp2.Body.Close()

	var jsonRespUSD yfResp

	body2, err := ioutil.ReadAll(resp2.Body)

	if err := json.Unmarshal(body2, &jsonRespUSD); err != nil {
		panic(err)
	}

	if jsonRespUSD.QuoteSummary.Error != nil {
		return
	}

	price4 = (jsonRespUSD.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw)

	return
}
