package handler

import (
	"discord-bot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"strings"
)

type Message struct {
	DB *mongo.Client
}

var lastMessageID string

const (
	prefix         = "리챔아"
	usageCommand   = "$사용법"
	helloResponses = "안녕하세요! 저를 사용해보고 싶으시면 '$사용법'을 채팅창에 쳐주세요!"
)

func (d *Message) MessageInfoMsg(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println("MessageInfoMsg called with content:", m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ID == lastMessageID {
		return
	}
	lastMessageID = m.ID
	content := strings.ToLower(m.Content)

	if strings.HasPrefix(content, prefix) {
		switch content {
		case prefix + " 안녕", prefix + " 반가워", prefix + " 안녕!", prefix + " 반가워!":
			s.ChannelMessageSend(m.ChannelID, helloResponses)
		}
		if content == "리챔아" {
			s.ChannelMessageSend(m.ChannelID, "안녕하세연?")
		} else if strings.HasPrefix(content, prefix+" 배워") {
			data := strings.Split(m.Content, " ")
			answerIndex := 0
			var datas strings.Builder

			for i, word := range data {
				if word == prefix {
					continue
				}
				if word == "배워" {
					answerIndex = i + 1
					break
				}
			}
			for i := answerIndex + 1; i < len(data); i++ {
				datas.WriteString(data[i] + " ")
			}
			database.CreateMSG(s, m, d.DB, data[answerIndex], datas.String())
		} else if strings.HasPrefix(content, prefix) {
			word := strings.Split(m.Content, " ")

			results, haveResults := database.GetMSG(d.DB, word[1])
			log.Println(results)
			if haveResults {
				if len(results) == 0 {
					return
				} else {
					randNum := rand.Intn(len(results))
					s.ChannelMessageSend(m.ChannelID, results[randNum].MSG)
					s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced(fmt.Sprintf("쓴 새끼 %s", results[randNum].Author), fmt.Sprintf("만든 날짜%s", results[randNum].CreatedAt), 16705372))
					return
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "아직 없는 명령어이니 만들어 보아요!")
				return
			}
		}
	}
}

func (d *Message) SendingEmbed(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ID == lastMessageID {
		return
	}
	switch m.Content {
	case usageCommand:
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("명령어 목록", "듀\n듀댜?\n\n리챔아 배워 ooo\nooo을 하면 ooo을 함(?)", 16705372))
	case "듀":
		s.ChannelMessageSend(m.ChannelID, "듀다다듀")
	case "듀다?":
		s.ChannelMessageSend(m.ChannelID, "듀댜됴디")
	}
}
