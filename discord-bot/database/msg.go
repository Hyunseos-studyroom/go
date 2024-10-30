package database

import (
	"context"
	"discord-bot/types"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
	"strings"
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

func DeleteMSG(s *discordgo.Session, m *discordgo.MessageCreate, db *mongo.Client) error {
	index := strings.Split(m.Content, " ")
	count, err := strconv.Atoi(index[2])
	if err != nil {
		return err
	}

	// Discord API의 제한에 따라 최대 100개 메시지를 한 번에 삭제할 수 있습니다.
	if count > 100 {
		count = 100
	}

	messages, err := s.ChannelMessages(m.ChannelID, count, "", "", "")
	if err != nil {
		return err
	}

	// 삭제할 메시지 IDs를 수집합니다.
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	// 메시지를 한 번에 삭제합니다.
	if len(messageIDs) > 0 {
		if err := s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs); err != nil {
			return err
		}
	}

	// 이미지 파일을 열고 전송합니다.
	image, err := os.Open("Untitled.jpg")
	if err != nil {
		return err
	}
	defer image.Close() // 파일을 사용한 후 반드시 닫아줍니다.

	_, err = s.ChannelFileSend(m.ChannelID, "Untitled.jpg", image)
	if err != nil {
		return err
	}

	return nil
}
