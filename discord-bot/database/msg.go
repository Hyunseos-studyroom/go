package database

import (
	"context"
	"discord-bot/types"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func CreateMSG(s *discordgo.Session, m *discordgo.MessageCreate, db *mongo.Client, title, data string) {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "시간zon 에러")
		return
	}

	currentTime := time.Now().In(location).Format("2006-08-15 15:05")

	upload := &types.CreateMSG{
		Title:     title,
		MSG:       data,
		Author:    m.Author.Username,
		CreatedAt: currentTime,
	}

	collection := db.Database("msg").Collection("msg")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err = collection.InsertOne(ctx, upload)
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "저장 실��")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "완벽하게 숙지 했다.")
}
