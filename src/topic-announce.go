package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func announceTopic(discord *discordgo.Session, guild string) {
	// 1時間
	t := time.NewTicker(2 * time.Hour)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//Every 2 hour
			log.Print("announceTopic")
			channels, err := discord.GuildChannels(guild)
			if err != nil {
				log.Printf("Error getting GuildChannels: %s", err)
				break
			}
			sendChannels := make([]*discordgo.Channel, 0)
			// 全チャンネルのトピックを取得
			for _, c := range channels {
				// テキストチャンネル以外は無視
				if c.Type != discordgo.ChannelTypeGuildText {
					continue
				}
				// トピックがないチャンネルは無視
				if c.Topic == "" {
					continue
				}
				// メッセージ100件取得
				messages, err := discord.ChannelMessages(c.ID, 100, "", "", "")
				if err != nil {
					log.Printf("Error getting ChannelMessages: %s", err)
					break
				}
				// 100メッセージ取得してBotのメッセージがなかったらトピックを送信(100メッセージごとに送信)
				for i, m := range messages {
					log.Printf("%v,m.Author.ID: %v == %v", i, m.Author.ID, discord.State.User.ID)
					if m.Author.ID == discord.State.User.ID {
						log.Printf("true")
						break
					}
					if i == len(messages)-1 {
						log.Printf("false")
						sendChannels = append(sendChannels, c)
					}
				}
			}
			//メッセージ送出
			log.Printf("sendChannels: %v", sendChannels)
			for _, c := range sendChannels {
				discord.ChannelMessageSend(c.ID, "【定期広報】チャンネル内ルール :\n"+c.Topic)
				// 通知爆弾にならないよう10分間隔で送信
				time.Sleep(10 * time.Minute)
			}
		}
	}
}
