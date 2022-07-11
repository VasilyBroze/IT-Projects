package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//ИМЯ В НУЖНОМ ПАДЕЖЕ
type rSuffix struct {
	NameV string `json:"В"`
	NameD string `json:"Д"`
	NameR string `json:"Р"`
	Code  int    `json:"code"`
}

//СТРУКТУРА ПАРСИНГА ИЗ ПОЖЕЛАНИЙ
type TextFirstPart struct {
	Congratulation string `json:"Поздравляю"`
}

//СТРУКТУРА ПАРСИНГА ИЗ ПОЖЕЛАНИЙ
type TextSecondPart struct {
	WishYou string `json:"Желаю"`
}

//СТРУКТУРА ПАРСИНГА ИЗ ПОЖЕЛАНИЙ
type TextThirdPart struct {
	Sentiments string `json:"Пожелание"`
}

//СТРУКТУРА ПАРСИНГА ИЗ ГУГЛ ТАБЛИЦ
type Employee struct {
	Name string `json:"ФИО"`
	Date string `json:"Дата рождения"`
	//	Title       string `json:"Должность"`
	Department string `json:"Отдел"`
	//	PhoneNumber string `json:"Телефон"`
	Male string
}

func main() {

	//botToken := getBotToken()
	bot, err := tgbotapi.NewBotAPI("5336741653:AAEFq8-y8i9O3f2mg0IqKXWWkQZ7i2Ivb64") //БОТ ПОЗДРАВЛЯТОР ЛИИСОВИЧ
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	for {
		//ПОЗДРАВЛЕНИЕ ТОЛЬКО В ПЕРИОД 10-11
		currentTime := time.Now()
		if currentTime.Hour() == 14 {

			birthdayToday := getBirthdayJson()

			if len(birthdayToday) > 0 {
				for _, peoples := range birthdayToday {
					fmt.Println(peoples)
					msg := getBirthdayMsg(peoples)
					bot.Send(tgbotapi.NewMessage(-728590508, msg)) //ОТПРАВИТЬ В ТЕСТОВЫЙ ЧАТ
					//bot.Send(tgbotapi.NewMessage(678187421, msg)) //ОТПРАВИТЬ В ЛИЧНЫЙ ЧАТ
					time.Sleep(10 * time.Minute)
				}
			}
		} else {
			time.Sleep(1 * time.Hour)
		}
		time.Sleep(1 * time.Hour)
	}
}

//ПАРСИМ ЛЮДЕЙ У КОТОРЫХ СЕГОДНЯ ДЕНЬ РОЖДЕНИЯ, ОПРЕДЕЛЯЕМ ПОЛ
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

	//КОНВЕРТАЦИЯ МЕСЯЦА
	switch int(currentTime.Month()) {
	case 1:
		strMonth = "янв."
	case 2:
		strMonth = "февр"
	case 3:
		strMonth = "мар."
	case 4:
		strMonth = "апр."
	case 5:
		strMonth = "мая"
	case 6:
		strMonth = "июн."
	case 7:
		strMonth = "июл."
	case 8:
		strMonth = "авг."
	case 9:
		strMonth = "сент."
	case 10:
		strMonth = "окт."
	case 11:
		strMonth = "нояб."
	case 12:
		strMonth = "дек."

	}

	//КОНВЕРТАЦИЯ ДНЯ
	strDay = strconv.Itoa(currentTime.Day())

	strDate = strDay + " " + strMonth //ПРИВОДИМ ДАТУ К ВИДУ ГУГЛДОК

	//В ЦИКЛЕ ПО ВСЕМ ЛЮДЯМ ИЩЕМ ТЕХ У КОГО ДЕНЬ РОЖДЕНИЯ И ДОБАВЛЯЕМ ИХ В НОВУЮ СТРУКТУРУ
	for _, empl := range employes {
		if strings.HasPrefix(empl.Date, strDate) {
			shortName := strings.Split(empl.Name, " ")
			//ЕСЛИ ФИО ИЗ 3 СЛОВ - ОПРЕДЕЛЯЕМ ПОЛ ПО ОТЧЕСТВУ, УБИРАЕМ ОТЧЕСТВО
			if len(shortName) == 3 {
				switch {
				case
					strings.HasSuffix(shortName[2], "ч"):
					empl.Male = "М"
				case
					strings.HasSuffix(shortName[2], "а"):
					empl.Male = "Ж"
				default:
					empl.Male = "?"
				}
				empl.Name = shortName[1] + " " + shortName[0]
			} else {
				empl.Male = "?"
			}

			//ИЗМЕНЯЕМ НАЗВАНИЕ ОТДЕЛА НА БОЛЕЕ КОРОТКОЕ
			switch {
			case strings.Contains(empl.Department, "ПТО"):
				empl.Department = "Отдел ПТО"
			case strings.Contains(empl.Department, "(ПО)"):
				empl.Department = "Отдел IT"
			case strings.Contains(empl.Department, "ПНР"):
				empl.Department = "Отдел ПНР"
			case strings.Contains(empl.Department, "("): //СОКРАЩАЕМ НАЗВАНИЕ ОТДЕЛА ДО ПЕРВОЙ СКОБКИ
				dep := strings.Split(empl.Department, "(")
				if len(dep) > 0 {
					empl.Department = dep[0]
				}
			}

			employesBirthday = append(employesBirthday, empl)
		}
	}
	return employesBirthday
}

//ПОЛУЧАЕМ ИМЯ В НУЖНОМ ПАДЕЖЕ
func getPrettySuffix(people, padej string) string {
	name := people
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

	//ЕСЛИ НЕ ПОЛУЧИЛИ ИМЯ В НУЖНОМ ПАДЕЖЕ - ВОЗВРАЩАЕМ КАК ЕСТЬ
	if rodSuffix.Code != 0 {
		name = strings.Replace(people, "%20", " ", -1)
		fmt.Println("ОШИБКА СЕРВИСА ПАДЕЖЕЙ")
		return name
	}

	switch padej {
	case "V":
		name = rodSuffix.NameV
	case "D":
		name = rodSuffix.NameD
	case "R":
		name = rodSuffix.NameR
	}

	return name
}

//СЛУЧАЙНОЕ ЧИСЛО ДЛЯ ОПРЕДЕЛЕНИЯ ТЕКСТА СООБЩЕНИЯ
func random(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

//ПАРСИМ ТАБЛИЦУ С ТЕКСТОМ ПОЗДРАВЛЕНИЙ И РАСПРЕДЕЛЯЕМ ИХ ПО МАССИВАМ
func getCongratArrays() ([]TextFirstPart, []TextSecondPart, []TextThirdPart) {
	resp, _ := http.Get(fmt.Sprintf("https://tools.aimylogic.com/api/googlesheet2json?sheet=Text&id=1mV5gdMfZ385RugZQAYLJQfljFFg17kWsb0DmZmG98dI"))
	defer resp.Body.Close()

	//МАССИВЫ ДЛЯ ПАРСИНГА
	fTP := []TextFirstPart{}
	sTP := []TextSecondPart{}
	tTP := []TextThirdPart{}

	fTPraw := []TextFirstPart{}
	sTPraw := []TextSecondPart{}
	tTPraw := []TextThirdPart{}

	body, err := ioutil.ReadAll(resp.Body) //ПОЛУЧИЛИ JSON
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &fTP); err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := json.Unmarshal(body, &sTP); err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(body, &tTP); err != nil {
		fmt.Println(err)
	}

	//ФИЛЬТРУЕМ ПУСТЫЕ СТРОКИ
	for _, first := range fTP {
		if first.Congratulation != "" {
			fTPraw = append(fTPraw, first)
		}
	}

	for _, second := range sTP {
		if second.WishYou != "" {
			sTPraw = append(sTPraw, second)
		}
	}

	for _, third := range tTP {
		if third.Sentiments != "" {
			tTPraw = append(tTPraw, third)
		}
	}

	return fTPraw, sTPraw, tTPraw
}

//ГЕНЕРИРУЕМ СООБЩЕНИЕ ПО ГУГЛ ТАБЛИЦЕ С ЗАГОТОВКАМИ
func getBirthdayMsg(peoples Employee) string {
	//МАССИВЫ СТРУКТУР ЧАСТЕЙ ПОЗДРАВЛЕНИЯ
	fTP, sTP, tTP := getCongratArrays()

	var text1, text2, text3, text4, text5 string

	//ГЕНЕРИРУЕМ СЛУЧАЙНОЕ ЧИСЛО, И ПО НЕМУ ПОДСТАВЛЯЕМ ЧАСТЬ ТЕКСТА
	text1 = fTP[random(len(fTP))].Congratulation

	//ЕСЛИ В ПЕРВОЙ ЧАСТИ УКАЗАН ПОЛ, ПРОВЕРЯЕМ ПОЛ СОТРУДНИКА
	for strings.HasSuffix(text1, " *Ж") && peoples.Male == "М" {
		text1 = fTP[random(len(fTP))].Congratulation
		time.Sleep(77 * time.Microsecond)
	}
	for strings.HasSuffix(text1, " *М") && peoples.Male == "Ж" {
		text1 = fTP[random(len(fTP))].Congratulation
		time.Sleep(77 * time.Microsecond)
	}
	//ЕСЛИ ПОЛ ОПРЕДЕЛИТЬ НЕ УДАЛОСЬ - НЕ ИСПОЛЬЗУЕМ НАЧАЛЬНЫЕ ФРАЗЫ В КОТОРЫХ ОН УКАЗАН
	for peoples.Male == "?" && (strings.HasSuffix(text1, " *М") || strings.HasSuffix(text1, " *Ж")) {
		text1 = fTP[random(len(fTP))].Congratulation
		time.Sleep(77 * time.Microsecond)
	}

	//УДАЛЯЕМ УКАЗАТЕЛИ ПОЛА В НАЧАЛЬНОЙ ФРАЗЕ
	if strings.HasSuffix(text1, " *Ж") {
		text1 = strings.Replace(text1, " *Ж", "", 1)
	}
	if strings.HasSuffix(text1, " *М") {
		text1 = strings.Replace(text1, " *М", "", 1)
	}
	//ПОЛУЧАЕМ ИМЯ В НУЖНОМ ПАДЕЖЕ В ЗАВИСИМОСТИ ОТ УКАЗАТЕЛЯ ПАДЕЖА В ГУГЛ ТАБЛИЦЕ
	if strings.HasSuffix(text1, " *В") {
		text1 = strings.Replace(text1, " *В", "", 1)
		peoples.Name = getPrettySuffix(peoples.Name, "V")
	}
	if strings.HasSuffix(text1, " *Д") {
		text1 = strings.Replace(text1, " *Д", "", 1)
		peoples.Name = getPrettySuffix(peoples.Name, "D")
	}
	if strings.HasSuffix(text1, " *Р") {
		text1 = strings.Replace(text1, " *Р", "", 1)
		peoples.Name = getPrettySuffix(peoples.Name, "R")
	}

	//ПОЛУЧАЕМ ОТДЕЛ В НУЖНОМ ПАДЕЖЕ
	peoples.Department = getPrettySuffix(peoples.Department, "R")

	text2 = sTP[random(len(sTP))].WishYou
	//ПОЖЕЛАНИЯ ГЕНЕРИРУЕМ ТАК, ЧТОБЫ ОНИ НЕ ПОВТОРЯЛИСЬ
	text3 = tTP[random(len(tTP))].Sentiments
	for text4 == "" || text4 == text3 {
		text4 = tTP[random(len(tTP))].Sentiments
	}
	for text5 == "" || text5 == text4 || text5 == text3 {
		text5 = tTP[random(len(tTP))].Sentiments
	}
	msg := fmt.Sprintf("%v %v из %v. %v %v, %v и %v!", text1, peoples.Name, peoples.Department, text2, text3, text4, text5)
	return msg
}
