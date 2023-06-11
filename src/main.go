package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	token *string = flag.String("token", "", "Bot Token")
	guild *string = flag.String("guild", "", "Guild ID")
	debug *bool   = flag.Bool("debug", false, "debug mode")
)

const (
	prefix  = "!jsnog-bot"
	message = "【定期広報】チャンネル内ルール :\n"
)

func main() {
	log.Print("Starting jsnog-discordbot")
	flag.Parse()
	if *token == "" || *guild == "" {
		log.Fatal("not specified nessesarry argument. please check command options(-token,-guild)")
	}

	// スレッド関連があるのでAPIVersionを10に固定
	discordgo.APIVersion = "10"

	d, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal("Error creating Discord Session: ", err)
	}
	err = d.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}
	defer d.Close()
	// 無限ループなので並列化
	var wg sync.WaitGroup
	wg.Add(1)
	// トピック送信関数
	go func() {
		var t time.Duration
		if *debug == true {
			t = 10 * time.Second
		} else {
			t = 1 * time.Hour
		}
		timer := time.NewTicker(t)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				//Every 1hour
				err := SendTopic(d, *guild)
				log.Printf("SendTopic")
				if err != nil {
					log.Printf("Error sending topic: %s", err)
				}
			}
		}
	}()
	wg.Wait()
}
