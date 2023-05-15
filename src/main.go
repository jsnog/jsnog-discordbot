package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

var (
	token *string = flag.String("token", "", "Bot Token")
	guild *string = flag.String("guild", "", "Guild ID")
)

func main() {
	log.Print("Starting jsnog-discordbot")
	flag.Parse()
	if *token == "" || *guild == "" {
		log.Fatal("not specified nessesarry argument. please check command options(-token,-guild)")
	}

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
	go announceTopic(d, *guild)
	wg.Wait()

	defer d.Close()
}
