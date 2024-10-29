package database

import (
	"context"
	"discord-bot/types"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
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

func GetMSG(db *mongo.Client, word string) ([]types.CreateMSG, bool) {
	if word == "배워" {
		return nil, false
	} else {
		collection := db.Database("msg").Collection("msg")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.D{{"title", word}}
		cur, err := collection.Find(ctx, filter)
		if err != nil {
			log.Println(err)
			return nil, false
		}
		defer cur.Close(ctx)

		var results []types.CreateMSG
		if err := cur.All(ctx, &results); err != nil {
			log.Println(err)
			return nil, false
		}

		return results, true
	}
}
