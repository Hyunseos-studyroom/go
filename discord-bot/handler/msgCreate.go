package handler

import (
	"discord-bot/database"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type Message struct {
	DB *mongo.Client
}

func (d *Message) MessageInfoMsg(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.ToLower(m.Content) // 메시지 내용을 소문자로 변환

	if strings.HasPrefix(content, "리챔아") {
		switch content {
		case "리챔아 안녕", "리챔아 반가워", "리챔아 안녕!", "리챔아 반가워!":
			s.ChannelMessageSend(m.ChannelID, "안녕하세요! 저를 사용해보고 싶으시면 '$사용법'을 채팅창에 쳐주세요!")
		}

		if strings.HasPrefix(m.Content, "리챔아 배워") {
			data := strings.Split(m.Content, " ")
			var answerIndex int
			var datas string
			for i, word := range data {
				if word == "리챔아" {
					continue
				}
				if word == "배워" {
					answerIndex = i + 1
					continue
				}
			}
			for i := answerIndex + 1; i < len(data); i++ {
				datas += data[i] + " "
			}
			database.CreateMSG(s, m, d.DB, data[answerIndex], datas)
		}
	}
}

func (d *Message) SendingEmbed(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "$사용법":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("명령어 목록", "듀\n듀댜?\n\n리챔아 배워 ooo\nooo을 하면 ooo을 함(?)", 16705372))
	case "듀":
		s.ChannelMessageSend(m.ChannelID, "듀다다듀")
	case "듀다?":
		s.ChannelMessageSend(m.ChannelID, "듀댜됴디")
	}

}
