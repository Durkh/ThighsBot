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
	bot.Handle("/start", func(m *tb.Message) {
		_, err := bot.Send(m.Sender,
			"Use the unique command to send you a random thigh picture\n\nAs for now, there's only one command:"+
				"\n\t/thighs - Send you a thigh pic"+
				"\n\nhope you like the bot, we'll be trying to improve it as the time passes")

		if err != nil {
			log.Fatal(err)
		}
	})

	wg.Wait()

	bot.Handle("/thighs", func(m *tb.Message) {
		p := &tb.Photo{File: tb.FromURL(getRandomLink())}
		_, err := bot.Send(m.Sender, p)
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
