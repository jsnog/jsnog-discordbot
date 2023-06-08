package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	prefix = "!jsnog-bot"
)

func announceTopic(discord *discordgo.Session, guild string) {
	// 1時間
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
			//Every 2 hour
			log.Print("announceTopic")
			channels, err := discord.GuildChannels(guild)
			if err != nil {
				log.Printf("Error getting GuildChannels: %s", err)
				break
			}
			// 全チャンネルのトピックを取得
			var topic string
			for _, c := range channels {
				// トピックがないチャンネルは無視
				if c.Topic == "" {
					continue
				}
				options := make(map[string]string)
				// Prefixがあった場合はそれに続くオプション通り、なかった場合はデフォルト動作する
				lastline := getLastLine(c.Topic)
				topic = c.Topic
				log.Printf("lastline: %v", strings.HasPrefix(prefix, lastline))
				if strings.HasPrefix(prefix, lastline) {
					// prefixがある場合最終行を除外
					topic = strings.Join(strings.SplitAfterN(c.Topic, "\n", -1), "\n")
					args := strings.Split(strings.TrimPrefix(prefix, lastline), ",")
					if len(args) < 4 {
						continue
					}
					for _, arg := range args {
						log.Printf("arg: %v", arg)
						options[strings.Split(arg, "=")[0]] = strings.Split(arg, "=")[1]
					}
				}
				if options["enable"] == "false" {
					log.Printf("enable=false")
					continue
				}
				foreach, err := strconv.Atoi(options["foreach"])
				if err != nil || foreach > 101 {
					foreach = 100
				}
				if options["enableinthreads"] == "false" {
					if c.Type == discordgo.ChannelTypeGuildPrivateThread || c.Type == discordgo.ChannelTypeGuildPublicThread {
						continue
					}
				}
				// メッセージ100件取得
				messages, err := discord.ChannelMessages(c.ID, foreach, "", "", "")
				if err != nil {
					log.Printf("Error getting ChannelMessages: %s", err)
					break
				}
				// メッセージ取得してBotのメッセージがなかったらトピックを送信(foreachメッセージごとに送信)
				for i, m := range messages {
					log.Printf("%v,m.Author.ID: %v == %v", i, m.Author.ID, discord.State.User.ID)
					// 非送信チャンネルの場合はその場で切り上げ
					if m.Author.ID == discord.State.User.ID {
						break
					}
					// 送信対象を追加
					if i == len(messages)-1 {
						discord.ChannelMessageSend(c.ID, message+topic)
						log.Printf("channel: %v send: %v", c.Name, message+topic)
					}
				}
			}
		}
	}
}

func getLastLine(str string) string {
	log.Printf("str: %v", str)
	log.Printf(strings.Split(str, "\n")[len(strings.Split(str, "\n"))-1])
	return strings.Split(str, "\n")[len(strings.Split(str, "\n"))-1]
}
