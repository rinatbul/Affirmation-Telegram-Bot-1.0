package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

// Структура для получения данных из API
type CatFact struct {
	Fact string `json:"fact"`
}

func getCatFact() (string, error) {
	resp, err := http.Get("https://catfact.ninja/fact")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var catFact CatFact
	if err := json.NewDecoder(resp.Body).Decode(&catFact); err != nil {
		return "", err
	}
	return catFact.Fact, nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI("---")
	if err != nil {
		log.Panic(err)
	}

	//Включаем режим отладки
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//Создаем новый апдейт конфиг
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//Настраиваем клавиатуру
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Get a Cat Fact"),
		),
	)

	// Обрабатываем обновления от бота
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		// Получаем имя пользователя
		userName := update.Message.From.FirstName
		if userName == "" {
			userName = update.Message.From.UserName
		}

		//Генерируем ответ на основе команды пользователя
		var responseText string
		if update.Message.Text == "Get a Cat Fact" {
			fact, err := getCatFact()
			if err != nil {
				responseText = "Sorry, I couldn't fetch a cat fact."
			} else {
				responseText = fact
			}
		} else {
			responseText = "Hello, " + userName + "!"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}
