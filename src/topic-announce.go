package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// channel
type channel struct {
	channel *discordgo.Channel
	// enableInThreadsがfalseの場合Threadsは定義しないこと
	threads         []*discordgo.Channel
	topic           string
	enableInThreads bool
	foreach         int
}

func SendTopic(d *discordgo.Session, guild string) error {
	// channels
	chs, err := d.GuildChannels(guild)
	if err != nil {
		return err
	}
	// ActiveThreads
	actiThs, err := d.GuildThreadsActive(guild)
	if err != nil {
		return err
	}
	sendchs, err := seekchannels(chs, actiThs)
	if err != nil {
		return err
	}

	// 送信関数(c:channel, t:topic, fe:foreach)
	send := func(c *discordgo.Channel, t string, fe int) {
		msgs, err := d.ChannelMessages(c.ID, fe, "", "", "")
		if err != nil {
			log.Printf("Error getting ChannelMessages: %s", err)
		}
		for i, m := range msgs {
			if m.Author.ID == d.State.User.ID {
				break
			}
			if i == len(msgs)-1 {
				_, err := d.ChannelMessageSend(c.ID, message+t)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
			}
		}
	}
	for _, c := range sendchs {
		// 通常チャンネルへの送信
		send(c.channel, c.topic, c.foreach)
		// チャンネル内スレッドへの送信
		if c.enableInThreads {
			for _, t := range c.threads {
				send(t, c.topic, c.foreach)
			}
		}
	}
	return nil
}

func seekchannels(dChs []*discordgo.Channel, dActiThs *discordgo.ThreadsList) ([]channel, error) {
	// channelにスレッドが紐づいていないため手動で紐づけ!
	threads := make(map[string][]*discordgo.Channel)
	for _, t := range dActiThs.Threads {
		threads[t.ParentID] = append(threads[t.ParentID], t)
	}
	// channels
	chs := make([]channel, 0)
dchsLoop:
	for _, c := range dChs {
		// トピックがないチャンネルは無視
		if c.Topic != "" {
			t := strings.Split(c.Topic, "\n")
			lastline := t[len(t)-1]
			ch := channel{}
			// デフォルト設定
			ch.channel = c
			ch.topic = c.Topic
			ch.enableInThreads = false
			ch.foreach = 100

			// prefixの有無
			if !strings.HasPrefix(lastline, prefix) {
				chs = append(chs, ch)
			} else {
				// Prefixがあった場合
				args := strings.Split(strings.TrimPrefix(lastline, prefix), " ")
				for _, arg := range args {
					opt := strings.Split(arg, "=")
					switch opt[0] {
					// トピックアナウンスの有効/無効
					case "enable":
						if opt[1] == "false" {
							continue dchsLoop
						}
					// チャンネル内スレッドのトピックアナウンスの有効/無効
					case "enableinthreads":
						if opt[1] == "true" {
							ch.enableInThreads = true
						} else {
							ch.enableInThreads = false
						}
					// どれだけの頻度でトピックをアナウンスするか
					case "foreach":
						foreach, err := strconv.Atoi(opt[1])
						if err != nil || foreach > 101 || foreach < 1 {
							ch.foreach = 100
						} else {
							ch.foreach = foreach
						}
					}
				}
				ch.topic = strings.Join(t[:len(t)-1], "\n")
				if ch.enableInThreads {
					// スレッドを紐づける
					ch.threads = threads[c.ID]
				}
				chs = append(chs, ch)
			}
		}
	}
	return chs, nil
}
