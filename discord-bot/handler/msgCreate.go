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
	"time"
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
		case prefix + " 자폭해":
			s.ChannelMessageSend(m.ChannelID, "5")
			time.Sleep(1 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "4")
			time.Sleep(1 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "3")
			time.Sleep(1 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "2")
			time.Sleep(1 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "1")
			time.Sleep(1 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "펑 (리챔 터지는 소리)")
		}
		if strings.HasPrefix(content, "리챔아 유저생성") {
			database.SetGamer(s, m, d.DB)
		}
		if strings.HasPrefix(content, "리챔아 홀짝 시작") {
			database.StartLRGame(s, m, d.DB)
		}
		if strings.HasPrefix(content, "리챔아 삭제해") {
			err := database.DeleteMSG(s, m)
			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, "에러남 ㅋㅋㅋ")
			}
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
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("명령어 목록", "듀\n듀댜?\n\n리챔아 배워 ooo xxx\nxxx을 하면 ooo을 함(?)\n\n리챔아 삭제해 99\n최대 100개까지 메세지 삭제 ㄱㄴ\n\n내월급\n당신의 이번달 월급은?", 16705372))
	case "듀":
		s.ChannelMessageSend(m.ChannelID, "듀다다듀")
	case "듀다?":
		s.ChannelMessageSend(m.ChannelID, "듀댜됴디")
	case "듀댜":
		s.ChannelMessageSend(m.ChannelID, "유아퇴행장애인")
	}
}

func (d *Message) ThisMonthSalary(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ID == lastMessageID {
		return
	}
	if m.Content == "내월급" {
		rand.Seed(time.Now().UnixNano())

		// 2000000에서 5000000 사이의 난수 생성
		randomNum := rand.Intn(3000000) + 2000000

		// 만원 단위로 반올림
		randomNum = (randomNum / 10000) * 10000

		// 포맷팅
		formattedNum := fmt.Sprintf("%d", randomNum)
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("이번달 당신의 월급은?", fmt.Sprintf("<@%s> 님의 월급은 %s원 입니다.", m.Author.ID, formattedNum), 16705372))
	}
}
