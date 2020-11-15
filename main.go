package main

import (
	"encoding/binary"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var text = [...]string{"A persistência é o caminho do êxito.", "No meio da dificuldade encontra-se a oportunidade.",
	"Eu faço da dificuldade a minha motivação. A volta por cima vem na continuação.",
	"Pedras no caminho? Eu guardo todas. Um dia vou construir um castelo.",
	"podem me atacar com paus e pedras, os paus eu chupo as pedra eu fumo."}

var (
	offset      uint32
	filePointer *os.File
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)

	//bot set up
	bot := MakeBot()
	//file opener
	go LoadFile(&wg)
	defer filePointer.Close()

	seed := time.Now()
	rand.Seed(seed.UnixNano())

	//handlers
	bot.Handle("/text", func(m *tb.Message) {
		_, err := bot.Send(m.Chat, getRandomText())
		if err != nil {
			log.Fatal(err)
		}
	})

	wg.Wait()

	bot.Handle("/photo", func(m *tb.Message) {
		p := &tb.Photo{File: tb.FromURL(getRandomLink())}
		_, err := bot.Send(m.Chat, p)

		if err != nil {
			log.Fatal(err)
		}
	})

	bot.Start()

}

func LoadFile(wg *sync.WaitGroup) {
	defer wg.Done()

	var err error

	filePointer, err = os.Open("links")
	if err != nil {
		log.Fatal(err)
	}

	err = binary.Read(filePointer, binary.LittleEndian, &offset)
	if err != nil {
		log.Fatal(err)
	}

}

func getRandomLink() string {

	buffer := make([]byte, 150)

	_, err := filePointer.Seek(int64(150*(rand.Int()%int(offset))), 0)
	if err != nil {
		log.Fatal(err)
	}

	err = binary.Read(filePointer, binary.LittleEndian, buffer)
	if err != nil {
		log.Fatal(err)
	}

	_, err = filePointer.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	var i uint16
	for i = 70; i < uint16(len(buffer)); i++ {
		if buffer[i] == '\n' {
			break
		}
	}

	return string(buffer[4:i])
}

func getRandomText() string {
	random := rand.Int() % 5

	return text[random]
}

func MakeBot() *tb.Bot {

	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: webhook,
	})

	if err != nil {
		log.Fatal(err)
	}

	return bot
}
