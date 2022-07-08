package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	//"strings"
)

//ИМЯ В РОДИТЕЛЬНОМ ПАДЕЖЕ
type rSuffix struct {
	Name string `json:"Р"`
}

//СТРУКТУРА ПАРСИНГА ИЗ ГУГЛ ТАБЛИЦ
type Employee struct {
	Name        string `json:"ФИО"`
	Date        string `json:"Дата рождения"`
	Title       string `json:"Должность"`
	Department  string `json:"Отдел"`
	PhoneNumber string `json:"Телефон"`
}

func main() {

	//botToken := getBotToken()
	bot, err := tgbotapi.NewBotAPI("5336741653:AAEFq8-y8i9O3f2mg0IqKXWWkQZ7i2Ivb64")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	for {
		//ПОЗДРАВЛЕНИЕ ТОЛЬКО В ПЕРИОД 10-12
		currentTime := time.Now()
		if currentTime.Hour() == 2 {

			birthdayToday := getBirthdayJson()
			birthdayToday[0].Name = getPrettySuffix(birthdayToday[0].Name)

			if len(birthdayToday) > 0 {
				for _, peoples := range birthdayToday {
					msg := fmt.Sprintf("Сегодня день рождения у %v", peoples.Name)
					bot.Send(tgbotapi.NewMessage(678187421, msg))
					time.Sleep(11 * time.Second)
				}
			}
			/*
				msg := fmt.Sprintf("Сегодня день рождения у %v", int(currentTime.Month()))
				bot.Send(tgbotapi.NewMessage(678187421, msg))
				time.Sleep(11 * time.Second)
				msg = fmt.Sprintf("Время2: %v", currentTime.Day())
				bot.Send(tgbotapi.NewMessage(678187421, msg))
				time.Sleep(10 * time.Second)
			*/

		} else {
			time.Sleep(1 * time.Hour)
		}
		time.Sleep(1 * time.Hour)
	}
}

//ПАРСИМ ЛЮДЕЙ У КОТОРЫХ СЕГОДНЯ ДЕНЬ РОЖДЕНИЯ
func getBirthdayJson() []Employee {
	resp, _ := http.Get(fmt.Sprintf("https://tools.aimylogic.com/api/googlesheet2json?sheet=Users&id=1mV5gdMfZ385RugZQAYLJQfljFFg17kWsb0DmZmG98dI"))
	defer resp.Body.Close()

	employes := []Employee{}

	body, err := ioutil.ReadAll(resp.Body) //ПОЛУЧИЛИ JSON
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if err := json.Unmarshal(body, &employes); err != nil {
		fmt.Println(err)
		return nil
	}

	employesBirthday := []Employee{} //СТРУКТУРА ЛЮДЕЙ С ДНЁМ РОЖДЕНИЯ
	currentTime := time.Now()

	var strMonth, strDay, strDate string

	//УЖАСНАЯ КОНВЕРТАЦИЯ
	if int(currentTime.Month()) < 10 {
		strMonth = fmt.Sprintf("0%v", int(currentTime.Month()))
	} else {
		strMonth = strconv.Itoa(int(currentTime.Month()))
	}

	//УЖАСНАЯ КОНВЕРТАЦИЯ
	if int(currentTime.Day()) < 10 {
		strDay = fmt.Sprintf("0%v", currentTime.Day())
	} else {
		strDay = strconv.Itoa(currentTime.Day())
	}

	strDate = strDay + "." + strMonth //ПРИВОДИМ ДАТУ К ВИДУ ГУГЛДОК

	//В ЦИКЛЕ ПО ВСЕМ ЛЮДЯМ ИЩЕМ ТЕХ У КОГО ДЕНЬ РОЖДЕНИЯ И ДОБАВЛЯЕМ ИХ В НОВУЮ СТРУКТУРУ
	for _, empl := range employes {
		if strings.HasPrefix(empl.Date, strDate) {
			shortName := strings.Split(empl.Name, " ")
			//ЕСЛИ ИМЯ БЕЗ ОСОБЕННОСТЕЙ - УБИРАЕМ ОТЧЕСТВО
			if len(shortName) == 3 {
				empl.Name = shortName[1] + " " + shortName[0]
			}
			employesBirthday = append(employesBirthday, empl)
		}
	}
	return employesBirthday
}

//ПОЛУЧИТЬ ИМЯ В РОДИТЕЛЬНОМ ПАДЕЖЕ
func getPrettySuffix(people string) string {

	people = strings.Replace(people, " ", "%20", -1)
	resp, err := http.Get(fmt.Sprint("http://ws3.morpher.ru/russian/declension?s=" + people + "&format=json"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	rodSuffix := rSuffix{}

	body, err := ioutil.ReadAll(resp.Body) //ПОЛУЧИЛИ JSON
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &rodSuffix); err != nil {
		fmt.Println(err)
	}

	people = rodSuffix.Name

	return people
}

/*

}

/*
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
*/
