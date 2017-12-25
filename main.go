package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var chatID = os.Getenv("CHAT_ID")
var botToken = os.Getenv("BOT_TOKEN")
var authToken = os.Getenv("AUTH_TOKEN")

var (
	bot *tgbotapi.BotAPI
	err error
)

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !intentAuth(r) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	log.Print(body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	chat, err := strconv.ParseInt(chatID, 10, 64)
	msg := tgbotapi.NewMessage(chat, string(body))
	bot.Send(msg)

	w.WriteHeader(http.StatusCreated)
}

func intentAuth(r *http.Request) bool {
	authHeaders := r.Header["Authorization"]
	if len(authHeaders) < 1 {
		return false
	}

	return authHeaders[0] == authToken
}

func main() {
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	http.HandleFunc("/events", eventHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
